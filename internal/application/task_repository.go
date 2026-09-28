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
