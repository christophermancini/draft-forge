package auth

import (
	"context"
	"testing"
	"time"

	"github.com/yourusername/draft-forge/internal/models"
)

func TestStartAuthGeneratesURLAndState(t *testing.T) {
	store := &stubUserStore{}
	gh := &stubGitHubClient{}
	tokens := NewTokenManager("access", "refresh", time.Minute, time.Hour)
	svc := NewService(store, gh, tokens, "client-id", "https://example.com/callback", "state-secret")
	svc.now = func() time.Time { return time.Unix(0, 123) }

	start, err := svc.StartAuth()
	if err != nil {
		t.Fatalf("StartAuth error: %v", err)
	}
	if start.State == "" {
		t.Fatal("expected state to be set")
	}
	if start.AuthURL == "" {
		t.Fatal("expected auth_url to be set")
	}
}

func TestCompleteAuthSuccess(t *testing.T) {
	store := &stubUserStore{}
	gh := &stubGitHubClient{
		accessToken: "gh-token",
		user: GitHubUser{
			ID:        99,
			Login:     "octo",
			Email:     "octo@example.com",
			AvatarURL: "http://avatar",
		},
	}
	tokens := NewTokenManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	svc := NewService(store, gh, tokens, "client-id", "http://cb", "state-secret")

	state := svc.signState("nonce")
	result, err := svc.CompleteAuth(context.Background(), "code123", state)
	if err != nil {
		t.Fatalf("CompleteAuth error: %v", err)
	}
	if result.User.GitHubID != 99 || result.User.Username != "octo" {
		t.Fatalf("unexpected user %+v", result.User)
	}
	if result.Token.AccessToken == "" || result.Token.RefreshToken == "" {
		t.Fatal("expected tokens to be set")
	}
}

func TestCompleteAuthInvalidState(t *testing.T) {
	store := &stubUserStore{}
	gh := &stubGitHubClient{}
	tokens := NewTokenManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	svc := NewService(store, gh, tokens, "client-id", "http://cb", "state-secret")

	_, err := svc.CompleteAuth(context.Background(), "code123", "badstate")
	if err == nil || err != ErrInvalidState {
		t.Fatalf("expected ErrInvalidState, got %v", err)
	}
}

type stubUserStore struct {
	user models.User
}

func (s *stubUserStore) UpsertGitHubUser(ctx context.Context, gh GitHubUser, accessToken, refreshToken string) (models.User, error) {
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
	return s.user, nil
}

type stubGitHubClient struct {
	accessToken string
	user        GitHubUser
}

func (c *stubGitHubClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	return c.accessToken, nil
}

func (c *stubGitHubClient) GetUser(ctx context.Context, accessToken string) (GitHubUser, error) {
	return c.user, nil
}
