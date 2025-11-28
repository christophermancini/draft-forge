package models

import (
	"encoding/json"
	"time"
)

type AgentRun struct {
	ID          int64           `json:"id"`
	ProjectID   int64           `json:"project_id"`
	AgentType   string          `json:"agent_type"`
	Trigger     string          `json:"trigger"`
	Status      string          `json:"status"`
	Results     json.RawMessage `json:"results,omitempty"`
	Error       string          `json:"error_message,omitempty"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}
