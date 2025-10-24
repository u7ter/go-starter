.PHONY: help build run test clean docker-build docker-up docker-down docker-dev-up docker-dev-down migrate-up migrate-down db-backup swagger lint dev

# Default target
help:
	@echo "Available commands:"
	@echo "  make build            - Build the application"
	@echo "  make run              - Run the application"
	@echo "  make dev              - Run with live reload (requires Air)"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-up        - Start Docker containers (production)"
	@echo "  make docker-down      - Stop Docker containers (production)"
	@echo "  make docker-dev-up    - Start development containers with live reload"
	@echo "  make docker-dev-down  - Stop development containers"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback database migrations"
	@echo "  make db-backup        - Backup database"
	@echo "  make swagger          - Generate Swagger documentation"
	@echo "  make lint             - Run linters"
	@echo "  make fmt              - Format code"
	@echo "  make install-tools    - Install development tools"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app cmd/app/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run the application
run:
	@echo "Running application..."
	go run cmd/app/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t go-starter:latest .

# Start Docker containers (production)
docker-up:
	@echo "Starting Docker containers (production)..."
	@if [ ! -f .env ]; then echo "Error: .env file not found. Copy .env.example to .env first."; exit 1; fi
	docker-compose --env-file .env up -d

# Stop Docker containers (production)
docker-down:
	@echo "Stopping Docker containers (production)..."
	docker-compose down

# Start development containers with live reload
docker-dev-up:
	@echo "Starting development containers with live reload..."
	@if [ ! -f .env ]; then echo "Error: .env file not found. Copy .env.example to .env first."; exit 1; fi
	docker compose --env-file .env -f docker-compose.dev.yml up --build

# Stop development containers
docker-dev-down:
	@echo "Stopping development containers..."
	docker compose --env-file .env -f docker-compose.dev.yml down

# Run database migrations up
migrate-up:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go -direction=up

# Run database migrations down
migrate-down:
	@echo "Rolling back database migrations..."
	go run cmd/migrate/main.go -direction=down

# Backup database
db-backup:
	@echo "Backing up database..."
	@mkdir -p backups
	@TIMESTAMP=$$(date +%Y%m%d_%H%M%S); \
	docker-compose exec -T postgres pg_dump -U $${DB_USER:-app} $${DB_NAME:-appdb} > backups/backup_$$TIMESTAMP.sql && \
	echo "Database backup created: backups/backup_$$TIMESTAMP.sql" || \
	echo "Database backup failed. Make sure Docker containers are running."

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/app/main.go -o docs

# Run linters
lint:
	@echo "Running linters..."
	go fmt ./...
	go vet ./...
	@if command -v staticcheck > /dev/null; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed. Run: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Visit: https://golangci-lint.run/usage/install/"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/cosmtrek/air@latest
	@echo "Development tools installed!"
	@echo "You may also want to install golangci-lint:"
	@echo "Visit: https://golangci-lint.run/usage/install/"

# Development: run with hot reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "air not installed. Install with: make install-tools"; \
		echo "Or run 'make run' for standard mode"; \
	fi
