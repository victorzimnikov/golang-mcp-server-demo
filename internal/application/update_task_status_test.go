package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type fakeUpdateTaskStatusRepository struct {
	task               *domain.Task
	updatedTask        *domain.Task
	status             domain.TaskStatus
	expectedVersion    int64
	source             domain.Source
	getTaskCalls       int
	updateStatusCalls  int
	getTaskErr         error
	updateStatusErr    error
	getTaskTaskID      int64
	updateStatusTaskID int64
}

func (f *fakeUpdateTaskStatusRepository) GetTaskByID(
	_ context.Context,
	taskID int64,
) (*domain.Task, error) {
	f.getTaskCalls++
	f.getTaskTaskID = taskID

	return f.task, f.getTaskErr
}

func (f *fakeUpdateTaskStatusRepository) UpdateTaskStatus(
	_ context.Context,
	taskID int64,
	status domain.TaskStatus,
	expectedVersion int64,
	source domain.Source,
) (*domain.Task, error) {
	f.updateStatusCalls++
	f.updateStatusTaskID = taskID
	f.status = status
	f.expectedVersion = expectedVersion
	f.source = source

	return f.updatedTask, f.updateStatusErr
}

func TestUpdateTaskStatusInvalidTaskID(t *testing.T) {
	repository := &fakeUpdateTaskStatusRepository{}

	useCase := application.NewUpdateTaskStatus(repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          0,
		Status:          domain.TaskStatusInProgress,
		Source:          domain.SourceCodex,
		ExpectedVersion: 1,
	})
	if !errors.Is(err, application.ErrTaskIDMustBeGreaterThanZero) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 0 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 0)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}
}

func TestUpdateTaskStatusInvalidExpectedVersion(t *testing.T) {
	repository := fakeUpdateTaskStatusRepository{}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		ExpectedVersion: 0,
		Status:          domain.TaskStatusInProgress,
		Source:          domain.SourceCodex,
	})
	if !errors.Is(err, application.ErrExpectedVersionMustBeGreaterThanZero) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 0 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 0)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}
}

func TestUpdateTaskStatusInvalidSource(t *testing.T) {
	repository := fakeUpdateTaskStatusRepository{}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		ExpectedVersion: 1,
		Status:          domain.TaskStatusInProgress,
		Source:          domain.Source("unknown"),
	})
	if !errors.Is(err, application.ErrTaskSourceInvalid) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 0 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 0)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}
}

func TestUpdateTaskStatusInvalidStatus(t *testing.T) {
	savedTask := &domain.Task{
		ID:          1,
		ProjectID:   1,
		Title:       "Title",
		Description: "Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityLow,
		Source:      domain.SourceClaude,
		Version:     1,
	}

	repository := fakeUpdateTaskStatusRepository{task: savedTask}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		ExpectedVersion: 1,
		Status:          domain.TaskStatus("unknown"),
		Source:          domain.SourceCodex,
	})
	if !errors.Is(err, domain.ErrInvalidTaskStatusTransition) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 1 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 1)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}

	if repository.getTaskTaskID != 1 {
		t.Errorf("get task taskID got %d, want %d", repository.getTaskTaskID, 1)
	}

	if savedTask.Status != domain.TaskStatusTodo {
		t.Errorf("task status got %s, want %s", savedTask.Status, domain.TaskStatusTodo)
	}
}

func TestUpdateTaskStatusGetTaskRepositoryError(t *testing.T) {
	getTaskErr := errors.New("get task error")

	repository := fakeUpdateTaskStatusRepository{getTaskErr: getTaskErr}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		Status:          domain.TaskStatusInProgress,
		ExpectedVersion: 1,
		Source:          domain.SourceClaude,
	})
	if !errors.Is(err, getTaskErr) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 1 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 1)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}

	if repository.getTaskTaskID != 1 {
		t.Errorf("get task taskID got %d, want %d", repository.getTaskTaskID, 1)
	}
}

func TestUpdateTaskStatusVersionConflict(t *testing.T) {
	savedTask := &domain.Task{
		ID:          1,
		ProjectID:   1,
		Title:       "Title",
		Description: "Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityLow,
		Source:      domain.SourceClaude,
		Version:     2,
	}

	repository := fakeUpdateTaskStatusRepository{task: savedTask}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		ExpectedVersion: 1,
		Status:          domain.TaskStatusInProgress,
		Source:          domain.SourceCodex,
	})

	if !errors.Is(err, application.ErrTaskVersionConflict) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 1 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 1)
	}

	if repository.updateStatusCalls != 0 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 0)
	}

	if repository.getTaskTaskID != 1 {
		t.Errorf("get task taskID got %d, want %d", repository.getTaskTaskID, 1)
	}

	if savedTask.Status != domain.TaskStatusTodo {
		t.Errorf("get task status got %s, want %s", savedTask.Status, domain.TaskStatusTodo)
	}
}

func TestUpdateTaskStatusSuccess(t *testing.T) {
	savedTask := &domain.Task{
		ID:          1,
		ProjectID:   1,
		Title:       "Title",
		Description: "Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityLow,
		Source:      domain.SourceClaude,
		Version:     1,
	}
	updatedTask := &domain.Task{
		ID:          1,
		ProjectID:   1,
		Title:       "Title",
		Description: "Description",
		Status:      domain.TaskStatusInProgress,
		Priority:    domain.TaskPriorityLow,
		Source:      domain.SourceClaude,
		Version:     2,
	}

	repository := fakeUpdateTaskStatusRepository{task: savedTask, updatedTask: updatedTask}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		ExpectedVersion: 1,
		Status:          domain.TaskStatusInProgress,
		Source:          domain.SourceCodex,
	})
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	if task != updatedTask {
		t.Fatalf("task got %v, want %v", task, updatedTask)
	}

	if repository.getTaskCalls != 1 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 1)
	}

	if repository.updateStatusCalls != 1 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 1)
	}

	if repository.getTaskTaskID != 1 {
		t.Errorf("get task taskID got %d, want %d", repository.getTaskTaskID, 1)
	}

	if repository.updateStatusTaskID != 1 {
		t.Errorf("get update status task taskID got %d, want %d", repository.updateStatusTaskID, 1)
	}

	if savedTask.Status != domain.TaskStatusInProgress {
		t.Errorf("updated status task got %s, want %s", savedTask.Status, domain.TaskStatusInProgress)
	}

	if repository.status != domain.TaskStatusInProgress {
		t.Errorf("repository status got %s, want %s", repository.status, domain.TaskStatusInProgress)
	}

	if repository.source != domain.SourceCodex {
		t.Errorf("repository source got %s, want %s", repository.source, domain.SourceCodex)
	}

	if repository.expectedVersion != 1 {
		t.Errorf("repository expectedVersion got %d, want %d", repository.expectedVersion, 1)
	}
}

func TestUpdateTaskStatusUpdateRepositoryError(t *testing.T) {
	savedTask := &domain.Task{
		ID:          1,
		ProjectID:   1,
		Title:       "Title",
		Description: "Description",
		Status:      domain.TaskStatusTodo,
		Priority:    domain.TaskPriorityLow,
		Source:      domain.SourceClaude,
		Version:     1,
	}
	updateStatusErr := errors.New("update task status error")

	repository := fakeUpdateTaskStatusRepository{updateStatusErr: updateStatusErr, task: savedTask}

	useCase := application.NewUpdateTaskStatus(&repository)

	task, err := useCase.Execute(t.Context(), application.UpdateTaskStatusInput{
		TaskID:          1,
		Status:          domain.TaskStatusInProgress,
		ExpectedVersion: 1,
		Source:          domain.SourceClaude,
	})
	if !errors.Is(err, updateStatusErr) {
		t.Errorf("execute: %v", err)
	}

	if task != nil {
		t.Errorf("execute task must be nil")
	}

	if repository.getTaskCalls != 1 {
		t.Errorf("get task calls got %d, want %d", repository.getTaskCalls, 1)
	}

	if repository.updateStatusCalls != 1 {
		t.Errorf("update task status calls got %d, want %d", repository.updateStatusCalls, 1)
	}

	if repository.updateStatusTaskID != 1 {
		t.Errorf("repository updateStatusTaskID got %d, want %d", repository.updateStatusTaskID, 1)
	}

	if repository.status != domain.TaskStatusInProgress {
		t.Errorf("repository status got %s, want %s", repository.status, domain.TaskStatusInProgress)
	}

	if repository.source != domain.SourceClaude {
		t.Errorf("repository source got %s, want %s", repository.source, domain.SourceClaude)
	}

	if repository.expectedVersion != 1 {
		t.Errorf("repository expectedVersion got %d, want %d", repository.expectedVersion, 1)
	}
}
