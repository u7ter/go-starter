# Go Starter - RESTful API

[![CI](https://github.com/u7ter/go-starter/actions/workflows/ci.yml/badge.svg)](https://github.com/u7ter/go-starter/actions/workflows/ci.yml)
[![CD](https://github.com/u7ter/go-starter/actions/workflows/cd.yml/badge.svg)](https://github.com/u7ter/go-starter/actions/workflows/cd.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/u7ter/go-starter)](https://goreportcard.com/report/github.com/u7ter/go-starter)
[![codecov](https://codecov.io/gh/u7ter/go-starter/branch/main/graph/badge.svg)](https://codecov.io/gh/u7ter/go-starter)

A production-ready RESTful HTTP API built with Go, featuring JWT authentication, rate limiting, structured logging, and PostgreSQL database.

## Features

- **RESTful HTTP API** with clean architecture
- **JWT Authentication** (HS256) with secure token handling
- **Rate Limiting** using token bucket algorithm
- **Structured Logging** with Zap (JSON format)
- **PostgreSQL** database with parameterized queries
- **Database Migrations** using golang-migrate
- **Docker & Docker Compose** support
- **Live Reload** in development with Air
- **CI/CD** with GitHub Actions
- **Swagger/OpenAPI** documentation
- **Graceful Shutdown** handling
- **Security Headers** (X-Content-Type-Options, X-Frame-Options, HSTS)
- **Health Check** endpoints

## Project Structure

```
.
├── cmd/
│   ├── app/           # Main application entry point
│   └── migrate/       # Database migration tool
├── internal/
│   ├── config/        # Configuration management
│   ├── handlers/      # HTTP handlers
│   ├── logger/        # Structured logging
│   ├── middleware/    # HTTP middleware (auth, rate limit, etc.)
│   ├── migrations/    # SQL migration files
│   ├── models/        # Data structures
│   ├── repositories/  # Database layer
│   └── services/      # Business logic
├── pkg/
│   └── database/      # Database connection utilities
├── docs/              # Swagger documentation
├── scripts/           # Helper scripts
├── Dockerfile         # Multi-stage Docker build
├── docker-compose.yml # Docker Compose configuration
├── Makefile           # Build automation
└── .env.example       # Environment variables template
```

## Prerequisites

- Go 1.21 or newer
- PostgreSQL 15+
- Docker & Docker Compose (for containerized deployment)

## Quick Start

### 1. Clone and Configure

```bash
# Copy environment variables
cp .env.example .env

# Edit .env with your configuration
# IMPORTANT: Change JWT_SECRET and DB_PASSWORD in production!
```

### 2. Local Development

#### Option A: With Live Reload (Recommended for Development)

```bash
# Install development tools including Air
make install-tools

# Start PostgreSQL (via Docker)
docker-compose up -d postgres

# Run migrations
make migrate-up

# Run with live reload
make dev
```

The application will automatically restart when you modify Go files!

#### Option B: Standard Mode

```bash
# Install dependencies
go mod download

# Start PostgreSQL (via Docker)
docker-compose up -d postgres

# Run migrations
make migrate-up

# Run the application
make run
```

The API will be available at `http://localhost:8080`

### 3. Development with Docker (Live Reload)

```bash
# Start development environment with live reload
make docker-dev-up

# View logs
docker-compose -f docker-compose.dev.yml logs -f app-dev

# Stop development environment
make docker-dev-down
```

This starts PostgreSQL, runs migrations, and starts the app with Air for live reload!

### 4. Production Deployment with Docker

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up -d
```

## API Endpoints

### Health Check
- `GET /healthz` - Health check (checks database connectivity)
- `GET /ready` - Readiness check

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login and receive JWT token

### Swagger Documentation
- `GET /swagger/index.html` - API documentation (development mode only)

## Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### Using Authenticated Endpoints

```bash
curl -X GET http://localhost:8080/protected-endpoint \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Makefile Commands

```bash
make help              # Show all available commands
make build             # Build the application
make run               # Run the application
make dev               # Run with live reload (requires Air)
make test              # Run tests
make test-coverage     # Run tests with coverage report
make docker-build      # Build Docker image
make docker-up         # Start Docker containers (production)
make docker-down       # Stop Docker containers (production)
make docker-dev-up     # Start development containers with live reload
make docker-dev-down   # Stop development containers
make migrate-up        # Run database migrations
make migrate-down      # Rollback database migrations
make db-backup         # Backup database
make swagger           # Generate Swagger documentation
make lint              # Run linters (go fmt, go vet, staticcheck)
make fmt               # Format code
make install-tools     # Install development tools (Air, Swag, etc.)
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `app` |
| `DB_PASSWORD` | Database password | *required* |
| `DB_NAME` | Database name | `appdb` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `JWT_SECRET` | JWT signing secret | *required* |
| `RATE_LIMIT_RPS` | Rate limit (requests/sec) | `10` |
| `RATE_LIMIT_BURST` | Rate limit burst | `20` |
| `LOG_LEVEL` | Logging level | `info` |
| `ENV` | Environment (development/production) | `development` |

## Database Migrations

Migrations are stored in `internal/migrations/` directory.

### Create a New Migration

```bash
# Create migration files manually
touch internal/migrations/000002_add_users_index.up.sql
touch internal/migrations/000002_add_users_index.down.sql
```

### Run Migrations

```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Or use the migrate tool directly
go run cmd/migrate/main.go -direction=up
```

## Database Backup

```bash
# Backup database to backups/ directory
make db-backup
```

This creates a timestamped SQL dump file.

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage in browser
open coverage.html
```

## Security Features

1. **JWT Authentication**: Tokens expire after 24 hours
2. **Password Hashing**: Using bcrypt with default cost
3. **SQL Injection Protection**: All queries are parameterized
4. **Rate Limiting**: IP-based request limiting
5. **Security Headers**: X-Content-Type-Options, X-Frame-Options, HSTS (production)
6. **Input Validation**: Using go-playground/validator

## Production Deployment

### Docker Production Build

The Dockerfile uses multi-stage builds with a distroless final image for minimal attack surface:

```bash
docker build -t go-starter:latest .
docker run -p 8080:8080 --env-file .env go-starter:latest
```

### Production Configuration

1. Set `ENV=production` in environment variables
2. Use strong `JWT_SECRET` (64+ random characters)
3. Use strong `DB_PASSWORD`
4. Enable SSL for database (`DB_SSLMODE=require`)
5. Configure proper rate limits
6. Disable Swagger in production (automatic)
7. Use HTTPS/TLS termination (nginx, load balancer, etc.)

## CI/CD with GitHub Actions

The project includes comprehensive CI/CD pipelines:

### Continuous Integration (CI)

Automatically runs on pull requests and pushes to `main` and `develop` branches:

- **Linting**: `go fmt`, `go vet`, `staticcheck`, `golangci-lint`
- **Testing**: Unit and integration tests with PostgreSQL
- **Coverage**: Code coverage reporting (with Codecov integration)
- **Build**: Multi-platform binary builds
- **Docker**: Docker image build verification
- **Security**: Gosec security scanning

### Continuous Deployment (CD)

Runs on pushes to `main` and version tags:

- **Docker Build & Push**: Builds and pushes images to GitHub Container Registry
- **Multi-platform**: Supports `linux/amd64` and `linux/arm64`
- **Releases**: Creates GitHub releases with binaries for:
  - Linux (amd64, arm64)
  - macOS (amd64, arm64)
  - Windows (amd64)
- **Deployment**: Example deployment configurations (uncomment and customize)

### Workflow Files

- `.github/workflows/ci.yml` - Continuous Integration
- `.github/workflows/cd.yml` - Continuous Deployment

### Setting Up CI/CD

1. Push your code to GitHub
2. CI will run automatically on pull requests
3. For CD, ensure `GITHUB_TOKEN` has proper permissions
4. For Codecov integration, add `CODECOV_TOKEN` secret
5. Customize deployment steps for your infrastructure

## Monitoring

- All requests are logged with structured JSON format
- Each request includes a unique `request_id` for tracing
- Health check endpoint for load balancer integration
- Database connection health monitoring

## Code Quality

```bash
# Format code
make fmt

# Run linters
make lint

# The project follows Go standards:
# - go fmt
# - go vet
# - staticcheck
```

## License

This is a starter template. Use it as you see fit for your projects.

## Contributing

Feel free to submit issues and pull requests.
