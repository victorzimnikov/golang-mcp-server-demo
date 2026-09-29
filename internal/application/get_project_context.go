package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

var (
	ErrProjectIDMustBeGreaterThanZero = errors.New("project ID must be greater than 0")
	ErrProjectNotFound                = errors.New("project not found")
)

type GetProjectContextInput struct {
	ProjectID int64
}

type ProjectContext struct {
	Project   domain.Project
	OpenTasks []domain.Task
}

type GetProjectContext struct {
	repository GetProjectContextRepository
}

func NewGetProjectContext(repository GetProjectContextRepository) *GetProjectContext {
	return &GetProjectContext{
		repository: repository,
	}
}

func (r *GetProjectContext) Execute(ctx context.Context, input GetProjectContextInput) (*ProjectContext, error) {
	if input.ProjectID <= 0 {
		return nil, fmt.Errorf("execute: %w", ErrProjectIDMustBeGreaterThanZero)
	}

	project, err := r.repository.GetProjectByID(ctx, input.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("execute get project: %w", err)
	}

	openTasks, err := r.repository.ListOpenTasksByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("execute list open tasks: %w", err)
	}

	if openTasks == nil {
		openTasks = make([]domain.Task, 0)
	}

	return &ProjectContext{
		Project:   *project,
		OpenTasks: openTasks,
	}, nil
}
