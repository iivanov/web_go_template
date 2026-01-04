# AI Agent Guidelines

Guidelines for AI coding assistants working with this Go project template.

## Build Requirements

**Critical:** This project uses Go's experimental JSON v2 encoder. Always use:
```bash
GOEXPERIMENT=jsonv2 go build ...
GOEXPERIMENT=jsonv2 go test ...
```

Use `make build` or `make test` which handle this automatically.

## Project Conventions

### Directory Structure

| Directory | Purpose |
|-----------|---------|
| `cmd/` | CLI commands (Cobra). Add new subcommands here. |
| `internal/app/` | Application bootstrap. Module composition happens in `app.go`. |
| `internal/shared/` | Code shared across bounded contexts (events, common types). |
| `internal/<context>/` | Bounded contexts. Each is self-contained. |
| `pkg/` | Reusable infrastructure packages (database, logger, messagebus, etc.). |

### Bounded Context Structure

Each bounded context under `internal/` follows this structure:
```
internal/<contextname>/
├── module.go           # fx module - MUST register all dependencies
├── config/             # Context-specific configuration structs
├── controllers/        # HTTP handlers (one file per endpoint)
├── services/           # Business logic
├── repositories/       # Data access layer
├── entities/           # Domain models / GORM entities
├── dto/                # Request/response DTOs
└── handlers/           # Message bus event handlers
```

### Dependency Injection (fx)

- Every package that provides dependencies MUST export a `Module` variable
- Use `fx.Provide()` for constructors
- Use `fx.Invoke()` for side effects (migrations, starting services)
- Register new modules in `internal/app/app.go` → `generalModules()`

### HTTP Handlers

Handlers must implement the `Route` interface:
```go
type Route interface {
    http.Handler
    Pattern() string  // e.g., "GET /api/users/{id}"
}
```

Registration pattern:
```go
// In module.go
webserver.AsRoute(controller.NewMyHandler)
```

### Message Bus Handlers

Handlers must implement the `Handler` interface:
```go
type Handler interface {
    Topic() string
    Handle(ctx context.Context, payload []byte) error
}
```

Registration pattern:
```go
// In module.go
messagebus.AsHandler(handlers.NewMyHandler)
```

### Events

Define events in `internal/shared/events/`:
```go
const TopicMyEvent = "domain.event"

type MyEvent struct {
    // fields
}

func (e MyEvent) Topic() string {
    return TopicMyEvent
}
```

### Configuration

- Config structs use `mapstructure` tags for Viper
- Use `default` tags from `github.com/creasty/defaults`
- Environment variables use prefix `APP_` with `_` as separator
- Add new config sections to `internal/app/config.go` → `Config` struct

## Common Tasks

### Add New Bounded Context

1. Create `internal/<contextname>/module.go`
2. Add to `internal/app/app.go`:
   ```go
   import "project_template/internal/<contextname>"
   
   func generalModules() []fx.Option {
       return []fx.Option{
           // ...
           <contextname>.Module,
       }
   }
   ```

### Add New HTTP Endpoint

1. Create handler file in `internal/<context>/controllers/`
2. Implement `Route` interface (Pattern + ServeHTTP)
3. Add to context's `module.go`: `webserver.AsRoute(controller.NewHandler)`

### Add New Event Handler

1. Define event in `internal/shared/events/` if new
2. Create handler in `internal/<context>/handlers/`
3. Implement `Handler` interface (Topic + Handle)
4. Add to context's `module.go`: `messagebus.AsHandler(handlers.NewHandler)`

### Add New Middleware

1. Create middleware function returning `func(http.Handler) http.Handler`
2. Register in `internal/app/app.go` → `middlewares()`:
   ```go
   fx.Provide(webserver.AsMiddleware(NewMyMiddleware))
   ```

### Add New Infrastructure Package

1. Create package in `pkg/<name>/`
2. Export `Module` variable with fx providers
3. Add to `internal/app/app.go` → `generalModules()`

### Add Database Migration

This project uses **Goose** for database migrations. Migrations are embedded in the binary and run automatically on startup (configurable via `APP_MIGRATIONS_RUN_ON_STARTUP=false`).

1. Create a new migration file:
   ```bash
   make migrate-create NAME=create_posts
   ```

2. Edit the generated file in `pkg/migrations/sql/`:
   ```sql
   -- +goose Up
   -- +goose StatementBegin
   CREATE TABLE posts (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       title VARCHAR(255) NOT NULL,
       created_at TIMESTAMPTZ DEFAULT NOW()
   );
   -- +goose StatementEnd

   -- +goose Down
   -- +goose StatementBegin
   DROP TABLE posts;
   -- +goose StatementEnd
   ```

3. Create matching GORM entity in `internal/<context>/entities/` (for queries only, no AutoMigrate):
   ```go
   type Post struct {
       ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
       Title     string
       CreatedAt time.Time
   }
   ```

4. Rebuild and run - migrations execute automatically on startup.

### Migration CLI Commands

Run migrations manually using the CLI:

```bash
# Run all pending migrations
./bin/gonewproject migrate up

# Rollback the last migration
./bin/gonewproject migrate down

# Show migration status
./bin/gonewproject migrate status
```

**Important:** Never use `db.AutoMigrate()`. All schema changes must go through Goose migrations.

## Code Style

- Use `log/slog` for logging (injected via fx)
- Use `encoding/json/v2` for JSON (requires GOEXPERIMENT=jsonv2)
- Handlers receive dependencies via constructor injection
- Keep handlers thin - delegate to services
- Services contain business logic
- Repositories handle data access only

## Testing

Run tests with:
```bash
make test
```

Or directly:
```bash
GOEXPERIMENT=jsonv2 go test -v ./...
```

## Verification Commands

```bash
# Build check
make build

# Run tests
make test

# Format code
make fmt

# Lint (requires golangci-lint)
make lint

# Start all infrastructure
make docker-up

# Run application
make run-serve
```
