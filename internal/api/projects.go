package api

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/models"
	"github.com/yourusername/draft-forge/internal/projects"
)

type ProjectService interface {
	Create(ctx context.Context, req projects.CreateRequest) (models.Project, projects.ScaffoldResult, error)
	List(ctx context.Context, userID int64) ([]models.Project, error)
}

type projectHandler struct {
	service *projects.Service
}

func NewProjectHandler(service *projects.Service) *projectHandler {
	return &projectHandler{service: service}
}

func (h *projectHandler) Register(app fiber.Router) {
	app.Post("/projects", h.create)
	app.Get("/projects", h.list)
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectType string `json:"project_type"`
	UseGitHub   bool   `json:"use_github"`
	GitHubOwner string `json:"github_owner"`
	Template    string `json:"template"`
}

func (h *projectHandler) create(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	var req createProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	project, scaffoldResult, err := h.service.Create(c.Context(), projects.CreateRequest{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		ProjectType: req.ProjectType,
		Template:    models.ProjectTemplate(req.Template),
		GitHubToken: func() string {
			if req.UseGitHub {
				return githubTokenFromContext(c)
			}
			return ""
		}(),
		GitHubOwner: req.GitHubOwner,
	})
	if err != nil {
		if err == projects.ErrInvalidName {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return err
	}

	meta := fiber.Map{}
	if scaffoldResult.Path != "" {
		meta["scaffold_path"] = scaffoldResult.Path
	}
	if scaffoldResult.RepoURL != "" {
		meta["repo_url"] = scaffoldResult.RepoURL
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": project,
		"meta": meta,
	})
}

func (h *projectHandler) list(c *fiber.Ctx) error {
	userID, err := extractUserID(c)
	if err != nil {
		return err
	}

	projects, err := h.service.List(c.Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"data": projects,
		"meta": fiber.Map{
			"count": len(projects),
		},
	})
}

func extractUserID(c *fiber.Ctx) (int64, error) {
	userIDVal := c.Locals("user_id")
	userID, ok := userIDVal.(int64)
	if !ok || userID <= 0 {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "unauthorized")
	}
	return userID, nil
}

func githubTokenFromContext(c *fiber.Ctx) string {
	tokenVal := c.Locals("github_token")
	if tokenVal == nil {
		return ""
	}
	if token, ok := tokenVal.(string); ok {
		return token
	}
	return ""
}
