# Variables
DOCKER_COMPOSE=docker-compose
DOCKER_COMPOSE_FILE=deployments/docker-compose.yaml

.PHONY: run
run: prune ## Build and run the application using Docker Compose
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) up --build

.PHONY: stop
stop: ## Stop all running containers
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down

.PHONY: clean
clean: ## Remove all stopped containers
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down --rmi all --volumes

.PHONY: prune
prune: ## Remove all containers and volumes
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_FILE) down -v

.PHONY: test
test: ## Run unit tests
	go test ./...

.PHONY: proto
proto: ## Generate gRPC code from .proto files
	protoc -I=api/proto --go_out=api/proto/pb --go-grpc_out=api/proto/pb api/proto/*.proto

.PHONY: build
build: ## Build the Go binaries
	go build -o bin/antibruteforce ./cmd/antibruteforce

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest

.PHONY: lint
lint: ## Run linters
	golangci-lint run

.PHONY: format
format: ## Format the code
	go fmt ./...

.PHONY: ci
ci: format lint test ## Run a continuous integration pipeline (format, lint, test)
