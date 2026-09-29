package mcpserver

import (
	"context"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
)

type ListProjectDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ListProjectsInput struct{}

type ListProjectsOutput struct {
	Projects []ListProjectDTO `json:"projects"`
}

func RegisterListProjects(server *mcp.Server, useCase *application.ListProjects) {
	falseHint := false

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_projects",
			Description: "Показывает список проектов. Используй, когда пользователь просит показать список проектов",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				OpenWorldHint:   &falseHint,
				DestructiveHint: &falseHint,
			},
		},
		func(
			ctx context.Context,
			request *mcp.CallToolRequest,
			input ListProjectsInput,
		) (*mcp.CallToolResult, ListProjectsOutput, error) {
			return listProjects(ctx, request, input, useCase)
		},
	)
}

func listProjects(
	ctx context.Context,
	request *mcp.CallToolRequest,
	input ListProjectsInput,
	useCase *application.ListProjects,
) (*mcp.CallToolResult, ListProjectsOutput, error) {
	listProjects, err := useCase.Execute(ctx)
	if err != nil {
		return nil, ListProjectsOutput{
			Projects: make([]ListProjectDTO, 0),
		}, err
	}

	list := make([]ListProjectDTO, len(listProjects))

	for idx := range listProjects {
		list[idx] = ListProjectDTO{
			ID:          listProjects[idx].ID,
			Name:        listProjects[idx].Name,
			Description: listProjects[idx].Description,
			CreatedAt:   listProjects[idx].CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:   listProjects[idx].UpdatedAt.UTC().Format(time.RFC3339),
		}
	}

	return nil, ListProjectsOutput{
		Projects: list,
	}, nil
}
