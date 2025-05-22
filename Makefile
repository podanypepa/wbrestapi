.PHONY: build run docker-up docker-down test

dev:
	go tool air

build:
	go build .

clean:
	rm -f ./wbrestapi

run:
	go run .

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
