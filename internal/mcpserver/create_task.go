package mcpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type CreateTaskInput struct {
	ProjectID   int64  `json:"project_id" jsonschema:"Идентификатор проекта"`
	Title       string `json:"title" jsonschema:"Краткое название задачи"`
	Description string `json:"description,omitempty" jsonschema:"Дополнительное описание задачи"`
	Priority    string `json:"priority" jsonschema:"Приоритет задачи. Допустимые значения: low, medium, high, urgent"`
	Source      string `json:"source" jsonschema:"Источник задачи. Допустимые значения: claude, codex, human"`
}

type CreateTaskOutput struct {
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

func RegisterCreateTask(server *mcp.Server, useCase *application.CreateTask) {
	falseHint := false

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "create_task",
			Description: "Создать задачу",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				OpenWorldHint:   &falseHint,
				DestructiveHint: &falseHint,
			},
		},
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
			input CreateTaskInput,
		) (*mcp.CallToolResult, CreateTaskOutput, error) {
			return createTask(ctx, request, input, useCase)
		},
	)
}

func createTask(
	ctx context.Context,
	request *mcp.CallToolRequest,
	input CreateTaskInput,
	useCase *application.CreateTask,
) (*mcp.CallToolResult, CreateTaskOutput, error) {
	if !domain.ValidateSource(domain.Source(input.Source)) {
		return nil, CreateTaskOutput{}, fmt.Errorf("create task: invalid source")
	}

	if !domain.ValidateTaskPriority(domain.TaskPriority(input.Priority)) {
		return nil, CreateTaskOutput{}, fmt.Errorf("create task: invalid task priority")
	}

	task, err := useCase.Execute(ctx, application.CreateTaskInput{
		ProjectID:   input.ProjectID,
		Title:       input.Title,
		Description: input.Description,
		Priority:    domain.TaskPriority(input.Priority),
		Source:      domain.Source(input.Source),
	})
	if err != nil {
		return nil, CreateTaskOutput{}, err
	}

	return nil, CreateTaskOutput{
		ID:          task.ID,
		ProjectID:   task.ProjectID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		Source:      string(task.Source),
		Version:     task.Version,
		CreatedAt:   task.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}
