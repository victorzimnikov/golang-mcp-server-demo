package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/application"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/mcpserver"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/storage/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dsn := "example.db?_foreign_keys=on&_busy_timeout=5000"

	db, err := sqlite.OpenDB(context.Background(), dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("successfully connected to SQLite database")

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    mcpserver.ServerName,
			Version: mcpserver.ServerVersion,
		},
		nil,
	)

	mux := http.NewServeMux()

	repository := sqlite.NewRepository(db)

	createTaskUseCase := application.NewCreateTask(repository)
	listProjectsUseCase := application.NewListProjects(repository)
	updateTaskUseCase := application.NewUpdateTaskStatus(repository)
	projectContextUseCase := application.NewGetProjectContext(repository)

	mcpserver.RegisterServerInfo(server)
	mcpserver.RegisterCreateTask(server, createTaskUseCase)
	mcpserver.RegisterListProjects(server, listProjectsUseCase)
	mcpserver.RegisterUpdateTaskStatus(server, updateTaskUseCase)
	mcpserver.RegisterGetProjectContext(server, projectContextUseCase)

	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server { return server },
		nil,
	)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		fmt.Fprint(w, "ok")
	})
	mux.Handle("/mcp", mcpHandler)

	return http.ListenAndServe(":8000", mux)
}
