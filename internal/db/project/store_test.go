package project

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/yourusername/draft-forge/internal/models"
)

func TestInsertProject(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	store := NewStore(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO projects (user_id, name, slug, description, project_type, github_repo_id, github_repo_name, github_repo_url)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $7, $8)
		RETURNING id, COALESCE(description, ''), github_repo_id, github_repo_name, github_repo_url, created_at, updated_at
	`)).
		WithArgs(int64(1), "Name", "name", "desc", "novel", nil, nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id", "description", "github_repo_id", "github_repo_name", "github_repo_url", "created_at", "updated_at"}).
			AddRow(int64(10), "desc", nil, nil, nil, time.Now(), time.Now()))

	project, err := store.InsertProject(context.Background(), models.Project{
		UserID:      1,
		Name:        "Name",
		Slug:        "name",
		Description: "desc",
		ProjectType: "novel",
	})
	if err != nil {
		t.Fatalf("InsertProject error: %v", err)
	}
	if project.ID != 10 || project.Description != "desc" {
		t.Fatalf("unexpected project %+v", project)
	}
}

func TestListProjects(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	store := NewStore(sqlxDB)
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, user_id, name, slug, COALESCE(description, '') AS description, project_type,
		       github_repo_id, github_repo_name, github_repo_url, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "slug", "description", "project_type", "github_repo_id", "github_repo_name", "github_repo_url", "created_at", "updated_at"}).
			AddRow(int64(1), int64(1), "Name", "name", "desc", "novel", int64(999), "octo/name", "https://github.com/octo/name", now, now))

	projects, err := store.ListProjects(context.Background(), 1)
	if err != nil {
		t.Fatalf("ListProjects error: %v", err)
	}
	if len(projects) != 1 || projects[0].Slug != "name" {
		t.Fatalf("unexpected projects %+v", projects)
	}
	if projects[0].GitHubRepo == nil || projects[0].GitHubRepo.URL == "" {
		t.Fatalf("expected github repo info, got %+v", projects[0].GitHubRepo)
	}
}
