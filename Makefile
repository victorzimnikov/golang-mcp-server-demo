.PHONY: dev_server watch_server fmt migrate_up check build_server

build_server:
	go build -o ./bin/golang-mcp-server-demo ./cmd/server

dev_server:
	set -a && . ./.env && set +a && go run ./cmd/server

watch_server:
	reflex -s -r '\.go$$' make dev_server

fmt:
	goimports -w .
	go fmt ./...

check:
	go vet ./... && go test ./...

migrate_up:
	set -a && . ./.env && set +a && go tool goose -dir db/migrations postgres "$$DATABASE_URL" up
