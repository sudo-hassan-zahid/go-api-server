APP_NAME := go-api-server
DOCKER_COMPOSE := docker compose
ENV ?= local
PORT ?= 8080

.PHONY: build
build:
	@echo "Starting Postgres container..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "Waiting 5s for Postgres to be ready..."
	sleep 5
	@echo "Generating swagger docs"
	swag init -g ./cmd/main.go
	sleep 3
	@echo "Starting Go server..."
	ENV=$(ENV) PORT=$(PORT) go run ./cmd/main.go


.PHONY: run
run:
	swag init -g ./cmd/main.go && go run ./cmd/main.go

.PHONY: swagger
swagger:
	swag init -g ./cmd/main.go