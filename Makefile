.PHONY: help build run-server run-client test clean docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go server
	go build -o bin/server ./server/cmd/main.go

run-server: ## Run the Go server
	go run server/cmd/main.go

run-client: ## Run the Vue.js frontend
	cd client && npm run dev

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf client/dist/
	rm -rf client/node_modules/

docker-up: ## Start all services with Docker Compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

install-deps: ## Install dependencies
	go mod download
	cd client && npm install

dev: ## Run both server and client in development
	@echo "Starting server and client..."
	@make -j2 run-server run-client
