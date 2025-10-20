# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Security Features

This application implements several security measures:

### 1. Rate Limiting
- Default: 100 requests per minute per IP
- Configurable via `RATE_LIMIT_MAX` and `RATE_LIMIT_WINDOW` environment variables
- Returns `429 Too Many Requests` when limit is exceeded

### 2. Input Validation
- All user inputs are validated using `go-playground/validator`
- UUID format validation for `external_id`
- Email format validation
- Name length constraints (2-100 characters)
- Required field validation

### 3. Database Security
- Prepared statements via GORM (SQL injection prevention)
- Connection pooling with configurable limits
- SSL mode support for database connections
- Credentials managed via environment variables (never hardcoded)

### 4. Application Security
- Panic recovery middleware
- Structured error messages (no sensitive data leakage)
- Graceful shutdown handling
- Health check endpoint for monitoring

### 5. Container Security
- Non-root user in Docker container
- Minimal Alpine-based image
- Multi-stage build for smaller attack surface
- Health checks configured

## Known Limitations

⚠️ **This is a demonstration project. For production use, consider implementing:**

### Authentication & Authorization
- The API currently has **no authentication** - all endpoints are public
- Recommended: JWT tokens, API keys, or OAuth2
- Role-based access control (RBAC)

### HTTPS/TLS
- Application does not enforce HTTPS
- Recommended: Use reverse proxy (nginx, Caddy) or load balancer with TLS termination

### Additional Security Measures for Production
- Implement CORS policies
- Add request signing/verification
- Implement audit logging
- Add API versioning
- Use secrets management (HashiCorp Vault, AWS Secrets Manager)
- Implement request/response encryption for sensitive data
- Add intrusion detection/prevention
- Implement DDoS protection
- Regular security audits and penetration testing

## Reporting a Vulnerability

If you discover a security vulnerability, please follow these steps:

1. **Do NOT** open a public issue
2. Email the maintainer directly (see GitHub profile)
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work on a fix as soon as possible.

## Security Updates

Security updates will be released as patch versions (e.g., 1.0.1) and announced in:
- GitHub Releases
- CHANGELOG.md
- Security advisories (for critical issues)

## Best Practices for Deployment

### Environment Variables
```bash
# Never commit these to version control
DB_PASSWORD=<strong-random-password>

# Use different credentials for each environment
# Development != Staging != Production
```

### Database
```bash
# Enable SSL in production
DB_SSL=require

# Use connection pooling wisely
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

### Rate Limiting
```bash
# Adjust based on expected traffic
RATE_LIMIT_MAX=100
RATE_LIMIT_WINDOW=1m
```

### Logging
```bash
# Use appropriate log level
LOG_LEVEL=info  # Use 'debug' only in development
```

## Dependency Security

- Dependencies are managed via `go.mod`
- Run `go tool govulncheck ./...` regularly to check for vulnerabilities
- Keep dependencies updated
- Review security advisories from:
  - https://github.com/advisories
  - https://pkg.go.dev/vuln/

## CI/CD Security

GitHub Actions workflow includes:
- Automated testing
- Dependency vulnerability scanning with `govulncheck`
- Linting with `golangci-lint`
- Code quality checks

## Compliance

This project does not claim compliance with any specific security standards (PCI-DSS, HIPAA, SOC2, etc.). 

For compliance requirements, additional measures will be necessary depending on your use case.
