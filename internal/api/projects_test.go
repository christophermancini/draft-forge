package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/auth"
	"github.com/yourusername/draft-forge/internal/models"
	"github.com/yourusername/draft-forge/internal/projects"
)

func TestCreateProjectHandler(t *testing.T) {
	tokenMgr := auth.NewTokenManager("a", "b", time.Minute, time.Hour)
	store := &stubProjectStore{}
	service := projects.NewService(store, nil)
	handler := NewProjectHandler(service)

	app := fiber.New()
	protected := app.Group("", AuthMiddleware(tokenMgr, &projectUserStore{
		user: models.User{ID: 1, GitHubID: 9, AccessToken: "gh-token"},
	}))
	handler.Register(protected)

	user := models.User{ID: 1, GitHubID: 9}
	token, _ := tokenMgr.SignAccessToken(user)

	body := []byte(`{"name":"Test Book","project_type":"novel","template":"novel"}`)
	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
}

func TestListProjectsHandler(t *testing.T) {
	tokenMgr := auth.NewTokenManager("a", "b", time.Minute, time.Hour)
	store := &stubProjectStore{projects: []models.Project{{ID: 1, UserID: 1, Name: "A", Slug: "a"}}}
	service := projects.NewService(store, nil)
	handler := NewProjectHandler(service)

	app := fiber.New()
	protected := app.Group("", AuthMiddleware(tokenMgr, &projectUserStore{
		user: models.User{ID: 1, GitHubID: 9, AccessToken: "gh-token"},
	}))
	handler.Register(protected)

	user := models.User{ID: 1, GitHubID: 9}
	token, _ := tokenMgr.SignAccessToken(user)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var payload struct {
		Data []models.Project `json:"data"`
		Meta map[string]any   `json:"meta"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(payload.Data) != 1 {
		t.Fatalf("expected 1 project, got %d", len(payload.Data))
	}
}

type stubProjectStore struct {
	projects []models.Project
}

func (s *stubProjectStore) InsertProject(ctx context.Context, p models.Project) (models.Project, error) {
	p.ID = 1
	s.projects = append(s.projects, p)
	return p, nil
}

func (s *stubProjectStore) ListProjects(ctx context.Context, userID int64) ([]models.Project, error) {
	return s.projects, nil
}

func (s *stubProjectStore) UpdateRepoInfo(ctx context.Context, projectID int64, repo models.RepoInfo) error {
	if len(s.projects) > 0 {
		s.projects[0].GitHubRepo = &repo
	}
	return nil
}

type projectUserStore struct {
	user models.User
}

func (s *projectUserStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	return s.user, nil
}

func (s *projectUserStore) UpsertGitHubUser(ctx context.Context, gh auth.GitHubUser, accessToken, refreshToken string) (models.User, error) {
	return s.user, nil
}
