package agent

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/yourusername/draft-forge/internal/models"
)

func (s *Store) ListRuns(ctx context.Context, projectID int64) ([]models.AgentRun, error) {
	query := `
		SELECT id, project_id, agent_type, trigger, status, results, error_message, started_at, completed_at, created_at
		FROM agent_runs
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	var runs []dbAgentRun
	if err := s.db.SelectContext(ctx, &runs, query, projectID); err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	out := make([]models.AgentRun, 0, len(runs))
	for _, r := range runs {
		out = append(out, r.toModel())
	}
	return out, nil
}
