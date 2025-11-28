package auth

type GitHubUser struct {
	ID        int64
	Login     string
	Email     string
	AvatarURL string
}
