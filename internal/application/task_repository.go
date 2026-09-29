package application

import (
	"context"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type TaskRepository interface {
	CreateTask(
		ctx context.Context,
		task *domain.Task,
	) (*domain.Task, error)
}

type UpdateTaskStatusRepository interface {
	GetTaskByID(ctx context.Context, taskID int64) (*domain.Task, error)

	UpdateTaskStatus(
		ctx context.Context,
		taskID int64,
		status domain.TaskStatus,
		expectedVersion int64,
		source domain.Source,
	) (*domain.Task, error)
}
