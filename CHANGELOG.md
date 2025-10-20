# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-10-20

### Added
- Initial release with hexagonal architecture
- RESTful API with Fiber v2 framework
- PostgreSQL database integration with GORM
- Input validation using go-playground/validator
- Rate limiting middleware (100 req/min default)
- Structured JSON logging with slog
- Panic recovery middleware
- Health check endpoint (`/healthz`)
- Graceful shutdown with configurable timeout
- Database connection pooling with configuration
- Custom error handling and domain-specific errors
- Comprehensive unit tests for all layers
- Integration tests with real HTTP server
- Docker support with multi-stage build
- Docker Compose setup for development
- GitHub Actions CI/CD pipeline with:
  - Automated testing
  - Linting with golangci-lint
  - Security scanning with govulncheck
  - Docker image building
- OpenAPI 3.0 API documentation
- Security policy documentation
- Configuration management via environment variables
- CLI scripts for user creation and retrieval

### Security
- SQL injection prevention via prepared statements
- Input validation for all user inputs
- Rate limiting to prevent abuse
- Non-root Docker container
- Structured error messages (no sensitive data leakage)
- Connection pooling with limits

### Documentation
- Comprehensive README with examples
- OpenAPI 3.0 specification
- Security policy (SECURITY.md)
- Architecture documentation
- Configuration guide
- Contributing guidelines

## [0.1.0] - Initial Development

### Added
- Basic project structure
- Simple user CRUD endpoints
- Database schema
- Docker configuration
