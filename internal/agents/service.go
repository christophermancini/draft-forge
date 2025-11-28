package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/draft-forge/internal/models"
)

// Supported agent types for the PoC.
var validAgentTypes = map[string]bool{
	"continuity": true,
	"style":      true,
	"timeline":   true,
	"fact":       true,
}

// Supported triggers.
var validTriggers = map[string]bool{
	"manual":    true,
	"commit":    true,
	"pr":        true,
	"scheduled": true,
}

var (
	ErrInvalidAgentType = errors.New("invalid agent type")
	ErrInvalidTrigger   = errors.New("invalid trigger")
	ErrProjectNotFound  = errors.New("project not found")
	ErrRunNotFound      = errors.New("run not found")
)

type RunRequest struct {
	ProjectID    int64
	AgentType    string
	Trigger      string
	FilesChanged []string
}

// RunStore defines persistence needs for agent runs and project existence.
type RunStore interface {
	InsertRun(ctx context.Context, run models.AgentRun) (models.AgentRun, error)
	MarkRunning(ctx context.Context, id int64, startedAt time.Time) error
	MarkCompleted(ctx context.Context, id int64, results json.RawMessage, completedAt time.Time) error
	MarkFailed(ctx context.Context, id int64, message string, completedAt time.Time) error
	GetRun(ctx context.Context, id int64) (models.AgentRun, error)
	ListRuns(ctx context.Context, projectID int64) ([]models.AgentRun, error)
	ProjectExists(ctx context.Context, projectID int64) (bool, error)
}

type Service struct {
	store       RunStore
	artifactDir string
	now         func() time.Time
}

func NewService(store RunStore, artifactDir string) *Service {
	return &Service{
		store:       store,
		artifactDir: artifactDir,
		now:         time.Now,
	}
}

// QueueRun validates input, persists a queued run, executes a stub agent, and returns the completed run.
// For PoC the execution happens synchronously and writes a JSON artifact to the repo path.
func (s *Service) QueueRun(ctx context.Context, req RunRequest) (models.AgentRun, error) {
	if !validAgentTypes[req.AgentType] {
		return models.AgentRun{}, ErrInvalidAgentType
	}
	trigger := req.Trigger
	if trigger == "" {
		trigger = "manual"
	}
	if !validTriggers[trigger] {
		return models.AgentRun{}, ErrInvalidTrigger
	}

	exists, err := s.store.ProjectExists(ctx, req.ProjectID)
	if err != nil {
		return models.AgentRun{}, fmt.Errorf("check project: %w", err)
	}
	if !exists {
		return models.AgentRun{}, ErrProjectNotFound
	}

	run := models.AgentRun{
		ProjectID: req.ProjectID,
		AgentType: req.AgentType,
		Trigger:   trigger,
		Status:    "queued",
		CreatedAt: s.now(),
	}

	run, err = s.store.InsertRun(ctx, run)
	if err != nil {
		return models.AgentRun{}, fmt.Errorf("insert run: %w", err)
	}

	if err := s.executeRun(ctx, run, req.FilesChanged); err != nil {
		return models.AgentRun{}, err
	}

	finalRun, err := s.store.GetRun(ctx, run.ID)
	if err != nil {
		return models.AgentRun{}, fmt.Errorf("fetch run: %w", err)
	}
	return finalRun, nil
}

// GetRun returns a run by ID.
func (s *Service) GetRun(ctx context.Context, id int64) (models.AgentRun, error) {
	return s.store.GetRun(ctx, id)
}

func (s *Service) ListRuns(ctx context.Context, projectID int64) ([]models.AgentRun, error) {
	return s.store.ListRuns(ctx, projectID)
}

func (s *Service) executeRun(ctx context.Context, run models.AgentRun, files []string) error {
	started := s.now()
	if err := s.store.MarkRunning(ctx, run.ID, started); err != nil {
		return fmt.Errorf("mark running: %w", err)
	}

	resultsPayload := map[string]any{
		"summary": fmt.Sprintf("%s agent completed", run.AgentType),
		"issues": []map[string]any{
			{
				"severity": "info",
				"message":  "PoC agent run completed successfully.",
			},
		},
		"stats": map[string]any{
			"files_analyzed": len(files),
			"tokens_used":    0,
		},
	}

	resultBytes, err := json.Marshal(resultsPayload)
	if err != nil {
		failErr := fmt.Errorf("marshal results: %w", err)
		_ = s.store.MarkFailed(ctx, run.ID, failErr.Error(), s.now())
		return failErr
	}

	if err := s.writeArtifact(run, resultBytes); err != nil {
		failErr := fmt.Errorf("write artifact: %w", err)
		_ = s.store.MarkFailed(ctx, run.ID, failErr.Error(), s.now())
		return failErr
	}

	if err := s.store.MarkCompleted(ctx, run.ID, resultBytes, s.now()); err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}

	return nil
}

func (s *Service) writeArtifact(run models.AgentRun, results []byte) error {
	if s.artifactDir == "" {
		return errors.New("artifact directory not configured")
	}

	if err := os.MkdirAll(s.artifactDir, 0o755); err != nil {
		return fmt.Errorf("create artifact dir: %w", err)
	}

	filename := fmt.Sprintf("run-%d-%s.json", run.ID, run.AgentType)
	path := filepath.Join(s.artifactDir, filename)

	payload := map[string]any{
		"run_id":     run.ID,
		"project_id": run.ProjectID,
		"agent_type": run.AgentType,
		"trigger":    run.Trigger,
		"timestamp":  s.now().UTC().Format(time.RFC3339),
		"results":    json.RawMessage(results),
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal artifact: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write artifact file: %w", err)
	}

	return nil
}
