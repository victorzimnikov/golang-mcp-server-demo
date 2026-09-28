package application

import (
	"context"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type ListProjects struct {
	projectRepository ProjectRepository
}

func NewListProjects(projectRepository ProjectRepository) *ListProjects {
	return &ListProjects{
		projectRepository: projectRepository,
	}
}

func (l *ListProjects) Execute(ctx context.Context) ([]domain.Project, error) {
	projects, err := l.projectRepository.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}

	return projects, nil
}
