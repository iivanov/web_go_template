.PHONY: help build run test clean docker-up docker-down docker-logs deps fmt lint

# Default target
help:
	@echo "Available commands:"
	@echo "  make build     - Build the project with GOEXPERIMENT=jsonv2"
	@echo "  make run       - Build and run the project"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make docker-up - Start PostgreSQL via Docker Compose"
	@echo "  make docker-down - Stop PostgreSQL via Docker Compose"
	@echo "  make docker-logs - Show PostgreSQL logs"
	@echo "  make deps      - Download dependencies"
	@echo "  make fmt       - Format Go code"
	@echo "  make lint      - Run linter"

# Build the project with GOEXPERIMENT=jsonv2
build:
	@echo "Building with GOEXPERIMENT=jsonv2..."
	GOEXPERIMENT=jsonv2 go build -o bin/gonewproject main.go

# Build and run the project
run: build
	@echo "Running the application..."
	./bin/gonewproject

# Run tests
test:
	@echo "Running tests..."
	GOEXPERIMENT=jsonv2 go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	go clean -cache

# Start PostgreSQL via Docker Compose
docker-up:
	@echo "Starting PostgreSQL..."
	docker-compose up -d postgres

# Stop PostgreSQL via Docker Compose
docker-down:
	@echo "Stopping PostgreSQL..."
	docker-compose down

# Show PostgreSQL logs
docker-logs:
	docker-compose logs -f postgres

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Format Go code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.57.2"; \
	fi
