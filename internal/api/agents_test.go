package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/agents"
	"github.com/yourusername/draft-forge/internal/models"
)

func TestQueueRunHandlerSuccess(t *testing.T) {
	app := fiber.New()
	handler := NewAgentHandler(&stubAgentService{
		queueFunc: func(ctx context.Context, req agents.RunRequest) (models.AgentRun, error) {
			return models.AgentRun{ID: 42, Status: "queued"}, nil
		},
	})
	handler.Register(app)

	body := []byte(`{"agent_type":"continuity","trigger":"pr"}`)
	req := httptest.NewRequest(http.MethodPost, "/projects/1/agents/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := payload["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data envelope, got %v", payload)
	}
	if data["id"] != float64(42) {
		t.Fatalf("expected run id 42, got %v", data["id"])
	}
}

func TestQueueRunHandlerInvalidAgent(t *testing.T) {
	app := fiber.New()
	handler := NewAgentHandler(&stubAgentService{
		queueFunc: func(ctx context.Context, req agents.RunRequest) (models.AgentRun, error) {
			return models.AgentRun{}, agents.ErrInvalidAgentType
		},
	})
	handler.Register(app)

	body := []byte(`{"agent_type":"bad"}`)
	req := httptest.NewRequest(http.MethodPost, "/projects/1/agents/run", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetRunHandlerNotFound(t *testing.T) {
	app := fiber.New()
	handler := NewAgentHandler(&stubAgentService{
		getFunc: func(ctx context.Context, id int64) (models.AgentRun, error) {
			return models.AgentRun{}, agents.ErrRunNotFound
		},
	})
	handler.Register(app)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/agents/runs/99", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}
}

func TestGetRunHandlerSuccess(t *testing.T) {
	app := fiber.New()
	handler := NewAgentHandler(&stubAgentService{
		getFunc: func(ctx context.Context, id int64) (models.AgentRun, error) {
			return models.AgentRun{
				ID:        id,
				ProjectID: 1,
				AgentType: "continuity",
				Status:    "completed",
			}, nil
		},
	})
	handler.Register(app)

	req := httptest.NewRequest(http.MethodGet, "/projects/1/agents/runs/5", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var payload struct {
		Data models.AgentRun `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Data.ID != 5 || payload.Data.ProjectID != 1 {
		t.Fatalf("unexpected run payload: %+v", payload.Data)
	}
}

type stubAgentService struct {
	queueFunc func(ctx context.Context, req agents.RunRequest) (models.AgentRun, error)
	getFunc   func(ctx context.Context, id int64) (models.AgentRun, error)
}

func (s *stubAgentService) QueueRun(ctx context.Context, req agents.RunRequest) (models.AgentRun, error) {
	if s.queueFunc == nil {
		return models.AgentRun{}, nil
	}
	return s.queueFunc(ctx, req)
}

func (s *stubAgentService) GetRun(ctx context.Context, id int64) (models.AgentRun, error) {
	if s.getFunc == nil {
		return models.AgentRun{}, nil
	}
	return s.getFunc(ctx, id)
}
