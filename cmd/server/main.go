package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/mcpserver"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "golang-mcp-server-demo",
			Version: "v0.1.0",
		},
		nil,
	)

	mcpserver.RegisterTool(server)

	return server.Run(context.Background(), &mcp.StdioTransport{})
}
