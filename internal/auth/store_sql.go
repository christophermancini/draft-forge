package auth

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/yourusername/draft-forge/internal/db/user"
	"github.com/yourusername/draft-forge/internal/models"
)

type Store interface {
	UpsertGitHubUser(ctx context.Context, gh GitHubUser, accessToken, refreshToken string) (models.User, error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
}

type SQLStore struct {
	inner *user.Store
}

func NewSQLUserStore(db *sql.DB) *SQLStore {
	return &SQLStore{inner: user.NewStore(sqlx.NewDb(db, "postgres"))}
}

func (s *SQLStore) UpsertGitHubUser(ctx context.Context, gh GitHubUser, accessToken, refreshToken string) (models.User, error) {
	return s.inner.UpsertGitHubUser(ctx, models.User{
		GitHubID:  gh.ID,
		Username:  gh.Login,
		Email:     &gh.Email,
		AvatarURL: gh.AvatarURL,
	}, accessToken, refreshToken)
}

func (s *SQLStore) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	return s.inner.GetUserByID(ctx, id)
}
