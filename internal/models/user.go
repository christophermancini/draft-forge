package models

type User struct {
	ID           int64   `json:"id"`
	GitHubID     int64   `json:"github_id"`
	Username     string  `json:"username"`
	Email        *string `json:"email,omitempty"`
	AvatarURL    string  `json:"avatar_url,omitempty"`
	AccessToken  string  `json:"-"`
	RefreshToken string  `json:"-"`
}
