package application

import (
	"context"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type ProjectRepository interface {
	ListProjects(ctx context.Context) ([]domain.Project, error)
}

type GetProjectContextRepository interface {
	GetProjectByID(
		ctx context.Context,
		projectID int64,
	) (*domain.Project, error)

	ListOpenTasksByProjectID(
		ctx context.Context,
		projectID int64,
	) ([]domain.Task, error)
}
