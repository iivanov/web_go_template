# Go Project Template

A modular Go application template demonstrating clean architecture with bounded contexts, dependency injection, event-driven communication, and observability.

## Tech Stack

- **Go 1.25+** with `GOEXPERIMENT=jsonv2` (encoding/json/v2)
- **Uber fx** - Dependency injection framework
- **Cobra** - CLI framework
- **Viper** - Configuration management
- **GORM** - ORM with PostgreSQL
- **Watermill** - Message bus abstraction
- **OpenTelemetry** - Distributed tracing and metrics
- **Docker Compose** - Local infrastructure (PostgreSQL, Jaeger, Prometheus, Grafana)

## Project Structure

```
.
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command with global flags
│   └── serve.go                  # HTTP server command
├── internal/                     # Private application code
│   ├── app/                      # Application bootstrap and config
│   │   ├── app.go                # fx module composition
│   │   └── config.go             # Configuration loading (Viper)
│   ├── shared/                   # Shared code across bounded contexts
│   │   └── events/               # Domain events definitions
│   ├── someboundedcontext/       # Example bounded context
│   │   ├── config/               # Context-specific config
│   │   ├── controllers/          # HTTP handlers
│   │   ├── dto/                  # Data transfer objects
│   │   ├── entities/             # Domain entities
│   │   ├── repositories/         # Data access layer
│   │   ├── services/             # Business logic
│   │   └── module.go             # fx module definition
│   └── secondboundedcontext/     # Another bounded context
│       ├── handlers/             # Message bus handlers
│       ├── services/             # Business logic
│       └── module.go             # fx module definition
├── pkg/                          # Reusable packages
│   ├── database/                 # GORM database setup
│   ├── logger/                   # Structured logging (slog)
│   ├── messagebus/               # Event bus abstraction (Watermill)
│   ├── telemetry/                # OpenTelemetry setup
│   └── webserver/                # HTTP server with routing
└── main.go                       # Entry point
```

## Architecture Overview

### Dependency Injection with fx

The application uses **Uber fx** for dependency injection. Each package exposes a `Module` variable:

```go
var Module = fx.Module("someboundedcontext",
    fx.Provide(
        services.NewUserService,
        repositories.NewUserRepository,
        webserver.AsRoute(controller.NewUserHandler),
    ),
)
```

Modules are composed in `internal/app/app.go`:

```go
func Serve(configFile string) {
    fx.New(
        fx.Options(generalModules()...),
        fx.Options(middlewares()...),
        fx.Provide(NewServeConfig(configFile)),
        webserver.Module,
    ).Run()
}
```

### Bounded Contexts

Business logic is organized into **bounded contexts** under `internal/`. Each context is self-contained with its own:
- **module.go** - fx module registering all dependencies
- **controllers/** - HTTP handlers
- **services/** - Business logic
- **repositories/** - Data access
- **entities/** - Domain models
- **dto/** - Request/response objects

### HTTP Routing

Routes implement the `Route` interface and are auto-registered via fx groups:

```go
type Route interface {
    http.Handler
    Pattern() string
}

// Register a route
webserver.AsRoute(controller.NewUserHandler)
```

Handler example:
```go
func (*UserHandler) Pattern() string {
    return "GET /api/users/{id}"
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Handle request
}
```

### Event-Driven Communication

Bounded contexts communicate via the message bus (Watermill-based):

**Publishing events:**
```go
publisher.Publish(ctx, events.UserCreatedEvent{
    UserID: user.ID,
    Name:   user.Name,
})
```

**Subscribing to events:**
```go
// Register handler in module.go
messagebus.AsHandler(handlers.NewUserCreatedHandler)

// Implement Handler interface
func (h *UserCreatedHandler) Topic() string {
    return events.TopicUserCreated
}

func (h *UserCreatedHandler) Handle(ctx context.Context, payload []byte) error {
    // Process event
}
```

### Configuration

Configuration is loaded via Viper with support for:
- YAML config files (`--config` flag)
- Environment variables (prefix: `APP_`)
- `.env` files

```bash
# Via config file
./bin/gonewproject serve --config config.yaml

# Via environment variables
APP_WEBSERVER_PORT=8080 ./bin/gonewproject serve
```

## Quick Start

### Prerequisites
- Go 1.25+
- Docker & Docker Compose

### Setup

1. **Start infrastructure:**
```bash
make docker-up
```

2. **Copy environment file:**
```bash
cp .env.example .env
```

3. **Build and run:**
```bash
make run-serve
```

### Available Services

| Service    | URL                    | Description           |
|------------|------------------------|-----------------------|
| App        | http://localhost:8080  | Application API       |
| Jaeger     | http://localhost:16686 | Distributed tracing   |
| Prometheus | http://localhost:9090  | Metrics               |
| Grafana    | http://localhost:3000  | Dashboards (admin/admin) |
| PostgreSQL | localhost:5433         | Database              |

## Makefile Commands

| Command          | Description                              |
|------------------|------------------------------------------|
| `make help`      | Show all available commands              |
| `make build`     | Build with GOEXPERIMENT=jsonv2           |
| `make run`       | Build and run                            |
| `make run-serve` | Build and run HTTP server                |
| `make test`      | Run tests                                |
| `make fmt`       | Format code                              |
| `make lint`      | Run golangci-lint                        |
| `make deps`      | Download and tidy dependencies           |
| `make clean`     | Remove build artifacts                   |
| `make docker-up` | Start Docker Compose services            |
| `make docker-down` | Stop Docker Compose services           |
| `make docker-logs` | View Docker Compose logs               |

## Development Guide

### Adding a New Bounded Context

1. Create directory structure:
```
internal/newcontext/
├── module.go
├── controllers/
├── services/
├── repositories/
├── entities/
└── dto/
```

2. Define the fx module in `module.go`:
```go
var Module = fx.Module("newcontext",
    fx.Provide(
        services.NewMyService,
        webserver.AsRoute(controllers.NewMyHandler),
    ),
)
```

3. Register in `internal/app/app.go`:
```go
func generalModules() []fx.Option {
    return []fx.Option{
        // ...existing modules
        newcontext.Module,
    }
}
```

### Adding a New HTTP Endpoint

1. Create handler in `controllers/`:
```go
type MyHandler struct {
    logger  *slog.Logger
    service *services.MyService
}

func NewMyHandler(logger *slog.Logger, service *services.MyService) *MyHandler {
    return &MyHandler{logger: logger, service: service}
}

func (*MyHandler) Pattern() string {
    return "POST /api/resource"
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. Register in module:
```go
webserver.AsRoute(controllers.NewMyHandler)
```

### Adding a New Event Handler

1. Define event in `internal/shared/events/`:
```go
const TopicOrderCreated = "order.created"

type OrderCreatedEvent struct {
    OrderID uuid.UUID `json:"order_id"`
}

func (e OrderCreatedEvent) Topic() string {
    return TopicOrderCreated
}
```

2. Create handler:
```go
type OrderCreatedHandler struct {
    logger *slog.Logger
}

func NewOrderCreatedHandler(logger *slog.Logger) messagebus.Handler {
    return &OrderCreatedHandler{logger: logger}
}

func (h *OrderCreatedHandler) Topic() string {
    return events.TopicOrderCreated
}

func (h *OrderCreatedHandler) Handle(ctx context.Context, payload []byte) error {
    // Process event
}
```

3. Register in module:
```go
messagebus.AsHandler(handlers.NewOrderCreatedHandler)
```

### Adding HTTP Middleware

1. Create middleware function:
```go
func NewMyMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Before
            next.ServeHTTP(w, r)
            // After
        })
    }
}
```

2. Register in `internal/app/app.go`:
```go
func middlewares() []fx.Option {
    return []fx.Option{
        fx.Provide(webserver.AsMiddleware(NewMyMiddleware)),
    }
}
```

## Build Requirements

This project uses Go's experimental JSON v2 encoder. Always build with:
```bash
GOEXPERIMENT=jsonv2 go build ...
```

The Makefile handles this automatically.
