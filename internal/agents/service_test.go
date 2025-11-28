package agents

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/draft-forge/internal/models"
)

func TestQueueRunWritesArtifactAndCompletes(t *testing.T) {
	tempDir := t.TempDir()
	store := newMockStore()
	svc := NewService(store, tempDir)
	fixedTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixedTime }

	run, err := svc.QueueRun(context.Background(), RunRequest{
		ProjectID:    1,
		AgentType:    "continuity",
		Trigger:      "pr",
		FilesChanged: []string{"chapters/01.md"},
	})
	if err != nil {
		t.Fatalf("QueueRun returned error: %v", err)
	}

	if run.Status != "completed" {
		t.Fatalf("expected status %s, got %s", "completed", run.Status)
	}
	if run.Results == nil {
		t.Fatal("expected results to be set")
	}
	if run.StartedAt == nil || run.CompletedAt == nil {
		t.Fatal("expected timestamps to be set")
	}

	artifactPath := filepath.Join(tempDir, "run-1-continuity.json")
	data, err := os.ReadFile(artifactPath)
	if err != nil {
		t.Fatalf("failed to read artifact: %v", err)
	}

	var payload struct {
		RunID     int64  `json:"run_id"`
		ProjectID int64  `json:"project_id"`
		AgentType string `json:"agent_type"`
		Trigger   string `json:"trigger"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("failed to unmarshal artifact: %v", err)
	}

	if payload.RunID != run.ID || payload.ProjectID != run.ProjectID {
		t.Fatalf("artifact contents mismatch run/project ids: %+v", payload)
	}
	if payload.AgentType != "continuity" || payload.Trigger != "pr" {
		t.Fatalf("artifact agent/trigger mismatch: %+v", payload)
	}
}

func TestQueueRunRejectsInvalidAgent(t *testing.T) {
	store := newMockStore()
	svc := NewService(store, t.TempDir())

	_, err := svc.QueueRun(context.Background(), RunRequest{
		ProjectID: 1,
		AgentType: "unknown",
	})
	if err == nil || err != ErrInvalidAgentType {
		t.Fatalf("expected ErrInvalidAgentType, got %v", err)
	}
}

type mockStore struct {
	nextID   int64
	runs     map[int64]models.AgentRun
	projects map[int64]bool
}

func newMockStore() *mockStore {
	return &mockStore{
		nextID:   1,
		runs:     make(map[int64]models.AgentRun),
		projects: map[int64]bool{1: true},
	}
}

func (m *mockStore) ProjectExists(_ context.Context, projectID int64) (bool, error) {
	return m.projects[projectID], nil
}

func (m *mockStore) InsertRun(_ context.Context, run models.AgentRun) (models.AgentRun, error) {
	run.ID = m.nextID
	m.nextID++
	run.CreatedAt = time.Now()
	m.runs[run.ID] = run
	return run, nil
}

func (m *mockStore) MarkRunning(_ context.Context, id int64, startedAt time.Time) error {
	run := m.runs[id]
	run.Status = "running"
	run.StartedAt = &startedAt
	m.runs[id] = run
	return nil
}

func (m *mockStore) MarkCompleted(_ context.Context, id int64, results json.RawMessage, completedAt time.Time) error {
	run := m.runs[id]
	run.Status = "completed"
	run.Results = results
	run.CompletedAt = &completedAt
	m.runs[id] = run
	return nil
}

func (m *mockStore) MarkFailed(_ context.Context, id int64, message string, completedAt time.Time) error {
	run := m.runs[id]
	run.Status = "failed"
	run.Error = message
	run.CompletedAt = &completedAt
	m.runs[id] = run
	return nil
}

func (m *mockStore) GetRun(_ context.Context, id int64) (models.AgentRun, error) {
	run, ok := m.runs[id]
	if !ok {
		return models.AgentRun{}, ErrRunNotFound
	}
	return run, nil
}
