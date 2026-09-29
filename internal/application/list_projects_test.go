package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type fakeProjectRepository struct {
	projects []domain.Project
	err      error
	calls    int
}

func (f *fakeProjectRepository) ListProjects(_ context.Context) ([]domain.Project, error) {
	f.calls++

	return f.projects, f.err
}

func TestListProjectsSuccess(t *testing.T) {
	projects := []domain.Project{
		{
			ID:          1,
			Name:        "First project",
			Description: "First project description",
		},
		{
			ID:          2,
			Name:        "Second project",
			Description: "Second project description",
		},
	}

	repository := &fakeProjectRepository{
		projects: projects,
	}

	useCase := application.NewListProjects(repository)

	result, err := useCase.Execute(t.Context())
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if repository.calls != 1 {
		t.Errorf("calls got %d, want %d", repository.calls, 1)
	}

	if len(result) != len(projects) {
		t.Fatalf("projects len got %d, want %d", len(result), len(projects))
	}

	for idx := range result {
		got := result[idx]
		want := projects[idx]

		if got != want {
			t.Errorf("project index %d got %v, want %v", idx, got, want)
		}
	}
}

func TestListProjectsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository unavailable")

	repository := &fakeProjectRepository{
		err: repositoryErr,
	}

	useCase := application.NewListProjects(repository)

	result, err := useCase.Execute(t.Context())
	if !errors.Is(err, repositoryErr) {
		t.Errorf("execute: %v", err)
	}

	if result != nil {
		t.Errorf("execute result got %v, want nil", result)
	}

	if repository.calls != 1 {
		t.Errorf("calls got %d, want %d", repository.calls, 1)
	}
}
