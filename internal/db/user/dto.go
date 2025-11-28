package user

import (
	"database/sql"

	"github.com/yourusername/draft-forge/internal/models"
)

type dbUser struct {
	ID           int64          `db:"id"`
	GitHubID     int64          `db:"github_id"`
	Username     string         `db:"username"`
	Email        sql.NullString `db:"email"`
	AvatarURL    sql.NullString `db:"avatar_url"`
	AccessToken  sql.NullString `db:"access_token"`
	RefreshToken sql.NullString `db:"refresh_token"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
}

func toDB(u models.User) dbUser {
	var email, avatar, access, refresh sql.NullString
	if u.Email != nil {
		email = sql.NullString{String: *u.Email, Valid: true}
	}
	if u.AvatarURL != "" {
		avatar = sql.NullString{String: u.AvatarURL, Valid: true}
	}
	if u.AccessToken != "" {
		access = sql.NullString{String: u.AccessToken, Valid: true}
	}
	if u.RefreshToken != "" {
		refresh = sql.NullString{String: u.RefreshToken, Valid: true}
	}

	return dbUser{
		ID:           u.ID,
		GitHubID:     u.GitHubID,
		Username:     u.Username,
		Email:        email,
		AvatarURL:    avatar,
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func (d dbUser) toModel() models.User {
	var email *string
	if d.Email.Valid {
		email = &d.Email.String
	}

	return models.User{
		ID:           d.ID,
		GitHubID:     d.GitHubID,
		Username:     d.Username,
		Email:        email,
		AvatarURL:    d.AvatarURL.String,
		AccessToken:  d.AccessToken.String,
		RefreshToken: d.RefreshToken.String,
	}
}
