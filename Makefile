APP_NAME := go-api-server
DOCKER_COMPOSE := docker compose

.PHONY: docker-up
docker-up:
	$(DOCKER_COMPOSE) up --build

.PHONY: docker-down
docker-down:
	$(DOCKER_COMPOSE) down

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: test
test:
	go test ./...

.PHONY: swagger
swagger:
	swag init -g ./cmd/main.go
