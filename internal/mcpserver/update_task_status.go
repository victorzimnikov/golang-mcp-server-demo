package mcpserver

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/domain"
)

type UpdateTaskStatusInput struct {
	TaskID          int64  `json:"task_id" jsonschema:"Обязательный идентификатор задачи. Не угадывай идентификатор, если пользователь его не указал"`
	Status          string `json:"status" jsonschema:"Обязательный целевой статус: todo, in_progress, blocked, done или cancelled. Не угадывай статус, если пользователь его не указал"`
	ExpectedVersion int64  `json:"expected_version" jsonschema:"Текущая версия задачи для защиты от конкурентных изменений"`
}

type UpdateTaskStatusOutput struct {
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

func RegisterUpdateTaskStatus(server *mcp.Server, useCase *application.UpdateTaskStatus) {
	falseHint := false
	trueHint := true

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "update_task_status",
			Description: "Обновляет статус задачи в проекте. Используй, когда пользователь просит обновить задачу. Статусы: в работе=in_progress, ожидает=todo, заблокирована=blocked, завершена=done, отменена=cancelled.",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    false,
				OpenWorldHint:   &falseHint,
				DestructiveHint: &trueHint,
				IdempotentHint:  true,
			},
		},
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
			input UpdateTaskStatusInput,
		) (*mcp.CallToolResult, UpdateTaskStatusOutput, error) {
			return updateTaskStatus(ctx, request, input, useCase)
		},
	)
}

func updateTaskStatus(
	ctx context.Context,
	request *mcp.CallToolRequest,
	input UpdateTaskStatusInput,
	useCase *application.UpdateTaskStatus,
) (*mcp.CallToolResult, UpdateTaskStatusOutput, error) {
	source := getSourceFromRequest(request)

	if !domain.ValidateTaskStatus(domain.TaskStatus(input.Status)) {
		return nil, UpdateTaskStatusOutput{}, fmt.Errorf("update task status: invalid status %q", input.Status)
	}

	task, err := useCase.Execute(ctx, application.UpdateTaskStatusInput{
		TaskID:          input.TaskID,
		Status:          domain.TaskStatus(input.Status),
		ExpectedVersion: input.ExpectedVersion,
		Source:          source,
	})
	if err != nil {
		return nil, UpdateTaskStatusOutput{}, err
	}

	return nil, UpdateTaskStatusOutput{
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
