package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yourusername/draft-forge/internal/models"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	now           func() time.Time
}

func NewTokenManager(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		now:           time.Now,
	}
}

type userClaims struct {
	UserID   int64 `json:"uid"`
	GitHubID int64 `json:"gid"`
	jwt.RegisteredClaims
}

func (tm *TokenManager) SignAccessToken(user models.User) (string, error) {
	claims := userClaims{
		UserID:   user.ID,
		GitHubID: user.GitHubID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tm.now().Add(tm.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(tm.now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.accessSecret)
}

func (tm *TokenManager) SignRefreshToken(user models.User) (string, error) {
	claims := userClaims{
		UserID:   user.ID,
		GitHubID: user.GitHubID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tm.now().Add(tm.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(tm.now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.refreshSecret)
}

func (tm *TokenManager) ParseAccessToken(tokenStr string) (userClaims, error) {
	return parseToken(tokenStr, tm.accessSecret)
}

func parseToken(tokenStr string, secret []byte) (userClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return userClaims{}, err
	}
	claims, ok := token.Claims.(*userClaims)
	if !ok || !token.Valid {
		return userClaims{}, errors.New("invalid token claims")
	}
	return *claims, nil
}
