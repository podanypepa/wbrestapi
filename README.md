# wbrestapi

A production-ready RESTful microservice built in Go using [Fiber v2](https://github.com/gofiber/fiber), [GORM](https://gorm.io), and PostgreSQL.

[![CI](https://github.com/podanypepa/wbrestapi/workflows/CI/badge.svg)](https://github.com/podanypepa/wbrestapi/actions)

---

## 🏗️ Architecture

This application follows the principles of **Hexagonal Architecture (Ports and Adapters)**. The layers are strictly decoupled with dedicated models for each layer:

- **Domain Layer** (`internal/domain`): Pure business logic and core entities. No infrastructure tags.
- **Application Layer** (`internal/application`): Use cases and interfaces (ports). Handles coordination.
- **Adapters Layer** (`internal/adapter`):
    - **Handler**: HTTP layer with dedicated **DTOs** (`user_dto.go`) and structured validation.
    - **Repository**: Persistence layer with dedicated **Entities** (`user_entity.go`) and database-specific error handling.
- **Configuration Layer** (`internal/config`): Centralized environment-based configuration.
- **Migrations** (`cmd/server/migrations`): Explicit SQL-based migrations managed by `golang-migrate` and embedded into the binary.

---

## 📚 API & Documentation

- **OpenAPI Spec**: Available in [api/openapi.yaml](./api/openapi.yaml).
- **Swagger UI**: Accessible at `http://localhost:3000/swagger` when the server is running.
- **Metrics**: Prometheus metrics available at `http://localhost:3000/metrics`.

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/healthz` | Health check |
| GET | `/metrics` | Prometheus metrics |
| GET | `/swagger/*` | Interactive API documentation |
| POST | `/save` | Create/Update user (requires valid UUID, age >= 15) |
| GET | `/:id` | Retrieve user by external UUID |

---

## ✅ Key Features

- 🛡️ **Security**: Helmet headers, CORS middleware, and input sanitization.
- 🚦 **Resilience**: Rate limiting, Graceful shutdown (SIGINT/SIGTERM), and Context propagation.
- 📝 **Observability**: Structured JSON logging (`slog`) and Prometheus metrics.
- 🔐 **Validation**: Multi-layer validation (DTO format + Domain business invariants).
- 🔄 **Database**: Explicit SQL migrations, connection pooling, and typed error handling (Postgres/SQLite).
- 🐳 **Optimized Docker**: Secure multi-stage build running as a non-root user.
- 🧪 **Testing**: High coverage with unit tests (mockery) and integration tests.

---

## 🚀 Getting Started

### 1. Requirements
- [Go 1.25+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- `make`

### 2. Setup
```bash
cp .env.localhost .env
```

### 3. Run with Docker
```bash
make docker-up
```
> App will be available at `http://localhost:3000`.

### 4. Run Locally
```bash
make tidy
make build
make run
# or development mode with hot-reload
make dev
```

---

## 🧪 Development & Quality

| Command | Description |
|---------|-------------|
| `make help` | Show all available commands |
| `make test` | Run all unit tests |
| `make lint` | Run exhaustive linters (golangci-lint, revive, etc.) |
| `make fmt` | Format source code |
| `make mocks` | Regenerate mocks for testing |

---

## 🏗️ Project Structure

```
wbrestapi/
├── api/                  # OpenAPI 3.0 specification
├── cmd/
│   └── server/          # Application entrypoint
│       └── migrations/  # SQL migration files (embedded)
├── internal/
│   ├── adapter/         # Infrastructure adapters
│   │   ├── handler/     # HTTP handlers & DTOs
│   │   └── repository/  # GORM implementations & Entities
│   ├── application/     # Application layer (Use Cases & Ports)
│   ├── config/          # Configuration management
│   └── domain/          # Pure domain models & invariants
├── Dockerfile           # Optimized multi-stage build
├── Makefile             # Development automation
└── README.md
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
