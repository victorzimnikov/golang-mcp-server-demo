package application

import (
	"context"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type ProjectRepository interface {
	ListProjects(ctx context.Context) ([]domain.Project, error)
}
