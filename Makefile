.PHONY: build run docker-up docker-down test

BIN_DIR := ./bin
BIN_NAME := app
ENTRYPOINT_PATH := ./cmd/server

dev:
	go tool air

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) $(ENTRYPOINT_PATH)

clean:
	rm -f $(BIN_DIR)/$(BIN_NAME)

run:
	go run $(ENTRYPOINT_PATH)

lint:
	go tool revive ./...
	go tool golangci-lint  run ./...
	go tool gocritic check ./...
	go tool govulncheck ./...

test:
	go test -v ./...

docker-up:
	docker compose --env-file ./.env.docker up -d

docker-down:
	docker compose --env-file ./.env.docker down
