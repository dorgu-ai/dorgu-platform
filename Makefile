.PHONY: help build run test clean frontend lint

# Variables
BINARY_NAME=dorgu-platform
SERVER_BINARY=bin/server
GO_FILES=$(shell find . -name '*.go' -type f -not -path "./vendor/*")
WEB_DIR=web
WEB_DIST=$(WEB_DIR)/dist

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: frontend ## Build backend with embedded frontend
	@echo "Building backend..."
	@mkdir -p bin
	go build -o $(SERVER_BINARY) ./cmd/server
	@echo "Build complete: $(SERVER_BINARY)"

run: ## Run backend in development mode
	@echo "Starting backend server..."
	go run ./cmd/server/main.go

test: ## Run Go tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

lint: ## Run linters
	@echo "Running linters..."
	golangci-lint run ./...

frontend: ## Build frontend
	@echo "Building frontend..."
	cd $(WEB_DIR) && npm install && npm run build
	@echo "Frontend build complete: $(WEB_DIST)"

frontend-dev: ## Run frontend in dev mode
	@echo "Starting frontend dev server..."
	cd $(WEB_DIR) && npm run dev

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf $(WEB_DIST)
	rm -f coverage.txt
	@echo "Clean complete"

deps: ## Download Go dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
