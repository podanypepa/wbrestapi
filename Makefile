.PHONY: build run docker-up docker-down test

dev:
	go tool air

build:
	go build .

run:
	go run .

lint:
	go tool revive ./...
	go tool golangci-lint  run ./...
	go tool gocritic check ./...
	go tool govulncheck ./...

test:
	go test -v -run ./...

docker-up:
	docker compose --env-file ./.env.docker up

docker-down:
	docker compose --env-file ./.env.docker down
