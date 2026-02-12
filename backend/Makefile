.PHONY: help run build test clean migrate-up migrate-status migrate-verify seed db-shell db-reset docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	go run cmd/api/main.go

build: ## Build the application
	go build -o bin/api cmd/api/main.go

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Download dependencies
	go mod download
	go mod tidy

migrate-up: ## Run database migrations up
	@go run ./cmd/migrate ./migrations

migrate-status: ## Check database tables
	@psql -h localhost -U postgres -d customermx -c "\dt"

migrate-verify: ## Verify seed data
	@psql -h localhost -U postgres -d customermx -c "SELECT b.name, COUNT(v.id) as vehicle_count FROM brands b LEFT JOIN vehicles v ON b.id = v.brand_id GROUP BY b.name;"

seed: ## Seed initial admin user
	@go run ./cmd/seed

db-shell: ## Open PostgreSQL shell
	@psql -h localhost -U postgres -d customermx

db-reset: ## Drop and recreate database (WARNING: destructive)
	@echo "⚠️  This will delete ALL data. Press Ctrl+C to cancel, Enter to continue..."
	@read confirm
	@psql -h localhost -U postgres -d postgres -c "DROP DATABASE IF EXISTS customermx;"
	@psql -h localhost -U postgres -d postgres -c "CREATE DATABASE customermx;"
	@echo "Database reset complete. Run 'make migrate-up' to apply migrations."

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

install-tools: ## Install development tools
	go install github.com/air-verse/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: ## Run linter
	golangci-lint run

watch: ## Run with hot reload (requires air)
	air
