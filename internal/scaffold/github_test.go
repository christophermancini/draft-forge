package scaffold

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/yourusername/draft-forge/internal/models"
	"github.com/yourusername/draft-forge/internal/projects"
)

func TestGitHubScaffolderCreatesRepoAndFiles(t *testing.T) {
	var repoCreated bool
	var fileCount int

	client := &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			switch {
			case req.Method == http.MethodPost && strings.Contains(req.URL.Path, "/user/repos"):
				repoCreated = true
				body := `{"html_url":"https://github.com/octo/test-novel","owner":{"login":"octo"}}`
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(strings.NewReader(body)),
					Header:     make(http.Header),
				}, nil
			case req.Method == http.MethodPut && strings.Contains(req.URL.Path, "/contents/"):
				fileCount++
				return &http.Response{
					StatusCode: http.StatusCreated,
					Body:       ioutil.NopCloser(bytes.NewReader([]byte("{}"))),
					Header:     make(http.Header),
				}, nil
			default:
				return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(strings.NewReader("bad request")), Header: make(http.Header)}, nil
			}
		}),
	}

	scaffolder := NewGitHubScaffolder(client)
	scaffolder.APIURL = "https://api.github.com" // not used by fake transport

	_, err := scaffolder.Scaffold(context.Background(), projects.ScaffoldInput{
		Project: models.Project{
			Name:        "Test Novel",
			Slug:        "test-novel",
			ProjectType: "novel",
			Description: "desc",
		},
		GitHubToken: "gh-token",
		Template:    models.TemplateNovel,
	})
	if err != nil {
		t.Fatalf("Scaffold error: %v", err)
	}

	if !repoCreated {
		t.Fatal("expected repo to be created")
	}
	if fileCount == 0 {
		t.Fatal("expected files to be uploaded")
	}
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
