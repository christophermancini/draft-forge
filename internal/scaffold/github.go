package scaffold

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/yourusername/draft-forge/internal/projects"
)

type GitHubScaffolder struct {
	Client *http.Client
	APIURL string
}

func NewGitHubScaffolder(client *http.Client) *GitHubScaffolder {
	if client == nil {
		client = http.DefaultClient
	}
	return &GitHubScaffolder{
		Client: client,
		APIURL: "https://api.github.com",
	}
}

func (g *GitHubScaffolder) Scaffold(ctx context.Context, input projects.ScaffoldInput) (projects.ScaffoldResult, error) {
	if input.GitHubToken == "" {
		return projects.ScaffoldResult{}, fmt.Errorf("github token required")
	}

	repoName := input.Project.Slug
	owner := input.GitHubOwner

	repoInfo, err := g.createRepo(ctx, owner, repoName, input.Project.Description, input.GitHubToken)
	if err != nil {
		return projects.ScaffoldResult{}, fmt.Errorf("create repo: %w", err)
	}
	owner = repoInfo.Owner
	repoURL := repoInfo.URL

	templateRoot := templateRootFor(input.Template)
	files, err := collectTemplateFiles(templateRoot, map[string]any{
		"Name":        input.Project.Name,
		"ProjectType": input.Project.ProjectType,
	})
	if err != nil {
		return projects.ScaffoldResult{}, fmt.Errorf("render templates: %w", err)
	}

	for path, content := range files {
		if err := g.createFile(ctx, owner, repoName, path, content, input.GitHubToken); err != nil {
			return projects.ScaffoldResult{}, fmt.Errorf("create file %s: %w", path, err)
		}
	}

	return projects.ScaffoldResult{RepoURL: repoURL}, nil
}

type repoInfo struct {
	URL   string
	Owner string
}

func (g *GitHubScaffolder) createRepo(ctx context.Context, owner, name, description, token string) (repoInfo, error) {
	type repoResp struct {
		HTMLURL string `json:"html_url"`
		Owner   struct {
			Login string `json:"login"`
		} `json:"owner"`
	}

	payload := map[string]any{
		"name":        name,
		"description": description,
		"private":     true,
	}

	var url string
	if owner != "" {
		url = fmt.Sprintf("%s/orgs/%s/repos", g.APIURL, owner)
	} else {
		url = fmt.Sprintf("%s/user/repos", g.APIURL)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, marshalBody(payload))
	if err != nil {
		return repoInfo{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return repoInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return repoInfo{}, fmt.Errorf("github repo create failed: %s", string(body))
	}

	var parsed repoResp
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return repoInfo{}, err
	}

	if parsed.Owner.Login != "" {
		owner = parsed.Owner.Login
	}
	if parsed.HTMLURL == "" {
		parsed.HTMLURL = fmt.Sprintf("https://github.com/%s/%s", owner, name)
	}
	return repoInfo{URL: parsed.HTMLURL, Owner: owner}, nil
}

func (g *GitHubScaffolder) createFile(ctx context.Context, owner, repo, path string, content []byte, token string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s", g.APIURL, owner, repo, path)
	payload := map[string]any{
		"message": "Initial commit via DraftForge",
		"content": base64.StdEncoding.EncodeToString(content),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, marshalBody(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github create file failed: %s", string(body))
	}
	return nil
}

// collectTemplateFiles renders all embedded files under a template root.
func collectTemplateFiles(root string, data map[string]any) (map[string][]byte, error) {
	out := make(map[string][]byte)

	err := fs.WalkDir(templatesFS, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		content, err := fs.ReadFile(templatesFS, path)
		if err != nil {
			return err
		}
		tmpl, err := template.New(rel).Parse(string(content))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", rel, err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("render template %s: %w", rel, err)
		}
		out[filepath.ToSlash(rel)] = buf.Bytes()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return out, nil
}

func marshalBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}
