package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type fakeTaskRepository struct {
	receivedTask *domain.Task
	result       *domain.Task
	err          error
	calls        int
}

func (f *fakeTaskRepository) CreateTask(
	_ context.Context,
	task *domain.Task,
) (*domain.Task, error) {
	f.calls++
	f.receivedTask = task

	return f.result, f.err
}

func TestCreateTaskSuccess(t *testing.T) {
	savedTask := &domain.Task{
		ID:          1,
		Title:       "Test title",
		Description: "Test description",
		Priority:    domain.TaskPriorityMedium,
		Source:      domain.SourceCodex,
		Status:      domain.TaskStatusTodo,
		Version:     1,
		ProjectID:   1,
	}

	repository := &fakeTaskRepository{
		result: savedTask,
	}

	useCase := application.NewCreateTask(repository)

	task, err := useCase.Execute(t.Context(), application.CreateTaskInput{
		ProjectID:   1,
		Title:       "   Test title   ",
		Description: "Test description",
		Priority:    domain.TaskPriorityMedium,
		Source:      domain.SourceCodex,
	})
	if err != nil {
		t.Fatalf("create new task: %v", err)
	}

	if task != savedTask {
		t.Errorf("task got '%v', want '%v'", task, savedTask)
	}

	if repository.receivedTask == nil {
		t.Fatalf("new task got 'nil', want '%v'", savedTask)
	}

	if repository.receivedTask.Title != "Test title" {
		t.Errorf("task title got '%s', want '%s'", repository.receivedTask.Title, "Test title")
	}

	if repository.receivedTask.Description != "Test description" {
		t.Errorf("task description got '%s', want '%s'", repository.receivedTask.Description, "Test description")
	}

	if repository.receivedTask.Status != domain.TaskStatusTodo {
		t.Errorf("task status got '%s', want '%s'", repository.receivedTask.Status, domain.TaskStatusTodo)
	}

	if repository.receivedTask.Version != 1 {
		t.Errorf("task version got '%d', want '%d'", repository.receivedTask.Version, 1)
	}

	if repository.receivedTask.ProjectID != 1 {
		t.Errorf("task projectID got '%d', want '%d'", repository.receivedTask.ProjectID, 1)
	}

	if repository.receivedTask.Priority != domain.TaskPriorityMedium {
		t.Errorf("task priority got '%s', want '%s'", repository.receivedTask.Priority, domain.TaskPriorityMedium)
	}

	if repository.receivedTask.Source != domain.SourceCodex {
		t.Errorf("task source got '%s', want '%s'", repository.receivedTask.Source, domain.SourceCodex)
	}

	if repository.calls != 1 {
		t.Errorf("calls got %d, want %d", repository.calls, 1)
	}
}

func TestCreateTaskValidationError(t *testing.T) {
	repository := &fakeTaskRepository{}

	useCase := application.NewCreateTask(repository)

	task, err := useCase.Execute(t.Context(), application.CreateTaskInput{
		ProjectID:   1,
		Title:       "          ",
		Description: "Description",
		Priority:    domain.TaskPriorityHigh,
		Source:      domain.SourceClaude,
	})

	if err == nil {
		t.Fatalf("create new task error must not be nil")
	}

	if task != nil {
		t.Fatalf("task must be nil")
	}

	if repository.calls != 0 {
		t.Errorf("calls got %d, want %d", repository.calls, 0)
	}

	if repository.receivedTask != nil {
		t.Errorf("receivedTask got %v, want nil", repository.receivedTask)
	}
}

func TestCreateTaskRepositoryError(t *testing.T) {
	var repositoryErr = errors.New("repository unavailable")

	repository := &fakeTaskRepository{
		err: repositoryErr,
	}
	useCase := application.NewCreateTask(repository)

	task, err := useCase.Execute(t.Context(), application.CreateTaskInput{
		ProjectID:   1,
		Title:       "Test title",
		Description: "Test description",
		Priority:    domain.TaskPriorityMedium,
		Source:      domain.SourceCodex,
	})
	if !errors.Is(err, repositoryErr) {
		t.Fatalf("execute: %v", err)
	}

	if task != nil {
		t.Fatalf("task must be nil")
	}

	if repository.calls != 1 {
		t.Errorf("calls got %d, want %d", repository.calls, 1)
	}

	if repository.receivedTask == nil {
		t.Errorf("receivedTask must not be nil")
	}
}
