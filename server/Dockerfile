# ============================================
# T026: Go Server Dockerfile (Multi-stage)
# Distroless final image for security
# ============================================

# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies including gcc for SQLite CGO
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Copy go mod files first (cache dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /server ./cmd/jail-server

# ============================================
# Stage 2: Runtime (Distroless for security)
# ============================================
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /server /app/server

# Copy any static assets if needed
# COPY --from=builder /app/assets /app/assets

# Non-root user (distroless default)
USER nonroot:nonroot

# Expose ports
EXPOSE 8080

# Health check
# Note: distroless doesn't have shell, so we use grpc_health_probe or similar
# For now, rely on orchestrator health checks

# Run
ENTRYPOINT ["/app/server"]
