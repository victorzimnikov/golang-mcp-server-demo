package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type fakeGetProjectContextRepository struct {
	project                *domain.Project
	openTasks              []domain.Task
	getProjectErr          error
	openTasksErr           error
	getProjectCalls        int
	listOpenTasksCalls     int
	getProjectProjectID    int64
	listOpenTasksProjectID int64
}

func (f *fakeGetProjectContextRepository) GetProjectByID(
	ctx context.Context,
	projectID int64,
) (*domain.Project, error) {
	f.getProjectCalls++
	f.getProjectProjectID = projectID

	return f.project, f.getProjectErr
}

func (f *fakeGetProjectContextRepository) ListOpenTasksByProjectID(
	ctx context.Context,
	projectID int64,
) ([]domain.Task, error) {
	f.listOpenTasksCalls++
	f.listOpenTasksProjectID = projectID

	return f.openTasks, f.openTasksErr
}

func TestGetProjectContextValidationError(t *testing.T) {
	repository := &fakeGetProjectContextRepository{}

	useCase := application.NewGetProjectContext(repository)

	result, err := useCase.Execute(t.Context(), application.GetProjectContextInput{
		ProjectID: 0,
	})
	if !errors.Is(err, application.ErrProjectIDMustBeGreaterThanZero) {
		t.Errorf("execute: %v", err)
	}

	if result != nil {
		t.Errorf("execute result got %v, want nil", result)
	}

	if repository.getProjectCalls != 0 {
		t.Errorf("get project calls got %d, want %d", repository.getProjectCalls, 0)
	}

	if repository.listOpenTasksCalls != 0 {
		t.Errorf("list open tasks calls got %d, want %d", repository.listOpenTasksCalls, 0)
	}
}

func TestGetProjectContextSuccess(t *testing.T) {
	project := domain.Project{
		ID:          1,
		Name:        "Test project",
		Description: "Test project description",
	}
	openTasks := []domain.Task{
		{
			ID:          1,
			ProjectID:   1,
			Title:       "First task",
			Description: "First task description",
			Priority:    domain.TaskPriorityLow,
			Source:      domain.SourceClaude,
		},
		{
			ID:          2,
			ProjectID:   1,
			Title:       "Second task",
			Description: "Second task description",
			Priority:    domain.TaskPriorityHigh,
			Source:      domain.SourceCodex,
		},
	}

	repository := &fakeGetProjectContextRepository{
		project:   &project,
		openTasks: openTasks,
	}

	useCase := application.NewGetProjectContext(repository)

	result, err := useCase.Execute(t.Context(), application.GetProjectContextInput{
		ProjectID: 1,
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if result == nil {
		t.Fatal("get result got nil, want application.ProjectContext")
	}

	if repository.getProjectCalls != 1 {
		t.Errorf("get project calls got %d, want %d", repository.getProjectCalls, 1)
	}

	if repository.listOpenTasksCalls != 1 {
		t.Errorf("list open tasks calls got %d, want %d", repository.listOpenTasksCalls, 1)
	}

	if result.Project != project {
		t.Errorf("result project got %v, want %v", result.Project, project)
	}

	if len(result.OpenTasks) != len(openTasks) {
		t.Fatalf("result open tasks len got %d, want %d", len(result.OpenTasks), len(openTasks))
	}

	for idx := range result.OpenTasks {
		got := result.OpenTasks[idx]
		want := openTasks[idx]

		if got != want {
			t.Errorf("open task index %d got %v, want %v", idx, got, want)
		}
	}

	if repository.getProjectProjectID != 1 {
		t.Errorf("get project projectID got %d, want %d", repository.getProjectProjectID, 1)
	}

	if repository.listOpenTasksProjectID != 1 {
		t.Errorf("list open tasks projectID got %d, want %d", repository.listOpenTasksProjectID, 1)
	}
}

func TestGetProjectContextProjectRepositoryError(t *testing.T) {
	projectErr := errors.New("get project unavailable")

	repository := &fakeGetProjectContextRepository{
		getProjectErr: projectErr,
	}

	useCase := application.NewGetProjectContext(repository)

	result, err := useCase.Execute(t.Context(), application.GetProjectContextInput{
		ProjectID: 1,
	})
	if !errors.Is(err, projectErr) {
		t.Errorf("execute %v", err)
	}

	if result != nil {
		t.Errorf("result got %v, want nil", result)
	}

	if repository.getProjectProjectID != 1 {
		t.Errorf("get project projectID got %d, want %d", repository.getProjectProjectID, 1)
	}

	if repository.listOpenTasksProjectID != 0 {
		t.Errorf("list open tasks projectID got %d, want %d", repository.listOpenTasksProjectID, 0)
	}

	if repository.getProjectCalls != 1 {
		t.Errorf("get project calls got %d, want %d", repository.getProjectCalls, 1)
	}

	if repository.listOpenTasksCalls != 0 {
		t.Errorf("list open tasks calls got %d, want %d", repository.listOpenTasksCalls, 0)
	}
}

func TestGetProjectContextOpenTasksRepositoryError(t *testing.T) {
	openTasksErr := errors.New("list open tasks unavailable")

	project := &domain.Project{ID: 1}

	repository := &fakeGetProjectContextRepository{
		project:      project,
		openTasksErr: openTasksErr,
	}

	useCase := application.NewGetProjectContext(repository)

	result, err := useCase.Execute(t.Context(), application.GetProjectContextInput{
		ProjectID: 1,
	})
	if !errors.Is(err, openTasksErr) {
		t.Errorf("execute %v", err)
	}

	if result != nil {
		t.Errorf("result got %v, want nil", result)
	}

	if repository.getProjectProjectID != 1 {
		t.Errorf("get project projectID got %d, want %d", repository.getProjectProjectID, 1)
	}

	if repository.listOpenTasksProjectID != 1 {
		t.Errorf("list open tasks projectID got %d, want %d", repository.listOpenTasksProjectID, 1)
	}

	if repository.getProjectCalls != 1 {
		t.Errorf("get project calls got %d, want %d", repository.getProjectCalls, 1)
	}

	if repository.listOpenTasksCalls != 1 {
		t.Errorf("list open tasks calls got %d, want %d", repository.listOpenTasksCalls, 1)
	}
}

func TestGetProjectContextNormalizesNilOpenTasks(t *testing.T) {
	project := &domain.Project{ID: 1}

	repository := &fakeGetProjectContextRepository{
		project: project,
	}

	useCase := application.NewGetProjectContext(repository)

	result, err := useCase.Execute(t.Context(), application.GetProjectContextInput{
		ProjectID: 1,
	})
	if err != nil {
		t.Fatalf("execute %v", err)
	}

	if result == nil {
		t.Fatal("result got nil, want application.ProjectContext{}")
	}

	if result.OpenTasks == nil {
		t.Fatal("list open tasks got nil, want []domain.Task{}")
	}

	if len(result.OpenTasks) != 0 {
		t.Errorf("list open tasks len got %d, want %d", len(result.OpenTasks), 0)
	}
}
