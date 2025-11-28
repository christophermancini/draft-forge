package project

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

func (s *Store) InsertProject(ctx context.Context, p models.Project) (models.Project, error) {
	dbp := toDBModel(p)

	query := `
		INSERT INTO projects (user_id, name, slug, description, project_type, github_repo_id, github_repo_name, github_repo_url)
		VALUES (:user_id, :name, :slug, NULLIF(:description, ''), :project_type, :github_repo_id, :github_repo_name, :github_repo_url)
		RETURNING id, COALESCE(description, ''), github_repo_id, github_repo_name, github_repo_url, created_at, updated_at
	`

	rows, err := s.db.NamedQueryContext(ctx, query, dbp)
	if err != nil {
		return models.Project{}, fmt.Errorf("insert project: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&dbp.ID,
			&dbp.Description,
			&dbp.GitHubRepoID,
			&dbp.GitHubRepoName,
			&dbp.GitHubRepoURL,
			&dbp.CreatedAt,
			&dbp.UpdatedAt,
		); err != nil {
			return models.Project{}, fmt.Errorf("scan project: %w", err)
		}
	}

	return dbp.toModel(), nil
}

func (s *Store) UpdateRepoInfo(ctx context.Context, projectID int64, repo models.RepoInfo) error {
	dbp := toDBModel(models.Project{GitHubRepo: &repo})

	query := `
		UPDATE projects
		SET github_repo_id = :github_repo_id,
		    github_repo_name = :github_repo_name,
		    github_repo_url = :github_repo_url,
		    updated_at = NOW()
		WHERE id = :id
	`

	params := map[string]interface{}{
		"id":               projectID,
		"github_repo_id":   dbp.GitHubRepoID,
		"github_repo_name": dbp.GitHubRepoName,
		"github_repo_url":  dbp.GitHubRepoURL,
	}

	if _, err := s.db.NamedExecContext(ctx, query, params); err != nil {
		return fmt.Errorf("update repo info: %w", err)
	}
	return nil
}

func (s *Store) ListProjects(ctx context.Context, userID int64) ([]models.Project, error) {
	query := `
		SELECT id, user_id, name, slug, COALESCE(description, '') AS description, project_type,
		       github_repo_id, github_repo_name, github_repo_url, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var dbProjects []dbProject
	if err := s.db.SelectContext(ctx, &dbProjects, query, userID); err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	projects := make([]models.Project, 0, len(dbProjects))
	for _, dbp := range dbProjects {
		projects = append(projects, dbp.toModel())
	}
	return projects, nil
}
