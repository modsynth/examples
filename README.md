# Modsynth Examples

> Sample projects and usage examples

Part of the [Modsynth](https://github.com/modsynth) ecosystem.

## Examples

### 1. Backend Example
A complete Go backend using Modsynth modules:
- auth-module for authentication
- db-module for database
- cache-module for caching
- logging-module for structured logging
- api-gateway for routing

**Location**: `backend-example/`

### 2. Frontend Example
A React application using Modsynth frontend modules:
- ui-components for UI
- api-client for HTTP requests
- auth-client for authentication
- state-management for Redux
- routing for navigation

**Location**: `frontend-example/`

### 3. Fullstack Example
A complete full-stack application combining all modules:
- Go backend with all backend modules
- React frontend with all frontend modules
- Real-time WebSocket communication
- Monitoring and analytics

**Location**: `fullstack-example/`

## Getting Started

Each example has its own README with setup instructions.

```bash
# Backend example
cd backend-example
go mod download
go run main.go

# Frontend example
cd frontend-example
npm install
npm start

# Fullstack example
cd fullstack-example
docker-compose up
```

## Features Demonstrated

- **Authentication** - JWT tokens, OAuth2.0
- **Database** - GORM with PostgreSQL
- **Caching** - Redis integration
- **State Management** - Redux Toolkit
- **Form Handling** - React Hook Form + Zod
- **Real-time** - WebSocket communication
- **Monitoring** - Prometheus metrics
- **Internationalization** - Multi-language support
- **Charts & Tables** - Data visualization

## Architecture

```
Backend (Go)
├── API Gateway (Gin)
├── Auth Module (JWT + OAuth)
├── Database (GORM)
├── Cache (Redis)
├── Logging (Zap)
├── Monitoring (Prometheus)
└── WebSocket Server

Frontend (React + TypeScript)
├── UI Components (Tailwind)
├── API Client (Axios)
├── Auth Client
├── State Management (Redux)
├── Routing (React Router)
├── Form Validation (Zod)
├── Charts (Chart.js)
├── Tables (TanStack)
└── i18n (i18next)
```

## Version

Current version: `v0.1.0`

## License

MIT
