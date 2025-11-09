# Modsynth Examples

> Sample projects and usage examples

Part of the [Modsynth](https://github.com/modsynth) ecosystem.

## Quick Start Examples

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

---

## Production Templates 🚀

**프로덕션 준비가 완료된 풀스택 애플리케이션 템플릿**

즉시 배포 가능한 완전한 애플리케이션 템플릿으로, Docker, CI/CD, 모니터링, 테스트가 모두 포함되어 있습니다.

### 1. E-Commerce API
완전한 기능을 갖춘 전자상거래 REST API
- 상품 관리, 주문 처리
- Stripe 결제 통합
- Elasticsearch 검색
- Prometheus 모니터링

**Location**: `production-templates/e-commerce-api/`

### 2. Task Management App
Trello/Asana 스타일 작업 관리 애플리케이션
- 칸반 보드 (드래그 앤 드롭)
- 실시간 협업 (WebSocket)
- 이메일 알림
- 다국어 지원

**Location**: `production-templates/task-management-app/`

### 3. Real-Time Chat
Slack/Discord 스타일 실시간 채팅 애플리케이션
- 실시간 메시징
- 음성/영상 통화 (WebRTC)
- 메시지 검색
- 파일 공유

**Location**: `production-templates/real-time-chat/`

**[Production Templates 전체 문서 보기 →](production-templates/README.md)**

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
