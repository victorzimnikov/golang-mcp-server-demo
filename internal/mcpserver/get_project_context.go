package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
)

type GetProjectContextInput struct {
	ProjectID int64 `json:"project_id" jsonschema:"Идентификатор проекта"`
}

type ProjectContextProjectDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProjectContextTaskDTO struct {
	ID          int64  `json:"id"`
	ProjectID   int64  `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	Source      string `json:"source"`
	Version     int64  `json:"version"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type GetProjectContextOutput struct {
	Project   ProjectContextProjectDTO `json:"project"`
	OpenTasks []ProjectContextTaskDTO  `json:"open_tasks"`
}

func RegisterGetProjectContext(server *mcp.Server, useCase *application.GetProjectContext) {
	falseHint := false

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_project_context",
			Description: "Возвращает проект и открытые задачи с их текущими версиями. Используй перед update_task_status, чтобы получить expected_version.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				OpenWorldHint:   &falseHint,
				DestructiveHint: &falseHint,
			},
		},
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
			input GetProjectContextInput,
		) (*mcp.CallToolResult, GetProjectContextOutput, error) {
			return getProjectContext(ctx, request, input, useCase)
		},
	)
}

func getProjectContext(ctx context.Context, request *mcp.CallToolRequest, input GetProjectContextInput, useCase *application.GetProjectContext) (*mcp.CallToolResult, GetProjectContextOutput, error) {
	res, err := useCase.Execute(ctx, application.GetProjectContextInput{
		ProjectID: input.ProjectID,
	})
	if err != nil {
		return nil, GetProjectContextOutput{}, err
	}

	openTasks := make([]ProjectContextTaskDTO, len(res.OpenTasks))

	for idx, item := range res.OpenTasks {
		openTasks[idx] = ProjectContextTaskDTO{
			ID:          item.ID,
			ProjectID:   item.ProjectID,
			Title:       item.Title,
			Description: item.Description,
			Status:      string(item.Status),
			Priority:    string(item.Priority),
			Source:      string(item.Source),
			Version:     item.Version,
			CreatedAt:   item.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   item.UpdatedAt.UTC().Format(time.RFC3339),
		}
	}

	return nil, GetProjectContextOutput{
		Project: ProjectContextProjectDTO{
			ID:          res.Project.ID,
			Name:        res.Project.Name,
			Description: res.Project.Description,
			CreatedAt:   res.Project.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   res.Project.UpdatedAt.UTC().Format(time.RFC3339),
		},
		OpenTasks: openTasks,
	}, nil
}
