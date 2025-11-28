package project

import (
	"database/sql"

	"github.com/yourusername/draft-forge/internal/models"
)

// dbProject holds DB-facing fields including nullable columns.
type dbProject struct {
	ID             int64          `db:"id"`
	UserID         int64          `db:"user_id"`
	Name           string         `db:"name"`
	Slug           string         `db:"slug"`
	Description    sql.NullString `db:"description"`
	ProjectType    string         `db:"project_type"`
	GitHubRepoID   sql.NullInt64  `db:"github_repo_id"`
	GitHubRepoName sql.NullString `db:"github_repo_name"`
	GitHubRepoURL  sql.NullString `db:"github_repo_url"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
}

func toDBModel(p models.Project) dbProject {
	var desc sql.NullString
	if p.Description != "" {
		desc = sql.NullString{String: p.Description, Valid: true}
	}

	var repoID sql.NullInt64
	var repoName sql.NullString
	var repoURL sql.NullString
	if p.GitHubRepo != nil {
		if p.GitHubRepo.ID != nil {
			repoID = sql.NullInt64{Int64: *p.GitHubRepo.ID, Valid: true}
		}
		if p.GitHubRepo.Name != "" {
			repoName = sql.NullString{String: p.GitHubRepo.Name, Valid: true}
		}
		if p.GitHubRepo.URL != "" {
			repoURL = sql.NullString{String: p.GitHubRepo.URL, Valid: true}
		}
	}

	return dbProject{
		ID:             p.ID,
		UserID:         p.UserID,
		Name:           p.Name,
		Slug:           p.Slug,
		Description:    desc,
		ProjectType:    p.ProjectType,
		GitHubRepoID:   repoID,
		GitHubRepoName: repoName,
		GitHubRepoURL:  repoURL,
	}
}

func (dbp dbProject) toModel() models.Project {
	var repo *models.RepoInfo
	if dbp.GitHubRepoID.Valid || dbp.GitHubRepoName.Valid || dbp.GitHubRepoURL.Valid {
		repo = &models.RepoInfo{}
		if dbp.GitHubRepoID.Valid {
			val := dbp.GitHubRepoID.Int64
			repo.ID = &val
		}
		if dbp.GitHubRepoName.Valid {
			repo.Name = dbp.GitHubRepoName.String
		}
		if dbp.GitHubRepoURL.Valid {
			repo.URL = dbp.GitHubRepoURL.String
		}
	}

	p := models.Project{
		ID:          dbp.ID,
		UserID:      dbp.UserID,
		Name:        dbp.Name,
		Slug:        dbp.Slug,
		Description: dbp.Description.String,
		ProjectType: dbp.ProjectType,
		GitHubRepo:  repo,
	}
	return p
}
