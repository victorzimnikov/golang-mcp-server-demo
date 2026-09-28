package application

import (
	"context"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type CreateTask struct {
	taskRepository TaskRepository
}

type CreateTaskInput struct {
	ProjectID   int64
	Title       string
	Description string
	Priority    domain.TaskPriority
	Source      domain.Source
}

func NewCreateTask(taskRepository TaskRepository) *CreateTask {
	return &CreateTask{
		taskRepository: taskRepository,
	}
}

func (c *CreateTask) Execute(
	ctx context.Context,
	input CreateTaskInput,
) (*domain.Task, error) {
	task, err := domain.NewTask(input.ProjectID, input.Title, input.Description, input.Priority, input.Source)
	if err != nil {
		return nil, fmt.Errorf("validate task: %w", err)
	}

	task, err = c.taskRepository.CreateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
