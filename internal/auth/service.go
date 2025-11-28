package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/yourusername/draft-forge/internal/models"
)

var (
	ErrMissingCode  = errors.New("missing code")
	ErrInvalidState = errors.New("invalid state")
)

type Service struct {
	store       Store
	github      GitHubClient
	tokens      *TokenManager
	clientID    string
	redirectURI string
	stateSecret []byte
	now         func() time.Time
}

func NewService(store Store, gh GitHubClient, tokens *TokenManager, clientID, redirectURI, stateSecret string) *Service {
	return &Service{
		store:       store,
		github:      gh,
		tokens:      tokens,
		clientID:    clientID,
		redirectURI: redirectURI,
		stateSecret: []byte(stateSecret),
		now:         time.Now,
	}
}

type AuthStart struct {
	AuthURL string `json:"auth_url"`
	State   string `json:"state"`
}

type AuthResult struct {
	User  models.User `json:"user"`
	Token TokenPair   `json:"token"`
}

func (s *Service) StartAuth() (AuthStart, error) {
	state := s.signState(fmt.Sprintf("%d", s.now().UnixNano()))
	authURL := s.buildAuthURL(state)
	return AuthStart{AuthURL: authURL, State: state}, nil
}

func (s *Service) CompleteAuth(ctx context.Context, code, state string) (AuthResult, error) {
	if code == "" {
		return AuthResult{}, ErrMissingCode
	}
	if !s.verifyState(state) {
		return AuthResult{}, ErrInvalidState
	}

	accessToken, err := s.github.ExchangeCode(ctx, code)
	if err != nil {
		return AuthResult{}, fmt.Errorf("exchange code: %w", err)
	}

	ghUser, err := s.github.GetUser(ctx, accessToken)
	if err != nil {
		return AuthResult{}, fmt.Errorf("fetch github user: %w", err)
	}

	user, err := s.store.UpsertGitHubUser(ctx, ghUser, accessToken, "")
	if err != nil {
		return AuthResult{}, fmt.Errorf("persist user: %w", err)
	}

	access, err := s.tokens.SignAccessToken(user)
	if err != nil {
		return AuthResult{}, fmt.Errorf("sign access token: %w", err)
	}
	refresh, err := s.tokens.SignRefreshToken(user)
	if err != nil {
		return AuthResult{}, fmt.Errorf("sign refresh token: %w", err)
	}

	return AuthResult{
		User: user,
		Token: TokenPair{
			AccessToken:  access,
			RefreshToken: refresh,
			ExpiresIn:    int64(s.tokens.accessTTL.Seconds()),
		},
	}, nil
}

func (s *Service) buildAuthURL(state string) string {
	v := url.Values{}
	v.Set("client_id", s.clientID)
	if s.redirectURI != "" {
		v.Set("redirect_uri", s.redirectURI)
	}
	v.Set("scope", "read:user user:email")
	v.Set("state", state)

	return "https://github.com/login/oauth/authorize?" + v.Encode()
}

func (s *Service) signState(value string) string {
	mac := hmac.New(sha256.New, s.stateSecret)
	_, _ = mac.Write([]byte(value))
	sum := mac.Sum(nil)
	return value + "." + hex.EncodeToString(sum)
}

func (s *Service) verifyState(state string) bool {
	parts := []rune(state)
	idx := -1
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i] == '.' {
			idx = i
			break
		}
	}
	if idx <= 0 || idx >= len(parts)-1 {
		return false
	}
	value := string(parts[:idx])
	sig := string(parts[idx+1:])

	expected := s.signState(value)
	return state == expected && sig != ""
}
