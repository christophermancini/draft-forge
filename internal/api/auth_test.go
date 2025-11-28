package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/yourusername/draft-forge/internal/auth"
	"github.com/yourusername/draft-forge/internal/models"
)

func TestAuthStartHandler(t *testing.T) {
	app, _, _ := setupAuthApp(t)
	req := httptest.NewRequest(http.MethodGet, "/auth/github/start", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if _, ok := payload["data"]; !ok {
		t.Fatalf("expected data envelope, got %v", payload)
	}
}

func TestAuthMeHandler(t *testing.T) {
	app, tokenMgr, store := setupAuthApp(t)

	user := models.User{ID: 1, GitHubID: 99, Username: "octo"}
	store.user = user

	token, err := tokenMgr.SignAccessToken(user)
	if err != nil {
		t.Fatalf("SignAccessToken error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func setupAuthApp(t *testing.T) (*fiber.App, *auth.TokenManager, *stubUserStore) {
	t.Helper()

	tokenMgr := auth.NewTokenManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	store := &stubUserStore{}
	gh := &stubGitHubClient{
		accessToken: "gh-token",
		user:        auth.GitHubUser{ID: 99, Login: "octo"},
	}
	svc := auth.NewService(store, gh, tokenMgr, "client-id", "http://cb", "state-secret")
	handler := NewAuthHandler(svc, tokenMgr, store)

	app := fiber.New()
	handler.Register(app)
	return app, tokenMgr, store
}

type stubUserStore struct {
	user models.User
}

func (s *stubUserStore) UpsertGitHubUser(ctx context.Context, gh auth.GitHubUser, accessToken, refreshToken string) (models.User, error) {
	s.user = models.User{
		ID:           1,
		GitHubID:     gh.ID,
		Username:     gh.Login,
		Email:        &gh.Email,
		AvatarURL:    gh.AvatarURL,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return s.user, nil
}

func (s *stubUserStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	if s.user.ID == 0 {
		s.user.ID = id
	}
	return s.user, nil
}

type stubGitHubClient struct {
	accessToken string
	user        auth.GitHubUser
}

func (c *stubGitHubClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	return c.accessToken, nil
}

func (c *stubGitHubClient) GetUser(ctx context.Context, accessToken string) (auth.GitHubUser, error) {
	return c.user, nil
}
