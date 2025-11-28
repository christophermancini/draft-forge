package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/yourusername/draft-forge/internal/models"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) UpsertGitHubUser(ctx context.Context, gh models.User, accessToken, refreshToken string) (models.User, error) {
	dbu := toDB(models.User{
		GitHubID:     gh.GitHubID,
		Username:     gh.Username,
		Email:        gh.Email,
		AvatarURL:    gh.AvatarURL,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	query := `
		INSERT INTO users (github_id, username, email, avatar_url, access_token, refresh_token)
		VALUES (:github_id, :username, NULLIF(:email, ''), NULLIF(:avatar_url, ''), :access_token, :refresh_token)
		ON CONFLICT (github_id)
		DO UPDATE SET username = EXCLUDED.username,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, github_id, username, email, avatar_url, access_token, refresh_token, created_at, updated_at
	`

	rows, err := s.db.NamedQueryContext(ctx, query, dbu)
	if err != nil {
		return models.User{}, fmt.Errorf("upsert user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&dbu.ID,
			&dbu.GitHubID,
			&dbu.Username,
			&dbu.Email,
			&dbu.AvatarURL,
			&dbu.AccessToken,
			&dbu.RefreshToken,
			&dbu.CreatedAt,
			&dbu.UpdatedAt,
		); err != nil {
			return models.User{}, fmt.Errorf("scan user: %w", err)
		}
	}

	return dbu.toModel(), nil
}

func (s *Store) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	var dbu dbUser
	if err := s.db.GetContext(ctx, &dbu, `
		SELECT id, github_id, username, email, avatar_url, access_token, refresh_token, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id); err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	return dbu.toModel(), nil
}
