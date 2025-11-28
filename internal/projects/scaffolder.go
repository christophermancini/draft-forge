package projects

import (
	"context"

	"github.com/yourusername/draft-forge/internal/models"
)

// ScaffoldInput carries parameters for scaffolding a project.
type ScaffoldInput struct {
	Project     models.Project
	GitHubToken string
	GitHubOwner string
	Template    models.ProjectTemplate
}

// ScaffoldResult describes the outcome of scaffolding.
type ScaffoldResult struct {
	Path    string `json:"path,omitempty"`     // local path if applicable
	RepoURL string `json:"repo_url,omitempty"` // GitHub repo URL if applicable
}

// Scaffolder creates the initial project structure (locally or remotely).
type Scaffolder interface {
	Scaffold(ctx context.Context, input ScaffoldInput) (ScaffoldResult, error)
}
