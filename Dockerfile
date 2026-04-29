# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations (remove debug info)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/server ./cmd/server

# Stage 2: Final lean image
FROM alpine:3.19

# Security: run as non-root user
RUN apk --no-cache add ca-certificates tzdata && \
    adduser -D -g '' appuser

WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder /app/bin/server ./server

# Copy static assets needed at runtime (OpenAPI spec)
COPY --from=builder /app/api ./api

# Set ownership to non-root user
RUN chown -R appuser:appuser /home/appuser

USER appuser

EXPOSE 3000

# Healthcheck to ensure container is healthy
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/healthz || exit 1

ENTRYPOINT ["./server"]
