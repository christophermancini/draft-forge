package api

import (
	"context"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/agents"
	"github.com/yourusername/draft-forge/internal/models"
)

type AgentService interface {
	QueueRun(ctx context.Context, req agents.RunRequest) (models.AgentRun, error)
	GetRun(ctx context.Context, id int64) (models.AgentRun, error)
	ListRuns(ctx context.Context, projectID int64) ([]models.AgentRun, error)
}

type AgentHandler struct {
	service AgentService
}

func NewAgentHandler(service AgentService) *AgentHandler {
	return &AgentHandler{service: service}
}

func (h *AgentHandler) Register(app fiber.Router) {
	app.Post("/projects/:projectID/agents/run", h.queueRun)
	app.Get("/projects/:projectID/agents/runs/:runID", h.getRun)
	app.Get("/projects/:projectID/agents/runs", h.listRuns)
}

type queueRunRequest struct {
	AgentType    string   `json:"agent_type"`
	Trigger      string   `json:"trigger"`
	FilesChanged []string `json:"files_changed"`
}

func (h *AgentHandler) queueRun(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("projectID"), 10, 64)
	if err != nil || projectID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	var req queueRunRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	run, err := h.service.QueueRun(c.Context(), agents.RunRequest{
		ProjectID:    projectID,
		AgentType:    req.AgentType,
		Trigger:      req.Trigger,
		FilesChanged: req.FilesChanged,
	})
	if err != nil {
		switch {
		case errors.Is(err, agents.ErrInvalidAgentType), errors.Is(err, agents.ErrInvalidTrigger):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case errors.Is(err, agents.ErrProjectNotFound):
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		default:
			return err
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"data": run,
		"meta": fiber.Map{
			"message": "Agent run queued successfully",
		},
	})
}

func (h *AgentHandler) getRun(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("projectID"), 10, 64)
	if err != nil || projectID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	runID, err := strconv.ParseInt(c.Params("runID"), 10, 64)
	if err != nil || runID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid run id")
	}

	run, err := h.service.GetRun(c.Context(), runID)
	if err != nil {
		if errors.Is(err, agents.ErrProjectNotFound) || errors.Is(err, agents.ErrRunNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "run not found")
		}
		return err
	}

	if run.ProjectID != projectID {
		return fiber.NewError(fiber.StatusNotFound, "run not found")
	}

	return c.JSON(fiber.Map{"data": run})
}

func (h *AgentHandler) listRuns(c *fiber.Ctx) error {
	projectID, err := strconv.ParseInt(c.Params("projectID"), 10, 64)
	if err != nil || projectID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid project id")
	}

	runs, err := h.service.ListRuns(c.Context(), projectID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": runs,
		"meta": fiber.Map{"count": len(runs)},
	})
}
