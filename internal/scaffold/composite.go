package scaffold

import (
	"context"
	"errors"

	"github.com/yourusername/draft-forge/internal/projects"
)

// CompositeScaffolder tries remote (GitHub) scaffolder when a token is provided, otherwise falls back to local.
type CompositeScaffolder struct {
	Remote projects.Scaffolder
	Local  projects.Scaffolder
}

func (c *CompositeScaffolder) Scaffold(ctx context.Context, input projects.ScaffoldInput) (projects.ScaffoldResult, error) {
	if input.GitHubToken != "" && c.Remote != nil {
		return c.Remote.Scaffold(ctx, input)
	}
	if c.Local != nil {
		return c.Local.Scaffold(ctx, input)
	}
	return projects.ScaffoldResult{}, errors.New("no scaffolder configured")
}
