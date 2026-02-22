.PHONY: help dev build test lint migrate-up migrate-down seed docker-up docker-down docker-logs frontend-dev frontend-build

# ─── Variables ────────────────────────────────────────────────────────────────
BACKEND_DIR  := ./backend
FRONTEND_DIR := ./frontend
DB_URL       ?= postgres://postgres:postgres@localhost:5432/cake_shop?sslmode=disable
BINARY       := cake-shop-api

# ─── Help ─────────────────────────────────────────────────────────────────────
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ──────────────────────────────────────────────────────────────
dev: ## Run backend in development mode (requires air or go run)
	cd $(BACKEND_DIR) && go run ./cmd/api

frontend-dev: ## Run frontend dev server
	cd $(FRONTEND_DIR) && npm run dev

# ─── Build ────────────────────────────────────────────────────────────────────
build: ## Build the backend binary
	cd $(BACKEND_DIR) && \
		CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY) ./cmd/api
	@echo "Binary built: $(BACKEND_DIR)/bin/$(BINARY)"

frontend-build: ## Build the frontend for production
	cd $(FRONTEND_DIR) && npm run build

# ─── Testing ─────────────────────────────────────────────────────────────────
test: ## Run backend unit tests
	cd $(BACKEND_DIR) && go test -v -race ./...

test-coverage: ## Run backend tests with coverage report
	cd $(BACKEND_DIR) && go test -coverprofile=coverage.out ./... && \
		go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: $(BACKEND_DIR)/coverage.html"

# ─── Code quality ────────────────────────────────────────────────────────────
lint: ## Run golangci-lint
	cd $(BACKEND_DIR) && golangci-lint run ./...

vet: ## Run go vet
	cd $(BACKEND_DIR) && go vet ./...

# ─── Database ────────────────────────────────────────────────────────────────
migrate-up: ## Apply all pending migrations
	migrate -path $(BACKEND_DIR)/db/migrations -database "$(DB_URL)" up

migrate-down: ## Rollback last migration
	migrate -path $(BACKEND_DIR)/db/migrations -database "$(DB_URL)" down 1

migrate-drop: ## Drop all tables (DANGEROUS)
	migrate -path $(BACKEND_DIR)/db/migrations -database "$(DB_URL)" drop -f

migrate-status: ## Show migration status
	migrate -path $(BACKEND_DIR)/db/migrations -database "$(DB_URL)" version

seed: ## Seed the database with sample data
	cd $(BACKEND_DIR) && go run ./db/seed/seed.go

# ─── sqlc ─────────────────────────────────────────────────────────────────────
sqlc-generate: ## Regenerate sqlc code
	cd $(BACKEND_DIR) && sqlc generate

sqlc-verify: ## Verify sqlc configuration
	cd $(BACKEND_DIR) && sqlc verify

# ─── Docker ──────────────────────────────────────────────────────────────────
docker-up: ## Start all services (PostgreSQL + API)
	docker compose up -d --build

docker-down: ## Stop all services
	docker compose down

docker-logs: ## Tail API logs
	docker compose logs -f api

docker-db: ## Start only PostgreSQL
	docker compose up -d postgres

docker-clean: ## Remove containers and volumes
	docker compose down -v --remove-orphans

# ─── Dependencies ────────────────────────────────────────────────────────────
deps-backend: ## Download Go dependencies
	cd $(BACKEND_DIR) && go mod download && go mod tidy

deps-frontend: ## Install frontend dependencies
	cd $(FRONTEND_DIR) && npm install

deps: deps-backend deps-frontend ## Install all dependencies

# ─── Setup ───────────────────────────────────────────────────────────────────
setup: ## Full project setup (install deps, start DB, run migrations, seed)
	@echo "==> Installing dependencies..."
	$(MAKE) deps
	@echo "==> Starting PostgreSQL..."
	$(MAKE) docker-db
	@echo "==> Waiting for database..."
	sleep 5
	@echo "==> Running migrations..."
	$(MAKE) migrate-up
	@echo "==> Seeding database..."
	$(MAKE) seed
	@echo "==> Setup complete! Run 'make dev' and 'make frontend-dev' in separate terminals."
