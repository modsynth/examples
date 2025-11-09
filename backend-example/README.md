# Backend Example

Complete Go backend using Modsynth modules.

## Modules Used

- auth-module
- db-module
- cache-module
- logging-module
- api-gateway
- monitoring-module

## Setup

```bash
go mod download
go run main.go
```

Server runs on http://localhost:8080

## Endpoints

- GET /health - Health check
- POST /auth/login - User login
- GET /users - List users (authenticated)
- GET /metrics - Prometheus metrics
