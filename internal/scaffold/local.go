package scaffold

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yourusername/draft-forge/internal/models"
	"github.com/yourusername/draft-forge/internal/projects"
)

// LocalScaffolder writes a project skeleton to the filesystem using embedded templates.
type LocalScaffolder struct {
	Root string
}

func NewLocalScaffolder(root string) *LocalScaffolder {
	return &LocalScaffolder{Root: root}
}

// Scaffold renders the embedded template (currently only "novel") under Root/<slug>.
func (s *LocalScaffolder) Scaffold(ctx context.Context, input projects.ScaffoldInput) (projects.ScaffoldResult, error) {
	_ = ctx // reserved for future cancellation/use

	project := input.Project
	base := filepath.Join(s.Root, project.Slug)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return projects.ScaffoldResult{}, fmt.Errorf("create base dir: %w", err)
	}

	templateRoot := templateRootFor(input.Template)
	data := map[string]any{
		"Name":        project.Name,
		"ProjectType": project.ProjectType,
	}

	err := fs.WalkDir(templatesFS, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == templateRoot {
			return nil
		}

		rel, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(base, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		contents, err := fs.ReadFile(templatesFS, path)
		if err != nil {
			return err
		}

		tmpl, err := template.New(rel).Parse(string(contents))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", rel, err)
		}

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", filepath.Dir(target), err)
		}

		f, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("create file %s: %w", target, err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return fmt.Errorf("render template %s: %w", rel, err)
		}

		return nil
	})
	if err != nil {
		return projects.ScaffoldResult{}, err
	}

	return projects.ScaffoldResult{Path: base}, nil
}

// ListTemplates returns available template names from the embedded FS.
func ListTemplates() ([]string, error) {
	var templates []string
	entries, err := fs.ReadDir(templatesFS, "templates")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			templates = append(templates, strings.TrimSpace(e.Name()))
		}
	}
	return templates, nil
}

func templateRootFor(t models.ProjectTemplate) string {
	if t == "" {
		t = models.TemplateNovel
	}
	return filepath.Join("templates", string(t))
}
