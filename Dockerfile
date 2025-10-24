# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o /app/bin/app \
    ./cmd/app

# Final stage - distroless
FROM gcr.io/distroless/static:nonroot

# Copy timezone data and certificates from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /app/bin/app /app

# Copy migrations
COPY --from=builder /app/internal/migrations /internal/migrations

# Use non-root user
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app"]
