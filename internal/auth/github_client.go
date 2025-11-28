package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type GitHubClient interface {
	ExchangeCode(ctx context.Context, code string) (string, error)
	GetUser(ctx context.Context, accessToken string) (GitHubUser, error)
}

type OAuthClient struct {
	httpClient   *http.Client
	clientID     string
	clientSecret string
	redirectURI  string
}

func NewOAuthClient(httpClient *http.Client, clientID, clientSecret, redirectURI string) *OAuthClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &OAuthClient{
		httpClient:   httpClient,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
	}
}

func (c *OAuthClient) ExchangeCode(ctx context.Context, code string) (string, error) {
	form := url.Values{}
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)
	form.Set("code", code)
	if c.redirectURI != "" {
		form.Set("redirect_uri", c.redirectURI)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://github.com/login/oauth/access_token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d from GitHub token endpoint", resp.StatusCode)
	}

	var payload struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode token response: %w", err)
	}
	if payload.Error != "" {
		return "", fmt.Errorf("github token error: %s", payload.Error)
	}
	if payload.AccessToken == "" {
		return "", fmt.Errorf("empty access token from GitHub")
	}
	return payload.AccessToken, nil
}

func (c *OAuthClient) GetUser(ctx context.Context, accessToken string) (GitHubUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("create user request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return GitHubUser{}, fmt.Errorf("fetch user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GitHubUser{}, fmt.Errorf("unexpected status %d from GitHub user API", resp.StatusCode)
	}

	var payload struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return GitHubUser{}, fmt.Errorf("decode user response: %w", err)
	}
	return GitHubUser{
		ID:        payload.ID,
		Login:     payload.Login,
		Email:     payload.Email,
		AvatarURL: payload.AvatarURL,
	}, nil
}
