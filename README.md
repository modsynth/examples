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

즉시 배포 가능한 완전한 애플리케이션 템플릿으로, Docker, 테스트, 문서화가 모두 포함되어 있습니다.

### 1. E-Commerce API ✅
완전한 기능을 갖춘 전자상거래 REST API (~3,000 LOC)
- **인증**: JWT (Access + Refresh 토큰), bcrypt 암호화
- **제품 관리**: CRUD, 카테고리, 이미지 지원
- **장바구니**: 세션 기반 장바구니 관리
- **주문**: 트랜잭션 기반 주문 처리
- **아키텍처**: Clean Architecture (Repository, Service, Handler)
- **테스트**: Table-driven tests, golangci-lint
- **Docker**: Multi-stage builds, docker-compose

**Location**: `production-templates/e-commerce-api/backend/`

### 2. Task Management App ✅
Kanban 스타일 작업 관리 애플리케이션 (~6,000 LOC)
- **칸반 보드**: 프로젝트별 보드, 태스크 드래그 앤 드롭
- **실시간 협업**: WebSocket 기반 실시간 업데이트
- **역할 관리**: Owner/Admin/Member/Viewer 계층
- **태스크 기능**: 댓글, 체크리스트, 라벨, 첨부파일
- **프로젝트**: 멤버 초대, 역할 변경, 아카이브
- **WebSocket Hub**: 프로젝트 기반 룸 관리
- **Docker**: PostgreSQL + Redis + App 오케스트레이션

**Location**: `production-templates/task-management-app/backend/`

### 3. Real-Time Chat ✅
Slack 스타일 실시간 채팅 애플리케이션 (~4,500 LOC)
- **실시간 메시징**: WebSocket 기반 즉시 전송
- **룸 타입**: Direct (1:1), Group, Public 채널
- **메시지 기능**: 반응(이모지), 답장, 수정/삭제
- **읽음 확인**: 메시지별 읽음 상태 추적
- **타이핑 표시기**: 실시간 타이핑 상태
- **사용자 상태**: Online/Away/Busy/Offline
- **읽지 않은 메시지**: 룸별 카운트 관리
- **Docker**: PostgreSQL + Redis + App

**Location**: `production-templates/realtime-chat/backend/`

---

**총 구현**: ~13,500 LOC | 100+ 파일 | 3개 완전한 템플릿

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
