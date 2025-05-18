.PHONY: build run docker-up docker-down test

dev:
	go tool air

build:
	go build .

lint:
	go tool revive ./...
	go tool golangci-lint  run ./...
	go tool gocritic check ./...
	go tool govulncheck ./...

docker-up:
	docker compose --env-file ./.env.docker up

docker-down:
	docker compose --env-file ./.env.docker down

docker-db:
	docker compose --env-file ./.env.docker start db

