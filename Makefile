.PHONY: dev_server watch_server fmt check build_server 

.PHONY: run_inspector

NVM_DIR ?= $(HOME)/.nvm
NODE_VERSION ?= 22

run_inspector:
	@. "$(NVM_DIR)/nvm.sh" && \
	nvm use $(NODE_VERSION) >/dev/null && \
	npx -y @modelcontextprotocol/inspector \
		--server-url http://127.0.0.1:8000/mcp \
		--transport http

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

# Database
.PHONY: migrate_status migrate_up migrate_down migrate_version

DB_PATH ?= ./example.db

migrate_up:
	go tool goose -dir migrations sqlite "$(DB_PATH)" up

migrate_down:
	go tool goose -dir migrations sqlite "$(DB_PATH)" down

migrate_status:
	go tool goose -dir migrations sqlite "$(DB_PATH)" status

migrate_version:
	go tool goose -dir migrations sqlite "$(DB_PATH)" version
