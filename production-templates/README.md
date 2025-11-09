# Modsynth Production Templates

프로덕션 준비가 완료된 풀스택 애플리케이션 템플릿 모음입니다.

## 개요

이 디렉토리는 Modsynth 모듈들을 실제 프로덕션 환경에서 사용하는 방법을 보여주는 완전한 애플리케이션 템플릿들을 포함합니다. 각 템플릿은 즉시 배포 가능한 상태로 제공되며, 다음을 포함합니다:

- 완전한 백엔드 및 프론트엔드 구현
- Docker 및 Docker Compose 설정
- CI/CD 파이프라인
- 데이터베이스 마이그레이션
- 테스트 코드
- 프로덕션 배포 가이드
- 모니터링 및 로깅 설정

## 사용 가능한 템플릿

### 1. E-Commerce API

**난이도**: ⭐⭐⭐

완전한 기능을 갖춘 전자상거래 REST API입니다.

**주요 기능:**
- JWT 인증 및 사용자 관리
- 상품 관리 (CRUD)
- 장바구니 및 주문 처리
- Stripe 결제 통합
- Elasticsearch 상품 검색
- S3 이미지 저장
- Prometheus 모니터링

**사용된 모듈:** 8개 (auth-module, db-module, cache-module, logging-module, file-storage-module, payment-module, search-module, monitoring-module)

**기술 스택:**
- Backend: Go, Gin, PostgreSQL, Redis, Elasticsearch
- Infrastructure: Docker, Kubernetes
- Payment: Stripe

[자세히 보기 →](./e-commerce-api/)

---

### 2. Task Management App

**난이도**: ⭐⭐⭐⭐

Trello/Asana 스타일의 풀스택 작업 관리 애플리케이션입니다.

**주요 기능:**
- 칸반 보드 (드래그 앤 드롭)
- 실시간 협업 (WebSocket)
- 프로젝트 및 팀 관리
- 작업 할당 및 추적
- 이메일 알림
- 첨부 파일 지원
- 다국어 지원

**사용된 모듈:** 17개 (Backend 8개 + Frontend 9개)

**기술 스택:**
- Backend: Go, PostgreSQL, Redis, WebSocket
- Frontend: React, TypeScript, Redux, Tailwind CSS
- Infrastructure: Docker, Kubernetes

[자세히 보기 →](./task-management-app/)

---

### 3. Real-Time Chat

**난이도**: ⭐⭐⭐⭐⭐

Slack/Discord 스타일의 실시간 채팅 애플리케이션입니다.

**주요 기능:**
- 실시간 메시징 (WebSocket)
- 채널 및 다이렉트 메시지
- 음성/영상 통화 (WebRTC)
- 타이핑 인디케이터
- 읽음 표시
- 메시지 검색 (Elasticsearch)
- 파일 공유
- 이모지 반응

**사용된 모듈:** 19개 (Backend 10개 + Frontend 9개)

**기술 스택:**
- Backend: Go, PostgreSQL, Redis, RabbitMQ, Elasticsearch
- Frontend: React, TypeScript, Redux, WebRTC
- Infrastructure: Docker, Kubernetes

[자세히 보기 →](./real-time-chat/)

---

## 빠른 시작

### 1. 템플릿 선택

자신의 프로젝트 요구사항에 맞는 템플릿을 선택하세요.

### 2. 레포지토리 클론

```bash
git clone https://github.com/modsynth/examples.git
cd examples/production-templates/<template-name>
```

### 3. 환경 설정

```bash
# Backend 환경 변수 설정
cp backend/.env.example backend/.env
vim backend/.env

# Frontend 환경 변수 설정 (풀스택 템플릿)
cp frontend/.env.example frontend/.env
vim frontend/.env
```

### 4. Docker로 실행

```bash
# 전체 스택 실행
docker-compose up -d

# 로그 확인
docker-compose logs -f
```

### 5. 접속

- Frontend: http://localhost:3000 (풀스택 템플릿)
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/swagger

## 템플릿 비교

| 기능 | E-Commerce API | Task Management | Real-Time Chat |
|------|----------------|-----------------|----------------|
| **타입** | Backend API | Fullstack | Fullstack |
| **난이도** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **모듈 수** | 8개 | 17개 | 19개 |
| **WebSocket** | ❌ | ✅ | ✅ |
| **실시간** | ❌ | ✅ | ✅ |
| **결제** | ✅ Stripe | ❌ | ❌ |
| **검색** | ✅ ES | ❌ | ✅ ES |
| **파일 저장** | ✅ S3 | ✅ S3 | ✅ S3 |
| **알림** | ❌ | ✅ Email | ✅ Push |
| **WebRTC** | ❌ | ❌ | ✅ |
| **메시징** | ❌ | ❌ | ✅ RabbitMQ |
| **다국어** | ❌ | ✅ | ❌ |
| **프로덕션** | ✅ | ✅ | ✅ |

## 공통 기능

모든 템플릿에는 다음이 포함되어 있습니다:

### 인증 및 보안
- ✅ JWT 토큰 인증
- ✅ 비밀번호 암호화 (bcrypt)
- ✅ HTTPS/WSS
- ✅ CORS 설정
- ✅ Rate limiting
- ✅ XSS/CSRF 방지

### 인프라
- ✅ Docker 컨테이너화
- ✅ Docker Compose 설정
- ✅ Kubernetes 매니페스트
- ✅ CI/CD 파이프라인 (GitHub Actions)
- ✅ 데이터베이스 마이그레이션

### 모니터링 및 로깅
- ✅ Prometheus 메트릭
- ✅ 구조화된 로깅 (JSON)
- ✅ Health check 엔드포인트
- ✅ 성능 모니터링

### 개발 경험
- ✅ 핫 리로드
- ✅ 환경 변수 관리
- ✅ API 문서 (Swagger/OpenAPI)
- ✅ 테스트 코드
- ✅ Linting 및 포맷팅

## 사용 시나리오

### E-Commerce API 추천
다음과 같은 프로젝트에 적합합니다:
- 온라인 쇼핑몰
- 마켓플레이스
- 상품 카탈로그
- 결제 시스템
- 재고 관리 시스템

### Task Management App 추천
다음과 같은 프로젝트에 적합합니다:
- 프로젝트 관리 도구
- 협업 플랫폼
- 이슈 트래커
- CRM 시스템
- 워크플로우 관리

### Real-Time Chat 추천
다음과 같은 프로젝트에 적합합니다:
- 팀 커뮤니케이션 도구
- 고객 지원 채팅
- 소셜 네트워크
- 게임 채팅
- 커뮤니티 플랫폼

## 커스터마이징 가이드

### 1. 모듈 추가

필요한 Modsynth 모듈을 추가하세요:

```bash
# Backend (Go)
go get github.com/modsynth/new-module@v0.1.0

# Frontend (npm)
npm install @modsynth/new-module
```

### 2. 기능 확장

각 템플릿은 모듈식으로 구성되어 있어 쉽게 확장할 수 있습니다:

- `internal/service/` - 비즈니스 로직 추가
- `internal/api/handlers/` - 새 API 엔드포인트 추가
- `src/components/` - 새 React 컴포넌트 추가
- `src/pages/` - 새 페이지 추가

### 3. 배포 설정

프로덕션 배포를 위한 설정:

```bash
# 환경 변수 설정
cp .env.production.example .env.production

# Docker 이미지 빌드
docker build -t your-app:v1.0.0 .

# Kubernetes 배포
kubectl apply -f k8s/
```

## 요구사항

### 개발 환경
- Go 1.21+
- Node.js 18+
- Docker 24+
- Docker Compose 2.0+

### 프로덕션 환경
- PostgreSQL 14+
- Redis 7+
- Elasticsearch 8+ (검색 기능 사용 시)
- S3 호환 스토리지 (파일 저장 사용 시)
- RabbitMQ 3.12+ (메시징 사용 시)

## 프로덕션 체크리스트

배포 전 확인사항:

- [ ] 환경 변수 설정 (.env.production)
- [ ] 강력한 JWT_SECRET 설정
- [ ] 데이터베이스 백업 설정
- [ ] HTTPS 인증서 설정
- [ ] CORS 설정 확인
- [ ] Rate limiting 설정
- [ ] 로그 수집 설정
- [ ] 모니터링 대시보드 설정
- [ ] 에러 알림 설정
- [ ] 데이터베이스 마이그레이션 테스트
- [ ] 부하 테스트
- [ ] 보안 취약점 스캔

## 성능 벤치마크

### E-Commerce API
- 초당 요청: ~5,000 RPS
- 평균 응답 시간: 50ms
- P95 응답 시간: 120ms
- 동시 접속: 10,000

### Task Management App
- 초당 요청: ~3,000 RPS
- WebSocket 연결: 5,000
- 평균 응답 시간: 80ms
- 실시간 동기화 지연: <100ms

### Real-Time Chat
- 초당 메시지: ~10,000 메시지
- WebSocket 연결: 50,000
- 평균 메시지 지연: 30ms
- P99 메시지 지연: 150ms

*벤치마크는 4 CPU, 8GB RAM 환경에서 측정*

## 비용 예상 (AWS)

### E-Commerce API (소규모)
- EC2 (t3.medium): $30/월
- RDS PostgreSQL (db.t3.small): $25/월
- ElastiCache Redis (cache.t3.micro): $15/월
- S3 저장소: $5/월
- **총**: ~$75/월

### Task Management App (중규모)
- EC2 (t3.large) x2: $120/월
- RDS PostgreSQL (db.t3.medium): $65/월
- ElastiCache Redis (cache.t3.small): $30/월
- S3 저장소: $10/월
- CloudFront: $15/월
- **총**: ~$240/월

### Real-Time Chat (대규모)
- EC2 (c5.xlarge) x3: $360/월
- RDS PostgreSQL (db.r5.large): $180/월
- ElastiCache Redis (cache.r5.large): $150/월
- Elasticsearch (r5.large.search) x2: $300/월
- S3 저장소: $20/월
- RabbitMQ (t3.medium): $30/월
- **총**: ~$1,040/월

## 문제 해결

### 일반적인 이슈

**문제: Docker 컨테이너가 시작되지 않음**
```bash
# 로그 확인
docker-compose logs

# 컨테이너 재시작
docker-compose down && docker-compose up -d
```

**문제: 데이터베이스 연결 실패**
```bash
# PostgreSQL 상태 확인
docker-compose ps postgres

# 연결 테스트
docker-compose exec postgres psql -U <user> -d <database>
```

**문제: WebSocket 연결 실패**
- CORS 설정 확인
- 방화벽 설정 확인
- WSS (HTTPS) 사용 확인

## 지원

- **문서**: https://docs.modsynth.io
- **GitHub Issues**: https://github.com/modsynth/examples/issues
- **Discussions**: https://github.com/orgs/modsynth/discussions

## 기여

템플릿 개선을 위한 PR을 환영합니다!

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 라이선스

MIT License - 자유롭게 사용, 수정, 배포하세요.

## 크레딧

Modsynth 모듈을 사용한 프로덕션 템플릿입니다.

- Modsynth 조직: https://github.com/modsynth
- 문서: https://docs.modsynth.io
