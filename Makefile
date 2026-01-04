.PHONY: help build run test clean docker-up docker-down docker-logs deps fmt lint migrate-create

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the project with GOEXPERIMENT=jsonv2
	@echo "Building with GOEXPERIMENT=jsonv2..."
	GOEXPERIMENT=jsonv2 go build -o bin/gonewproject main.go

run: build ## Build and run the project
	@echo "Running the application..."
	./bin/gonewproject

run-serve: build ## Build and run the project with serve
	@echo "Running the application with serve..."
	./bin/gonewproject serve

test: ## Run tests
	@echo "Running tests..."
	GOEXPERIMENT=jsonv2 go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf bin/
	go clean -cache

docker-up: ## Start PostgreSQL via Docker Compose
	@echo "Starting PostgreSQL..."
	docker compose up -d

docker-down: ## Stop PostgreSQL via Docker Compose
	@echo "Stopping PostgreSQL..."
	docker compose down

docker-logs: ## Show PostgreSQL logs
	docker compose logs -f

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b \$$(go env GOPATH)/bin v1.57.2"; \
	fi

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=create_posts)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=create_posts"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	FILE="pkg/migrations/sql/$${TIMESTAMP}_$(NAME).sql"; \
	echo "-- +goose Up" > $$FILE; \
	echo "-- +goose StatementBegin" >> $$FILE; \
	echo "" >> $$FILE; \
	echo "-- +goose StatementEnd" >> $$FILE; \
	echo "" >> $$FILE; \
	echo "-- +goose Down" >> $$FILE; \
	echo "-- +goose StatementBegin" >> $$FILE; \
	echo "" >> $$FILE; \
	echo "-- +goose StatementEnd" >> $$FILE; \
	echo "Created: $$FILE"
