# wbrestapi

A simple RESTful microservice built in Go using [Fiber v2](https://github.com/gofiber/fiber), [GORM](https://gorm.io), and PostgreSQL.

---

## Architecture

This application has been refactored to follow the principles of **Hexagonal Architecture (also known as Ports and Adapters)**. The goal of this architecture is to isolate the business logic (domain and use cases) from the infrastructure (HTTP handlers, database access, etc.), enabling better testability, maintainability, and flexibility.

### Key Components:

- **Domain Layer**: Contains core business models and logic (`internal/domain`).
- **Application Layer**: Defines use cases and interfaces (ports) required by the domain (`internal/application/usecase` and `internal/application/port`).
- **Adapters Layer**: Implements interfaces for infrastructure components, such as repositories and HTTP handlers (`internal/adapter`).
- **Entrypoint**: Application entrypoint (`cmd/server`) initializes the app and wires everything together.

This separation allows the core application logic to remain completely independent from frameworks, databases, and other external systems.

---

## 📄 Project Assignment

This project implements the following task:
**[View assignment](./assigment.md)**

### Functional Requirements

- `POST /save`: Store user data in the database.
- `GET /:id`: Retrieve user data by `external_id`.

---

## ⚙️ Requirements

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- `make` (optional but recommended)

---

## 🚀 Getting Started

### 1. Clone the repo

```bash
git clone https://github.com/podanypepa/wbrestapi.git
cd wbrestapi
```

### 2. Create `.env` file

Copy `.env.localhost` and modify if needed:

```bash
cp .env.example .env
```

Example content:

```env
PORT=3000
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=users
DB_HOST=localhost
DB_PORT=5432
DB_SSL=disable
```

There is some .env config examples for running on localhost or in Docker.

- [.env.localhost](./.env.localhost)
- [.env.docker](./.env.docker)

### 3. Run the project with Docker Compose

```bash
make docker-up
```

> This will start the app and PostgreSQL database, initialize schema, and make the app available at `http://localhost:3000`.

---

## 🧪 Running Tests

### Integration test

```bash
make test
# or
go test -v ./...
```

This test spins up the application and tests the HTTP endpoints (`/save`, `/{id}`) using a Go `http.Client`.

---

## 📦 Build and Run Locally (without Docker)

```bash
make build
make run

# or run in dev/watch mode
make dev
```

---

## 🧹 Stop and Clean Up

```bash
make docker-down
```

To reset DB data:

```bash
rm -rf ./data
```

---

## 🛠️ CLI Tools for testing

- [create_user.sh](./create_user.sh): shell script for creating a new user from the command line.
- [get_user.sh](./get_user.sh): shell script for retrieving users from database by uuid from command line.

---

## ✅ Features

- JSON REST API with Go + Fiber
- PostgreSQL + GORM
- Graceful shutdown (SIGINT/SIGTERM)
- Rejects new requests during shutdown
- Dockerized and portable
- Integration tested

