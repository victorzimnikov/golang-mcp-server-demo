package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	ServerName    = "golang-mcp-server-demo"
	ServerVersion = "v0.1.0"
)

type ServerInfoInput struct{}

type ServerInfoOutput struct {
	Name    string `json:"name" jsonschema:"имя сервера"`
	Version string `json:"version" jsonschema:"версия сервера"`
}

func RegisterServerInfo(server *mcp.Server) {
	falseHint := false

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "server_info",
			Description: "Показывает информация о сервере. Используй, когда пользователь просит показать информацию о сервере",
			Annotations: &mcp.ToolAnnotations{
				ReadOnlyHint:    true,
				OpenWorldHint:   &falseHint,
				DestructiveHint: &falseHint,
			},
		},
		serverInfo,
	)
}

func serverInfo(ctx context.Context, request *mcp.CallToolRequest, input ServerInfoInput) (*mcp.CallToolResult, ServerInfoOutput, error) {
	return nil, ServerInfoOutput{
		Name:    ServerName,
		Version: ServerVersion,
	}, nil
}
