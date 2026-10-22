APP_NAME := objexrg
MAIN := ./cmd/main.go
BINARY := bin/$(APP_NAME)
FRONTEND_DIR := frontend
BACKEND_REPLICAS ?= 2
FRONTEND_REPLICAS ?= 2
SCALE := --scale objexrg=$(BACKEND_REPLICAS) --scale frontend=$(FRONTEND_REPLICAS)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets
	@echo Usage: make [target]
	@echo.
	@echo   up            Build and start infra + backend + frontend behind nginx
	@echo   down          Stop and remove the Docker stack
	@echo   restart       Restart the whole Docker stack
	@echo   logs          Follow logs from every service
	@echo   ps            Show the status of every service
	@echo   images        Build the backend and frontend images
	@echo   infra         Start only infra (postgres + rustfs + bucket init)
	@echo   backend       Rebuild backend replicas (stack must be up)
	@echo   frontend      Rebuild frontend replicas (stack must be up)
	@echo   scale         Rescale replicas (gateway picks them up dynamically)
	@echo   build         Compile the Go binary into ./bin
	@echo   run           Run the API locally
	@echo   dev-frontend  Run the frontend dev server (Vite)
	@echo   tidy          Tidy go modules
	@echo   fmt           Format all Go files
	@echo   vet           Run go vet
	@echo   test          Run all tests
	@echo   check         Format, vet and test
	@echo   swag          Generate Swagger docs into ./cmd/docs
	@echo   clean         Remove build artifacts

.PHONY: up
up: ## Build and start infra + backend + frontend behind nginx
	docker compose up -d --build $(SCALE)
	@echo.
	@echo App:      http://localhost:8080
	@echo Swagger:  http://localhost:8080/swagger/index.html
	@echo Backend:  $(BACKEND_REPLICAS) replicas behind nginx
	@echo Frontend: $(FRONTEND_REPLICAS) replicas behind nginx
	@echo Run "make logs" to follow the logs.

.PHONY: down
down: ## Stop and remove the Docker stack
	docker compose down

.PHONY: restart
restart: down up ## Restart the whole Docker stack

.PHONY: logs
logs: ## Follow logs from every service
	docker compose logs -f

.PHONY: ps
ps: ## Show the status of every service
	docker compose ps

.PHONY: images
images: ## Build the backend and frontend images
	docker compose build

.PHONY: infra
infra: ## Start only infra (postgres + rustfs + bucket init)
	docker compose up -d postgres rustfs rustfs-init

.PHONY: backend
backend: ## Build and start only the backend replicas
	docker compose up -d --build --scale objexrg=$(BACKEND_REPLICAS) objexrg

.PHONY: frontend
frontend: ## Build and start only the frontend replicas
	docker compose up -d --build --scale frontend=$(FRONTEND_REPLICAS) frontend

.PHONY: scale
scale: ## Rescale backend/frontend replicas (gateway picks them up dynamically)
	docker compose up -d $(SCALE)

.PHONY: dev-frontend
dev-frontend: ## Run the frontend dev server (Vite)
	cd $(FRONTEND_DIR) && bun run dev

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: build
build: ## Compile the Go binary into ./bin
	go build -o $(BINARY) $(MAIN)

.PHONY: run
run: ## Run the API locally
	go run $(MAIN)

.PHONY: fmt
fmt: ## Format all Go files
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: test
test: ## Run all tests
	go test ./...

.PHONY: check
check: fmt vet test ## Format, vet and test

.PHONY: swag
swag: ## Generate Swagger docs into ./cmd/docs
	swag init -g $(MAIN) -o cmd/docs

.PHONY: docker-build
docker-build: images ## Alias for images

.PHONY: docker-up
docker-up: up ## Alias for up

.PHONY: docker-down
docker-down: down ## Alias for down

.PHONY: clean
clean: ## Remove build artifacts
	-rm -rf bin
	go clean -cache -testcache
