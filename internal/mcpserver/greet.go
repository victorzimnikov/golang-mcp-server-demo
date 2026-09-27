package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GreetInput struct {
	Name string `json:"name" jsonschema:"имя пользователя"`
}

type GreetOutput struct {
	Greeting string `json:"greeting" jsonschema:"текст приветствия"`
}

func RegisterTool(server *mcp.Server) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "greet",
			Description: "Приветствует пользователя по имени",
		},
		greet,
	)
}

func greet(ctx context.Context, request *mcp.CallToolRequest, input GreetInput) (*mcp.CallToolResult, GreetOutput, error) {
	return nil, GreetOutput{
		Greeting: "Привет, " + input.Name + "!",
	}, nil
}
