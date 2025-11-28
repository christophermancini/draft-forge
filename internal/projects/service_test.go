package projects

import (
	"context"
	"testing"

	"github.com/yourusername/draft-forge/internal/models"
)

func TestCreateGeneratesSlug(t *testing.T) {
	store := &mockStore{}
	scaff := &stubScaffolder{}
	svc := NewService(store, scaff)

	project, scaffoldRes, err := svc.Create(context.Background(), CreateRequest{
		UserID:      1,
		Name:        "My Novel!",
		Description: "desc",
		ProjectType: "novel",
	})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if project.Slug != "my-novel" {
		t.Fatalf("expected slug my-novel, got %s", project.Slug)
	}
	if store.inserted == nil {
		t.Fatal("expected store.InsertProject to be called")
	}
	if scaff.calledWith == nil || scaff.calledWith.Slug != "my-novel" {
		t.Fatal("expected scaffolder to be called with project")
	}
	if scaffoldRes.Path == "" {
		t.Fatal("expected scaffold result path")
	}
}

func TestCreateRejectsBlankName(t *testing.T) {
	store := &mockStore{}
	svc := NewService(store, nil)

	_, _, err := svc.Create(context.Background(), CreateRequest{UserID: 1, Name: " "})
	if err == nil || err != ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}

type mockStore struct {
	inserted *models.Project
}

func (m *mockStore) InsertProject(ctx context.Context, p models.Project) (models.Project, error) {
	p.ID = 1
	m.inserted = &p
	return p, nil
}

func (m *mockStore) UpdateRepoInfo(ctx context.Context, projectID int64, repo models.RepoInfo) error {
	return nil
}

func (m *mockStore) ListProjects(ctx context.Context, userID int64) ([]models.Project, error) {
	return nil, nil
}

type stubScaffolder struct {
	calledWith *models.Project
}

func (s *stubScaffolder) Scaffold(ctx context.Context, input ScaffoldInput) (ScaffoldResult, error) {
	s.calledWith = &input.Project
	return ScaffoldResult{Path: "/tmp/path"}, nil
}
