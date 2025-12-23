# Mighty Eagle - Makefile

.PHONY: help dev build test migrate-up migrate-down clean docker-up docker-down

help: ## Show this help message
	@echo "Mighty Eagle - Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev: ## Start development server
	cd services/api-go && go run main.go

build: ## Build production binary
	cd services/api-go && go build -o ../../bin/mighty-eagle-api main.go

test: ## Run tests
	cd services/api-go && go test ./... -v -cover

migrate-up: ## Run database migrations
	psql -h localhost -U postgres -d mighty_eagle -f services/api-go/migrations/001_initial_schema.sql

migrate-down: ## Rollback database migrations
	@echo "Manual rollback required - drop and recreate database"

clean: ## Clean build artifacts
	rm -rf bin/
	cd services/api-go && go clean

docker-up: ## Start Docker Compose services
	docker compose up -d

docker-down: ## Stop Docker Compose services
	docker compose down

docker-logs: ## View Docker Compose logs
	docker compose logs -f

docker-rebuild: ## Rebuild and restart Docker services
	docker compose down
	docker compose up -d --build

install-deps: ## Install Go dependencies
	cd services/api-go && go mod download && go mod tidy

fmt: ## Format Go code
	cd services/api-go && go fmt ./...

lint: ## Run linter
	cd services/api-go && golangci-lint run

.DEFAULT_GOAL := help
