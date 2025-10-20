# wbrestapi

A simple RESTful microservice built in Go using [Fiber v2](https://github.com/gofiber/fiber), [GORM](https://gorm.io), and PostgreSQL.

[![CI](https://github.com/podanypepa/wbrestapi/workflows/CI/badge.svg)](https://github.com/podanypepa/wbrestapi/actions)

---

## Architecture

This application follows the principles of **Hexagonal Architecture (Ports and Adapters)**. The goal is to isolate the business logic from infrastructure, enabling better testability, maintainability, and flexibility.

### Key Components:

- **Domain Layer**: Core business models and logic (`internal/domain`).
- **Application Layer**: Use cases and interfaces/ports (`internal/application`).
- **Adapters Layer**: Infrastructure implementations - HTTP handlers, database repositories (`internal/adapter`).
- **Configuration Layer**: Centralized configuration management (`internal/config`).
- **Entrypoint**: Application initialization and wiring (`cmd/server`).

This separation keeps the core application logic independent from frameworks, databases, and external systems.

---

## 📄 Project Assignment

This project implements the following task:
**[View assignment](./assigment.md)**

### Functional Requirements

- `POST /save`: Store user data in the database.
- `GET /:id`: Retrieve user data by `external_id`.

---

## 📚 API Documentation

API documentation is available in OpenAPI 3.0 format: **[OpenAPI Specification](./api/openapi.yaml)**

### API Endpoints

#### Health Check
```bash
GET /healthz
```
**Response:**
```json
{"status": "ok"}
```

#### Create User
```bash
POST /save
Content-Type: application/json

{
  "external_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "date_of_birth": "1990-01-15T00:00:00Z"
}
```

**Success Response (201):**
```json
{
  "id": 1,
  "external_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "date_of_birth": "1990-01-15T00:00:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid input or validation error
- `409 Conflict`: User with this external_id already exists
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

#### Get User
```bash
GET /{external_id}
```

**Success Response (200):** Same as Create User response

**Error Responses:**
- `404 Not Found`: User not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

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
cp .env.localhost .env
```

Environment config examples:
- [.env.localhost](./.env.localhost) - For local development
- [.env.docker](./.env.docker) - For Docker Compose

### 3. Run the project with Docker Compose

```bash
make docker-up
```

> This will start the app and PostgreSQL database, initialize schema, and make the app available at `http://localhost:3000`.

---

## 🧪 Running Tests

### Unit Tests

```bash
make test
# or
go test -v ./...
```

### Integration Tests

Integration tests start a real server and test HTTP endpoints:

```bash
cd cmd/server
go test -v
```

**Note:** Integration tests require PostgreSQL running on localhost:5432

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

- [create_user.sh](./create_user.sh): Create a new user from the command line
- [get_user.sh](./get_user.sh): Retrieve users by UUID from the command line

---

## ✅ Features

- ✨ JSON REST API with Go + Fiber
- 🐘 PostgreSQL + GORM
- 🔐 Input validation with structured errors
- 🚦 Rate limiting (100 req/min by default)
- 📝 Structured logging (JSON format)
- 🛡️ Panic recovery middleware
- 🏥 Health check endpoint
- ⚡ Graceful shutdown (SIGINT/SIGTERM)
- 🐳 Dockerized and portable
- ✅ Comprehensive test coverage (unit + integration)
- 📊 Database connection pooling
- 🔄 CI/CD with GitHub Actions
- 📖 OpenAPI 3.0 documentation
- 🎯 Hexagonal Architecture

---

## 🔧 Configuration

Configure via environment variables:

### Server Configuration
- `PORT` - Server port (default: 3000)
- `SHUTDOWN_TIMEOUT` - Graceful shutdown timeout (default: 5s)
- `RATE_LIMIT_MAX` - Max requests per window (default: 100)
- `RATE_LIMIT_WINDOW` - Rate limit time window (default: 1m)
- `LOG_LEVEL` - Logging level: info, debug (default: info)

### Database Configuration
- `DB_HOST` - Database host
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `DB_PORT` - Database port
- `DB_SSL` - SSL mode (disable, require, etc.)
- `DB_MAX_OPEN_CONNS` - Max open connections (default: 25)
- `DB_MAX_IDLE_CONNS` - Max idle connections (default: 5)
- `DB_CONN_MAX_LIFETIME` - Connection max lifetime (default: 5m)

---

## 🏗️ Project Structure

```
wbrestapi/
├── api/                      # API documentation
│   └── openapi.yaml         # OpenAPI 3.0 specification
├── cmd/
│   └── server/              # Application entrypoint
│       ├── main.go
│       └── integration_test.go
├── internal/
│   ├── adapter/             # Infrastructure adapters
│   │   ├── handler/         # HTTP handlers
│   │   └── repository/      # Database repositories
│   ├── application/         # Application layer
│   │   ├── port/           # Interfaces (ports)
│   │   └── usecase/        # Use cases
│   ├── config/             # Configuration
│   └── domain/             # Domain models and logic
├── .github/
│   └── workflows/          # CI/CD pipelines
├── Dockerfile
├── compose.yaml
├── Makefile
└── README.md
```

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
