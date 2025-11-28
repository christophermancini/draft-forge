package agent

import (
	"database/sql"

	"github.com/yourusername/draft-forge/internal/models"
)

type dbAgentRun struct {
	ID          int64          `db:"id"`
	ProjectID   int64          `db:"project_id"`
	AgentType   string         `db:"agent_type"`
	Trigger     string         `db:"trigger"`
	Status      string         `db:"status"`
	Results     []byte         `db:"results"`
	Error       sql.NullString `db:"error_message"`
	StartedAt   sql.NullTime   `db:"started_at"`
	CompletedAt sql.NullTime   `db:"completed_at"`
	CreatedAt   sql.NullTime   `db:"created_at"`
}

func (d dbAgentRun) toModel() models.AgentRun {
	run := models.AgentRun{
		ID:        d.ID,
		ProjectID: d.ProjectID,
		AgentType: d.AgentType,
		Trigger:   d.Trigger,
		Status:    d.Status,
	}
	if len(d.Results) > 0 {
		run.Results = d.Results
	}
	if d.Error.Valid {
		run.Error = d.Error.String
	}
	if d.StartedAt.Valid {
		run.StartedAt = &d.StartedAt.Time
	}
	if d.CompletedAt.Valid {
		run.CompletedAt = &d.CompletedAt.Time
	}
	if d.CreatedAt.Valid {
		run.CreatedAt = d.CreatedAt.Time
	}
	return run
}
