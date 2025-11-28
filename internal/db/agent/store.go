package agent

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yourusername/draft-forge/internal/models"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertRun(ctx context.Context, run models.AgentRun) (models.AgentRun, error) {
	query := `
		INSERT INTO agent_runs (project_id, agent_type, trigger, status, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id, created_at
	`
	var dbRun dbAgentRun
	err := s.db.QueryRowxContext(ctx, query, run.ProjectID, run.AgentType, run.Trigger, run.Status).
		Scan(&dbRun.ID, &dbRun.CreatedAt)
	if err != nil {
		return models.AgentRun{}, fmt.Errorf("insert agent run: %w", err)
	}
	dbRun.ProjectID = run.ProjectID
	dbRun.AgentType = run.AgentType
	dbRun.Trigger = run.Trigger
	dbRun.Status = run.Status
	return dbRun.toModel(), nil
}

func (s *Store) MarkRunning(ctx context.Context, id int64, startedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent_runs SET status = $1, started_at = $2 WHERE id = $3
	`, "running", startedAt, id)
	if err != nil {
		return fmt.Errorf("mark running: %w", err)
	}
	return nil
}

func (s *Store) MarkCompleted(ctx context.Context, id int64, results json.RawMessage, completedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent_runs SET status = $1, results = $2, completed_at = $3 WHERE id = $4
	`, "completed", results, completedAt, id)
	if err != nil {
		return fmt.Errorf("mark completed: %w", err)
	}
	return nil
}

func (s *Store) MarkFailed(ctx context.Context, id int64, message string, completedAt time.Time) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE agent_runs SET status = $1, error_message = $2, completed_at = $3 WHERE id = $4
	`, "failed", message, completedAt, id)
	if err != nil {
		return fmt.Errorf("mark failed: %w", err)
	}
	return nil
}

func (s *Store) GetRun(ctx context.Context, id int64) (models.AgentRun, error) {
	var dbRun dbAgentRun
	err := s.db.GetContext(ctx, &dbRun, `
		SELECT id, project_id, agent_type, trigger, status, results, error_message, started_at, completed_at, created_at
		FROM agent_runs
		WHERE id = $1
	`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.AgentRun{}, models.ErrNotFound
		}
		return models.AgentRun{}, fmt.Errorf("get run: %w", err)
	}
	return dbRun.toModel(), nil
}

func (s *Store) ProjectExists(ctx context.Context, projectID int64) (bool, error) {
	var exists bool
	if err := s.db.GetContext(ctx, &exists, `SELECT EXISTS (SELECT 1 FROM projects WHERE id = $1)`, projectID); err != nil {
		return false, fmt.Errorf("check project exists: %w", err)
	}
	return exists, nil
}
