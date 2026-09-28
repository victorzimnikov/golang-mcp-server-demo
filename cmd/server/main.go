package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/victorzimnikov/golang-mcp-server-demo/internal/mcpserver"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dns := "database.db"

	db, err := sql.Open("sqlite", dns)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Successfully connected to SQLite database")

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    mcpserver.ServerName,
			Version: mcpserver.ServerVersion,
		},
		nil,
	)

	mux := http.NewServeMux()

	mcpserver.RegisterServerInfo(server)

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
