# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

#### Core Features
- RESTful HTTP API with Clean Architecture
- JWT authentication (HS256) with 24-hour token expiration
- Token bucket rate limiting (IP-based)
- Structured logging with Zap (JSON format)
- PostgreSQL database integration with pgx/v5
- Database migrations with golang-migrate
- Health check endpoints (`/healthz`, `/ready`)
- Graceful shutdown handling (SIGINT/SIGTERM)
- Input validation with go-playground/validator
- Password hashing with bcrypt
- Security headers middleware

#### Development Tools
- **Live Reload**: Air configuration for automatic recompilation
- **Docker Dev Environment**: `docker-compose.dev.yml` with live reload
- **Swagger/OpenAPI**: Auto-generated API documentation
- **Makefile**: Comprehensive build automation
- **golangci-lint**: Configuration for code quality

#### CI/CD
- **GitHub Actions CI**:
  - Linting (go fmt, go vet, staticcheck, golangci-lint)
  - Unit and integration testing with PostgreSQL
  - Code coverage reporting (Codecov integration)
  - Multi-platform builds
  - Docker image building
  - Security scanning with Gosec

- **GitHub Actions CD**:
  - Automated Docker image builds and pushes to GHCR
  - Multi-platform support (linux/amd64, linux/arm64)
  - GitHub releases with binaries for Linux, macOS, Windows
  - Example deployment configurations

#### Documentation
- `README.md` - Comprehensive project documentation
- `DEVELOPMENT.md` - Developer guide with workflows and best practices
- `IMPLEMENTATION.md` - Technical implementation checklist
- `CHANGELOG.md` - Version history and changes
- API documentation via Swagger UI (development mode)

#### Project Structure
```
cmd/
  app/          - Main application entry point
  migrate/      - Database migration tool
internal/
  config/       - Environment-based configuration
  handlers/     - HTTP request handlers
  logger/       - Structured logging
  middleware/   - Auth, rate limit, logging, security
  migrations/   - SQL migration files
  models/       - Data structures
  repositories/ - Database access layer
  services/     - Business logic
pkg/
  database/     - Database utilities
.github/
  workflows/    - CI/CD pipelines
```

#### Configuration
- Environment variable-based configuration
- `.env` file support for local development
- Validation of required configuration at startup
- Production/development mode switching

#### Security Features
- Parameterized SQL queries (SQL injection protection)
- JWT token validation
- Rate limiting per IP address
- Security headers (X-Content-Type-Options, X-Frame-Options, HSTS)
- Input validation
- Bcrypt password hashing

#### Docker
- Production Dockerfile with multi-stage build
- Distroless base image for minimal attack surface
- Development Dockerfile with Air
- Docker Compose for local development
- Docker Compose for production deployment
- Health checks for database

#### Database
- Users table with email, password_hash, timestamps
- Index on email for performance
- Migration system for schema versioning
- Database backup command in Makefile

#### API Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login (returns JWT)
- `GET /healthz` - Health check
- `GET /ready` - Readiness check
- `GET /swagger/index.html` - API documentation (dev only)

#### Make Commands
- `make help` - Show all available commands
- `make build` - Build application
- `make run` - Run application
- `make dev` - Run with live reload
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make docker-build` - Build Docker image
- `make docker-up` - Start production containers
- `make docker-down` - Stop production containers
- `make docker-dev-up` - Start development containers with live reload
- `make docker-dev-down` - Stop development containers
- `make migrate-up` - Apply database migrations
- `make migrate-down` - Rollback migrations
- `make db-backup` - Backup database to SQL file
- `make swagger` - Generate Swagger documentation
- `make lint` - Run all linters
- `make fmt` - Format code
- `make install-tools` - Install development tools

### Dependencies
- `gorilla/mux` - HTTP router
- `pgx/v5` - PostgreSQL driver
- `zap` - Structured logging
- `jwt/v5` - JWT tokens
- `bcrypt` - Password hashing
- `validator/v10` - Input validation
- `godotenv` - Environment variables
- `golang-migrate/migrate` - Database migrations
- `swag` - Swagger generation
- `air` - Live reload
- `uuid` - Request ID generation
- `golang.org/x/time/rate` - Rate limiting

### Best Practices Implemented
- Clean Architecture with layered structure
- Dependency injection
- Context propagation
- Error handling with wrapped errors
- Structured logging with request IDs
- Parameterized database queries
- Environment-based configuration
- Graceful shutdown
- Docker multi-stage builds
- CI/CD automation
- Code quality checks
- Security scanning

## [1.0.0] - Initial Release (Not yet released)

This version represents the initial production-ready implementation of the Go Starter API.

[Unreleased]: https://github.com/YOUR_USERNAME/go-starter/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/YOUR_USERNAME/go-starter/releases/tag/v1.0.0
