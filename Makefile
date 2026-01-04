APP_NAME := go-api-server
DOCKER_COMPOSE := docker compose
ENV ?= local
PORT ?= 8080

.PHONY: run
run:
	@echo "Starting Postgres container..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "Waiting 5s for Postgres to be ready..."
	sleep 5
	@echo "Starting Go server..."
	ENV=$(ENV) PORT=$(PORT) go run ./cmd/main.go
