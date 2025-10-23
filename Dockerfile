# Build stage
FROM golang:1.24.9-alpine AS builder

# Install build dependencies
RUN apk update && apk add --no-cache git build-base

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
ENV CGO_ENABLED=1
RUN go build -ldflags="-w -s" -o citizens-data-webservice ./cmd/citizens-data-webservice/main.go

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk update && apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/citizens-data-webservice .

# Copy config files
COPY --chown=appuser:appuser config ./config

# Switch to non-root user
USER appuser

# Set environment variables
ENV CONFIG_PATH=config/prod.yaml

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application
ENTRYPOINT ["./citizens-data-webservice"]