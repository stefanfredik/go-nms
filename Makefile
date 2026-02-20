.PHONY: help build test run clean docker migrate lint dev dev-infra dev-down dev-logs dev-status install-air

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_DIR=bin
GOPATH ?= $(shell go env GOPATH)
SERVICES=api-gateway device-service collector-service worker config-service alert-service notification-service analytics-service

# =============================================================================
# DEVELOPMENT ENVIRONMENT
# =============================================================================

dev-infra: ## Jalankan infrastruktur (postgres, redis, nats, influxdb) di Docker
	@echo "🚀 Menjalankan infrastruktur development..."
	docker compose -f docker-compose.dev.yml up -d
	@echo "✅ Infrastruktur siap:"
	@echo "   PostgreSQL : localhost:5499"
	@echo "   Redis      : localhost:6399"
	@echo "   NATS       : localhost:4299"
	@echo "   InfluxDB   : localhost:8099"
	@echo ""
	@echo "Selanjutnya jalankan: make dev"

dev: ## Jalankan api-gateway dengan hot-reload (air) - butuh `make dev-infra` terlebih dahulu
	@echo "🔥 Menjalankan api-gateway dengan hot-reload..."
	@$(GOPATH)/bin/air

dev-down: ## Hentikan dan hapus infrastruktur development
	@echo "🛑 Menghentikan infrastruktur..."
	docker compose -f docker-compose.dev.yml down

dev-logs: ## Tampilkan log infrastruktur development
	docker compose -f docker-compose.dev.yml logs -f

dev-status: ## Tampilkan status container development
	docker compose -f docker-compose.dev.yml ps

install-air: ## Install air (hot-reload tool untuk Go)
	@echo "📦 Menginstall air..."
	go install github.com/air-verse/air@latest
	@echo "✅ air berhasil diinstall di $(GOPATH)/bin/air"

# Help command
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build all services
build: ## Build all services
	@echo "Building all services..."
	@mkdir -p $(BINARY_DIR)
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build -o $(BINARY_DIR)/$$service ./cmd/$$service; \
	done
	@echo "Build complete!"

# Build specific service
build-%: ## Build specific service (e.g., make build-api-gateway)
	@echo "Building $*..."
	@mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$* ./cmd/$*

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test -v -race ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run specific service tests
test-%: ## Run tests for specific package (e.g., make test-device)
	@echo "Running tests for $*..."
	go test -v -race ./internal/$*/...

# Run linter
lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run ./...

# Format code
fmt: ## Format code with go fmt
	@echo "Formatting code..."
	go fmt ./...

# Run go mod tidy
tidy: ## Run go mod tidy
	@echo "Tidying modules..."
	go mod tidy

# Download dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

# Generate mocks
generate-mocks: ## Generate test mocks
	@echo "Generating mocks..."
	@go install github.com/golang/mock/mockgen@latest
	mockgen -source=internal/device/repository/device_repo.go -destination=internal/device/repository/mock/device_repo_mock.go
	@echo "Mocks generated!"

# Database migrations
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path scripts/migrations -database "postgres://nms_user:password@localhost:5432/nms_dev?sslmode=disable" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path scripts/migrations -database "postgres://nms_user:password@localhost:5432/nms_dev?sslmode=disable" down

migrate-create: ## Create new migration (usage: make migrate-create name=create_users_table)
	@echo "Creating migration $(name)..."
	migrate create -ext sql -dir scripts/migrations -seq $(name)

# Docker commands
docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $$service image..."; \
		docker build -f deployments/docker/Dockerfile.$$service -t nms/$$service:latest .; \
	done

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images..."
	@for service in $(SERVICES); do \
		docker push nms/$$service:latest; \
	done

# Docker Compose commands
compose-up: ## Start services with docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.dev.yml up -d

compose-down: ## Stop services with docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.dev.yml down

compose-logs: ## Show docker-compose logs
	docker-compose -f deployments/docker-compose/docker-compose.dev.yml logs -f

# Run development environment
run-dev: compose-up ## Start development environment
	@echo "Development environment started!"
	@echo "PostgreSQL: localhost:5432"
	@echo "InfluxDB: localhost:8086"
	@echo "Redis: localhost:6379"
	@echo "NATS: localhost:4222"

# Run specific service
run-%: ## Run specific service (e.g., make run-api-gateway)
	@echo "Running $*..."
	go run ./cmd/$*/main.go

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean complete!"

# Install tools
install-tools: ## Install development tools
	@echo "Installing tools..."
	go install github.com/golang/mock/mockgen@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Tools installed!"

# Generate API documentation
generate-docs: ## Generate API documentation
	@echo "Generating API documentation..."
	swag init -g cmd/api-gateway/main.go -o api/swagger
	@echo "Documentation generated!"

# Benchmark tests
bench: ## Run benchmark tests
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Check for vulnerabilities
security: ## Run security checks
	@echo "Running security checks..."
	gosec ./...

# Full check before commit
pre-commit: fmt lint test ## Run all checks before commit
	@echo "All checks passed!"
