# go-api-server

A high performance API server built with Go and the Fiber framework. This project provides a robust foundation for building scalable web services with features like authentication, role based access control, and rate limiting.

## Project Overview

This server is designed to handle common API requirements out of the box. It uses GORM for database interactions, Redis for caching and rate limiting, and provides automatic Swagger documentation generation.

## Core Features

- User Authentication: Secure signup, login, and password reset flows using JWT.
- Email Verification: Integrated support for verifying user email addresses.
- Role Based Access Control (RBAC): Permission management with defined roles like Admin and User.
- Rate Limiting: Protection against brute force and DDoS attacks on auth and public routes.
- Database Management: Automatic migrations and efficient querying with GORM.
- API Documentation: Interactive Swagger UI for exploring and testing endpoints.
- Health Monitoring: Dedicated endpoints for server and database health status.
- Environment Configuration: Flexible setup using environment variables and .env files.
- Containerization: Full Docker support for easy deployment and local development.

## Technologies Used

- Go (v1.25.0)
- Fiber (Web Framework)
- PostgreSQL (Primary Database)
- Redis (Caching and Rate Limiting)
- GORM (ORM)
- JWT (Authentication)
- Swagger (API Documentation)
- Docker and Docker Compose
- Makefile (Task Automation)

## Prerequisites

Ensure you have the following installed on your system:

- Go v1.25.0 or later
- Docker and Docker Compose
- Make (optional but recommended)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/sudo-hassan-zahid/go-api-server.git
cd go-api-server
```

### 2. Environment Setup

The repository includes a Docker-ready `.env.example` file. Copy it to `.env` and update values if you need different local credentials, ports, or secrets.

```bash
cp .env.example .env
```

### 3. Running with Docker

Start the full Docker stack, including the Go app, Postgres, and Redis:

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080`.

Docker Compose runs this as one application named `go-api-server` with three containers:

- `go-api-server`: the Go API container.
- `go-api-server-postgres`: the private Postgres container.
- `go-api-server-redis`: the private Redis container.

Postgres is only available inside the Docker Compose network at `postgres:5432`; it is not published to your host machine, so it will not collide with any local database ports. Redis is also kept inside the Compose network at `redis:6379`.

The app container includes a healthcheck against `/api/health/server`.

Stop the stack:

```bash
docker compose down
```

## API Documentation

Swagger documentation is automatically generated. Once the server is running, you can access the interactive UI at:

`http://localhost:8080/swagger/index.html`

To regenerate the documentation after making changes to the API definitions, run:

```bash
make swagger
```

## Project Structure

- /cmd: Entry point of the application.
- /internal: Private application code (config, handler, service, repository, middleware).
- /routes: API route definitions and group setup.
- /docs: Generated Swagger documentation.
- /utils: Common utility functions and helpers.
- /tests: Unit and integration tests.

## Testing

Run the test suite using the standard Go test command:

```bash
go test ./...
```
