package projects

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/yourusername/draft-forge/internal/models"
)

var (
	ErrInvalidName = errors.New("invalid project name")
)

type CreateRequest struct {
	UserID      int64
	Name        string
	Description string
	ProjectType string
	GitHubToken string
	GitHubOwner string
	Template    models.ProjectTemplate
}

type Store interface {
	InsertProject(ctx context.Context, p models.Project) (models.Project, error)
	UpdateRepoInfo(ctx context.Context, projectID int64, repo models.RepoInfo) error
	ListProjects(ctx context.Context, userID int64) ([]models.Project, error)
}

// Scaffolder creates the initial project structure (local or GitHub-backed).
type Service struct {
	store      Store
	scaffolder Scaffolder
}

func NewService(store Store, scaffolder Scaffolder) *Service {
	return &Service{store: store, scaffolder: scaffolder}
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (models.Project, ScaffoldResult, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return models.Project{}, ScaffoldResult{}, ErrInvalidName
	}

	slug := slugify(name)
	if slug == "" {
		return models.Project{}, ScaffoldResult{}, ErrInvalidName
	}

	project := models.Project{
		UserID:      req.UserID,
		Name:        name,
		Slug:        slug,
		Description: strings.TrimSpace(req.Description),
		ProjectType: strings.TrimSpace(req.ProjectType),
	}

	created, err := s.store.InsertProject(ctx, project)
	if err != nil {
		return models.Project{}, ScaffoldResult{}, fmt.Errorf("insert project: %w", err)
	}

	var scaffoldResult ScaffoldResult
	if s.scaffolder != nil {
		result, err := s.scaffolder.Scaffold(ctx, ScaffoldInput{
			Project:     created,
			GitHubToken: strings.TrimSpace(req.GitHubToken),
			GitHubOwner: strings.TrimSpace(req.GitHubOwner),
			Template:    req.Template,
		})
		if err != nil {
			return models.Project{}, ScaffoldResult{}, fmt.Errorf("scaffold project: %w", err)
		}
		scaffoldResult = result

		if result.RepoURL != "" {
			repoInfo := models.RepoInfo{
				URL:  result.RepoURL,
				Name: created.Slug,
			}
			if err := s.store.UpdateRepoInfo(ctx, created.ID, repoInfo); err != nil {
				return models.Project{}, ScaffoldResult{}, fmt.Errorf("update repo info: %w", err)
			}
			created.GitHubRepo = &repoInfo
		}
	}

	return created, scaffoldResult, nil
}

func (s *Service) List(ctx context.Context, userID int64) ([]models.Project, error) {
	projects, err := s.store.ListProjects(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

var slugPattern = regexp.MustCompile(`[^a-z0-9-]+`)

func slugify(input string) string {
	lower := strings.ToLower(strings.TrimSpace(input))
	replaced := strings.ReplaceAll(lower, " ", "-")
	replaced = slugPattern.ReplaceAllString(replaced, "")
	replaced = strings.Trim(replaced, "-")
	return replaced
}
