APP_NAME := objexrg
MAIN := ./cmd/main.go
BINARY := bin/$(APP_NAME)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets
	@echo Usage: make [target]
	@echo.
	@echo   build         Compile the binary into ./bin
	@echo   run           Run the API locally
	@echo   tidy          Tidy go modules
	@echo   fmt           Format all Go files
	@echo   vet           Run go vet
	@echo   test          Run all tests
	@echo   check         Format, vet and test
	@echo   swag          Generate Swagger docs into ./cmd/docs
	@echo   docker-build  Build the Docker image
	@echo   docker-up     Start the Docker stack
	@echo   docker-down   Stop the Docker stack
	@echo   clean         Remove build artifacts

.PHONY: tidy
tidy: ## Tidy go modules
	go mod tidy

.PHONY: build
build: ## Compile the binary into ./bin
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
docker-build: ## Build the Docker image
	docker compose build

.PHONY: docker-up
docker-up: ## Start the Docker stack
	docker compose up --build

.PHONY: docker-down
docker-down: ## Stop the Docker stack
	docker compose down

.PHONY: clean
clean: ## Remove build artifacts
	-rm -rf bin
	go clean -cache -testcache
