package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mux := http.NewServeMux()

	client := mcp.NewClient(
		&mcp.Implementation{
			Name:    "project-context-hub-ai-service",
			Version: "v0.1.0",
		},
		nil,
	)

	transport := &mcp.StreamableClientTransport{
		Endpoint:             "http://127.0.0.1:8000/mcp",
		DisableStandaloneSSE: true,
	}

	session, err := client.Connect(context.Background(), transport, nil)
	if err != nil {
		return fmt.Errorf("client connect: %w", err)
	}
	defer session.Close()

	listToolsResult, err := session.ListTools(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("list tools: %w", err)
	}

	log.Printf("Tools count: %d\n", len(listToolsResult.Tools))

	for idx, tool := range listToolsResult.Tools {
		log.Printf("%d. %s\n", idx+1, tool.Name)
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		fmt.Fprint(w, "ok")
	})

	return http.ListenAndServe(":8001", mux)
}
