# Development Guide

This guide covers development workflows, tools, and best practices for working on this project.

## Table of Contents

- [Development Environment Setup](#development-environment-setup)
- [Live Reload with Air](#live-reload-with-air)
- [Testing](#testing)
- [Docker Development](#docker-development)
- [CI/CD Pipelines](#cicd-pipelines)
- [Best Practices](#best-practices)

## Development Environment Setup

### Prerequisites

- Go 1.23+ (or the version specified in go.mod)
- Docker & Docker Compose
- Make
- Git

### Initial Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-starter
   ```

2. **Install development tools**
   ```bash
   make install-tools
   ```

   This installs:
   - Air (live reload)
   - Swag (Swagger documentation)
   - Staticcheck (linter)

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

4. **Start the database**
   ```bash
   docker-compose up -d postgres
   ```

5. **Run migrations**
   ```bash
   make migrate-up
   ```

## Live Reload with Air

Air provides automatic recompilation and reloading during development.

### Local Development with Air

```bash
# Start the app with live reload
make dev
```

Air watches for file changes and automatically:
- Recompiles the application
- Restarts the server
- Shows build errors in the terminal

### Configuration

Air configuration is in `.air.toml`. Key settings:

- **Build command**: `go build -o ./tmp/main ./cmd/app`
- **Watched directories**: All except `assets`, `tmp`, `vendor`, `testdata`
- **Watched extensions**: `.go`, `.tpl`, `.tmpl`, `.html`
- **Excluded files**: `*_test.go`

### Docker Development with Air

For a fully containerized development environment:

```bash
# Start dev environment with live reload
make docker-dev-up

# View logs
docker-compose -f docker-compose.dev.yml logs -f app-dev

# Stop dev environment
make docker-dev-down
```

The Docker dev setup:
- Mounts source code as a volume
- Runs Air inside the container
- Automatically reloads on code changes
- Includes PostgreSQL and automatic migrations

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# View coverage report in browser
open coverage.html
```

### Writing Tests

Place test files next to the code they test:

```
internal/
  services/
    auth_service.go
    auth_service_test.go
```

Example test structure:

```go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestAuthService_Login(t *testing.T) {
    // Setup
    // ...

    // Test
    result, err := authService.Login(ctx, req)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Integration Tests

Integration tests require a database. Use `testcontainers-go` or ensure PostgreSQL is running:

```bash
# Start test database
docker-compose up -d postgres

# Run tests
make test
```

## Docker Development

### Development vs Production

**Development** (`docker-compose.dev.yml`):
- Mounts source code as volume
- Uses Air for live reload
- Debug logging enabled
- Swagger UI enabled

**Production** (`docker-compose.yml`):
- Compiled binary in distroless image
- No source code mounted
- Production logging
- Swagger UI disabled

### Common Docker Commands

```bash
# Build development image
docker-compose -f docker-compose.dev.yml build

# Start development environment
make docker-dev-up

# View logs
docker-compose -f docker-compose.dev.yml logs -f

# Execute commands in container
docker-compose -f docker-compose.dev.yml exec app-dev sh

# Restart a service
docker-compose -f docker-compose.dev.yml restart app-dev

# Clean up
make docker-dev-down
docker system prune -f
```

## CI/CD Pipelines

### GitHub Actions Workflows

#### CI Workflow (`.github/workflows/ci.yml`)

Triggered on:
- Pull requests to `main` or `develop`
- Pushes to `main` or `develop`

Steps:
1. **Lint**: Code formatting and static analysis
2. **Test**: Unit and integration tests with PostgreSQL
3. **Build**: Compile application
4. **Docker**: Build Docker image
5. **Security**: Gosec security scan

#### CD Workflow (`.github/workflows/cd.yml`)

Triggered on:
- Pushes to `main` branch
- Version tags (`v*`)

Steps:
1. **Build & Push**: Docker image to GitHub Container Registry
2. **Multi-platform**: Builds for amd64 and arm64
3. **Release**: Creates GitHub release with binaries
4. **Deploy**: Example deployment steps (customize)

### Running CI Checks Locally

Before pushing:

```bash
# Format code
make fmt

# Run linters
make lint

# Run tests
make test

# Build application
make build
```

### Setting Up CI/CD

1. **Push to GitHub**
   ```bash
   git remote add origin <your-repo-url>
   git push -u origin main
   ```

2. **Configure Secrets** (in GitHub repository settings):
   - `CODECOV_TOKEN` - For coverage reporting (optional)
   - `GITHUB_TOKEN` - Automatically provided

3. **Enable GitHub Container Registry**
   - Settings → Actions → General
   - Enable "Read and write permissions" for workflows

4. **Monitor Workflows**
   - Visit Actions tab in your repository
   - View workflow runs and logs

## Best Practices

### Code Style

1. **Use `go fmt`** before committing
   ```bash
   make fmt
   ```

2. **Run linters** regularly
   ```bash
   make lint
   ```

3. **Follow Go naming conventions**
   - Exported: `UserRepository`
   - Unexported: `validateToken`
   - Acronyms: `HTTPServer`, `userID`

### Git Workflow

1. **Create feature branch**
   ```bash
   git checkout -b feature/user-profile
   ```

2. **Make changes and commit**
   ```bash
   git add .
   git commit -m "feat: add user profile endpoint"
   ```

3. **Push and create PR**
   ```bash
   git push origin feature/user-profile
   ```

### Commit Messages

Follow conventional commits:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `test:` Test changes
- `refactor:` Code refactoring
- `chore:` Maintenance tasks

Examples:
```
feat: add user registration endpoint
fix: handle database connection timeout
docs: update API documentation
test: add tests for auth service
```

### Database Migrations

1. **Create migration files**
   ```bash
   # Naming: NNNNNN_description.up.sql and NNNNNN_description.down.sql
   touch internal/migrations/000002_add_users_name.up.sql
   touch internal/migrations/000002_add_users_name.down.sql
   ```

2. **Write SQL in up migration**
   ```sql
   ALTER TABLE users ADD COLUMN name VARCHAR(255);
   ```

3. **Write rollback in down migration**
   ```sql
   ALTER TABLE users DROP COLUMN name;
   ```

4. **Apply migration**
   ```bash
   make migrate-up
   ```

5. **Test rollback**
   ```bash
   make migrate-down
   make migrate-up
   ```

### API Development

1. **Add Swagger annotations**
   ```go
   // @Summary Create user
   // @Tags users
   // @Accept json
   // @Produce json
   // @Param request body CreateUserRequest true "User data"
   // @Success 201 {object} User
   // @Router /users [post]
   func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
       // ...
   }
   ```

2. **Generate Swagger docs**
   ```bash
   make swagger
   ```

3. **Test endpoint**
   ```bash
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com"}'
   ```

### Debugging

#### Local Debugging

1. **Enable debug logging**
   ```bash
   # In .env
   LOG_LEVEL=debug
   ```

2. **Use Delve debugger**
   ```bash
   # Install
   go install github.com/go-delve/delve/cmd/dlv@latest

   # Debug
   dlv debug ./cmd/app
   ```

#### Docker Debugging

1. **View container logs**
   ```bash
   docker-compose -f docker-compose.dev.yml logs -f app-dev
   ```

2. **Execute shell in container**
   ```bash
   docker-compose -f docker-compose.dev.yml exec app-dev sh
   ```

3. **Inspect database**
   ```bash
   docker-compose exec postgres psql -U app -d appdb
   ```

### Performance

1. **Profile CPU**
   ```bash
   go test -cpuprofile=cpu.prof -bench=.
   go tool pprof cpu.prof
   ```

2. **Profile Memory**
   ```bash
   go test -memprofile=mem.prof -bench=.
   go tool pprof mem.prof
   ```

3. **Benchmark**
   ```go
   func BenchmarkAuthService_Login(b *testing.B) {
       for i := 0; i < b.N; i++ {
           authService.Login(ctx, req)
       }
   }
   ```

## Troubleshooting

### Air not reloading

1. Check `.air.toml` configuration
2. Ensure files are saved
3. Check Air logs for errors
4. Try `make clean` then `make dev`

### Database connection issues

1. Check PostgreSQL is running:
   ```bash
   docker-compose ps postgres
   ```

2. Verify connection settings in `.env`

3. Check logs:
   ```bash
   docker-compose logs postgres
   ```

### Port already in use

```bash
# Find process using port 8080
lsof -ti:8080

# Kill process
kill -9 $(lsof -ti:8080)
```

### Docker issues

```bash
# Clean up everything
docker-compose down -v
docker system prune -af

# Rebuild
make docker-dev-up
```

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Air GitHub](https://github.com/cosmtrek/air)
- [Swagger Specification](https://swagger.io/specification/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Docker Documentation](https://docs.docker.com/)

## Getting Help

- Check existing issues and PRs
- Review documentation in `/docs`
- Run `make help` for available commands
- Check CI logs for test failures
