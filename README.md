# wbrestapi

A simple RESTful microservice built in Go using [Fiber v2](https://github.com/gofiber/fiber), [GORM](https://gorm.io), and PostgreSQL.

## 📄 Project Assignment

This project implements the following task:  
**[View assignment (PDF)](./Golang%20Zadani%20(1).pdf)**

### Functional Requirements

- `POST /save`: Store user data in the database.
- `GET /:id`: Retrieve user data by `external_id`.

---

## ⚙️ Requirements

- [Go 1.21+](https://golang.org/dl/)
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
```

This test spins up the application and tests the HTTP endpoints (`/save`, `/{id}`) using a Go `http.Client`.

---

## 🛠 Project Structure

```
.
├── main.go             # Main application entry point
├── Dockerfile          # Multi-stage Docker build
├── docker-compose.yml  # Compose config with PostgreSQL
├── .env.example        # Sample environment variables
├── Makefile            # Convenient task runner
└── integration_test.go # HTTP-based integration test
```

---

## 📦 Build and Run Locally (without Docker)

```bash
make build
make run

# or run in dev/watch mode
make air
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

## ✅ Features

- JSON REST API with Go + Fiber
- PostgreSQL + GORM
- Graceful shutdown (SIGINT/SIGTERM)
- Rejects new requests during shutdown
- Dockerized and portable
- Integration tested

---

