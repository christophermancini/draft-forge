package models

type RepoInfo struct {
	ID   *int64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type Project struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	ProjectType string    `json:"project_type"`
	GitHubRepo  *RepoInfo `json:"github_repo,omitempty"`
}

type ProjectTemplate string

const (
	TemplateNovel ProjectTemplate = "novel"
)
