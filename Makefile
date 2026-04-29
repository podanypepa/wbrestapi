.PHONY: help build run dev test lint fmt tidy clean docker-up docker-down mocks

BIN_DIR := ./bin
BIN_NAME := app
ENTRYPOINT_PATH := ./cmd/server

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

dev: ## Run application in development mode with hot-reload
	go tool air

build: ## Build the application binary
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) $(ENTRYPOINT_PATH)

clean: ## Remove binary and build artifacts
	rm -f $(BIN_DIR)/$(BIN_NAME)
	rm -rf ./mocks

run: ## Run the application
	go run $(ENTRYPOINT_PATH)

fmt: ## Format source code
	go fmt ./...

tidy: ## Tidy up go.mod and go.sum
	go mod tidy

lint: ## Run all linters (revive, golangci-lint, gocritic, govulncheck)
	go tool revive ./...
	go tool golangci-lint run ./...
	go tool gocritic check ./...
	go tool govulncheck ./...

test: ## Run unit tests
	go test -v ./...

mocks: ## Generate mocks using mockery
	go tool mockery --all --recursive --inpackage

docker-up: ## Start application in docker containers
	docker compose --env-file ./.env.docker up -d

docker-down: ## Stop and remove docker containers
	docker compose --env-file ./.env.docker down
