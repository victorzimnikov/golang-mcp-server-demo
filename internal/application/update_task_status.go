package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

var (
	ErrTaskVersionConflict                  = errors.New("task version conflict")
	ErrTaskSourceInvalid                    = errors.New("task source invalid")
	ErrTaskIDMustBeGreaterThanZero          = errors.New("task ID must be greater than 0")
	ErrExpectedVersionMustBeGreaterThanZero = errors.New("expected version must be greater than 0")
	ErrTaskNotFound                         = errors.New("task not found")
)

type UpdateTaskStatusInput struct {
	TaskID          int64
	Status          domain.TaskStatus
	ExpectedVersion int64
	Source          domain.Source
}

type UpdateTaskStatus struct {
	repository UpdateTaskStatusRepository
}

func NewUpdateTaskStatus(repository UpdateTaskStatusRepository) *UpdateTaskStatus {
	return &UpdateTaskStatus{
		repository: repository,
	}
}

func (t *UpdateTaskStatus) Execute(ctx context.Context, input UpdateTaskStatusInput) (*domain.Task, error) {
	if input.TaskID <= 0 {
		return nil, ErrTaskIDMustBeGreaterThanZero
	}

	if input.ExpectedVersion <= 0 {
		return nil, ErrExpectedVersionMustBeGreaterThanZero
	}

	if !domain.ValidateSource(input.Source) {
		return nil, ErrTaskSourceInvalid
	}

	task, err := t.repository.GetTaskByID(ctx, input.TaskID)
	if err != nil {
		return nil, fmt.Errorf("execute get task: %w", err)
	}

	if task.Version != input.ExpectedVersion {
		return nil, ErrTaskVersionConflict
	}

	if err := task.ChangeStatus(input.Status); err != nil {
		return nil, fmt.Errorf("execute change status: %w", err)
	}

	newTask, err := t.repository.UpdateTaskStatus(ctx, input.TaskID, task.Status, input.ExpectedVersion, input.Source)
	if err != nil {
		return nil, fmt.Errorf("execute update status: %w", err)
	}

	return newTask, nil
}
