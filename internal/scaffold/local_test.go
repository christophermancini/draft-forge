package scaffold

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/draft-forge/internal/models"
	"github.com/yourusername/draft-forge/internal/projects"
)

func TestLocalScaffolderCreatesFiles(t *testing.T) {
	root := t.TempDir()
	scaffolder := NewLocalScaffolder(root)

	project := models.Project{Name: "Test Novel", Slug: "test-novel", ProjectType: "novel"}
	res, err := scaffolder.Scaffold(context.Background(), projects.ScaffoldInput{Project: project})
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}

	readme := filepath.Join(res.Path, "README.md")
	if _, err := os.Stat(readme); err != nil {
		t.Fatalf("expected README.md to exist: %v", err)
	}

	metadata := filepath.Join(res.Path, "manuscript", "metadata.yaml")
	b, err := os.ReadFile(metadata)
	if err != nil {
		t.Fatalf("read metadata: %v", err)
	}
	if !strings.Contains(string(b), "Test Novel") {
		t.Fatalf("expected metadata to include project name, got: %s", string(b))
	}
}
