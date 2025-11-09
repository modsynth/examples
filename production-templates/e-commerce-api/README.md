# E-Commerce API - Production Template

프로덕션 준비가 완료된 전자상거래 REST API 템플릿입니다.

## 사용된 Modsynth 모듈

### Backend
- **auth-module** - JWT 인증 및 사용자 관리
- **db-module** - PostgreSQL 데이터베이스 연결
- **cache-module** - Redis 캐싱
- **logging-module** - 구조화된 로깅
- **file-storage-module** - 상품 이미지 저장 (S3)
- **payment-module** - Stripe 결제 처리
- **search-module** - Elasticsearch 상품 검색
- **monitoring-module** - Prometheus 메트릭

## 기능

### 인증 및 사용자 관리
- JWT 기반 인증
- 사용자 등록/로그인
- 비밀번호 암호화
- 리프레시 토큰

### 상품 관리
- 상품 CRUD
- 카테고리 관리
- 재고 관리
- 이미지 업로드
- 전문 검색 (Elasticsearch)

### 주문 처리
- 장바구니 관리
- 주문 생성
- 결제 처리 (Stripe)
- 주문 상태 추적
- 주문 내역 조회

### 관리자 기능
- 상품 관리
- 주문 관리
- 사용자 관리
- 대시보드 통계

## 아키텍처

```
e-commerce-api/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP handlers
│   │   │   ├── auth.go
│   │   │   ├── products.go
│   │   │   ├── orders.go
│   │   │   └── admin.go
│   │   ├── middleware/          # Middleware
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── ratelimit.go
│   │   └── routes.go            # Route definitions
│   ├── domain/
│   │   ├── user.go              # Domain models
│   │   ├── product.go
│   │   ├── order.go
│   │   └── cart.go
│   ├── repository/              # Data access
│   │   ├── user_repo.go
│   │   ├── product_repo.go
│   │   └── order_repo.go
│   ├── service/                 # Business logic
│   │   ├── auth_service.go
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   └── payment_service.go
│   └── config/
│       └── config.go            # Configuration
├── migrations/                  # Database migrations
│   ├── 001_create_users.sql
│   ├── 002_create_products.sql
│   └── 003_create_orders.sql
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

## 설치 및 실행

### 1. 필수 요구사항

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Elasticsearch 8+ (선택사항)
- S3 호환 스토리지 (선택사항)
- Stripe 계정 (결제 기능 사용 시)

### 2. 환경 설정

```bash
# .env 파일 생성
cp .env.example .env

# 환경 변수 설정
vim .env
```

필수 환경 변수:
```env
# Server
PORT=8080
ENV=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=ecommerce
DB_PASSWORD=your-password
DB_NAME=ecommerce_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key-change-this
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d

# Stripe
STRIPE_SECRET_KEY=sk_test_...
STRIPE_WEBHOOK_SECRET=whsec_...

# S3 (선택사항)
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=ecommerce-images

# Elasticsearch (선택사항)
ES_ADDRESSES=http://localhost:9200
```

### 3. Docker로 실행

```bash
# 의존성 서비스 실행
docker-compose up -d postgres redis elasticsearch

# 데이터베이스 마이그레이션
go run cmd/migrate/main.go up

# 애플리케이션 실행
docker-compose up app
```

### 4. 로컬 개발

```bash
# 의존성 설치
go mod download

# 데이터베이스 마이그레이션
make migrate-up

# 개발 서버 실행
make dev

# 또는
go run cmd/server/main.go
```

## API 엔드포인트

### 인증
```
POST   /api/v1/auth/register        # 회원가입
POST   /api/v1/auth/login           # 로그인
POST   /api/v1/auth/refresh         # 토큰 갱신
POST   /api/v1/auth/logout          # 로그아웃
GET    /api/v1/auth/me              # 내 정보
```

### 상품
```
GET    /api/v1/products             # 상품 목록
GET    /api/v1/products/:id         # 상품 상세
POST   /api/v1/products             # 상품 생성 (관리자)
PUT    /api/v1/products/:id         # 상품 수정 (관리자)
DELETE /api/v1/products/:id         # 상품 삭제 (관리자)
GET    /api/v1/products/search      # 상품 검색
```

### 장바구니
```
GET    /api/v1/cart                 # 장바구니 조회
POST   /api/v1/cart/items           # 상품 추가
PUT    /api/v1/cart/items/:id       # 수량 변경
DELETE /api/v1/cart/items/:id       # 상품 제거
DELETE /api/v1/cart                 # 장바구니 비우기
```

### 주문
```
POST   /api/v1/orders               # 주문 생성
GET    /api/v1/orders               # 내 주문 목록
GET    /api/v1/orders/:id           # 주문 상세
PUT    /api/v1/orders/:id/cancel    # 주문 취소
```

### 결제
```
POST   /api/v1/payments/create-intent    # 결제 의도 생성
POST   /api/v1/payments/webhook          # Stripe 웹훅
```

### 관리자
```
GET    /api/v1/admin/orders         # 모든 주문 관리
PUT    /api/v1/admin/orders/:id     # 주문 상태 변경
GET    /api/v1/admin/stats          # 대시보드 통계
GET    /api/v1/admin/users          # 사용자 관리
```

## 테스트

```bash
# 전체 테스트
make test

# 커버리지 확인
make test-coverage

# 통합 테스트
make test-integration
```

## 배포

### Docker 배포

```bash
# 이미지 빌드
docker build -t ecommerce-api:v1.0.0 .

# 이미지 푸시
docker push your-registry/ecommerce-api:v1.0.0

# Kubernetes 배포
kubectl apply -f k8s/
```

### 환경별 설정

- **Development**: `.env.development`
- **Staging**: `.env.staging`
- **Production**: `.env.production`

## 모니터링

### Prometheus 메트릭
```
GET /metrics
```

주요 메트릭:
- `http_requests_total` - 총 요청 수
- `http_request_duration_seconds` - 요청 처리 시간
- `db_queries_total` - 데이터베이스 쿼리 수
- `cache_hits_total` - 캐시 히트 수
- `orders_total` - 총 주문 수
- `revenue_total` - 총 매출

### 헬스 체크
```
GET /health
GET /health/live
GET /health/ready
```

## 보안

### 구현된 보안 기능
- HTTPS 필수
- JWT 토큰 인증
- Rate limiting
- CORS 설정
- SQL Injection 방지 (GORM)
- XSS 방지
- 비밀번호 암호화 (bcrypt)
- 환경 변수 보호

### 권장 사항
- 환경 변수는 절대 커밋하지 마세요
- 프로덕션에서는 강력한 JWT_SECRET 사용
- HTTPS 강제 설정
- Rate limiting 적절히 조정
- 정기적인 보안 업데이트

## 성능 최적화

### 캐싱 전략
- 상품 목록: 5분 TTL
- 상품 상세: 1시간 TTL
- 사용자 세션: 15분 TTL

### 데이터베이스 인덱스
```sql
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
```

### 커넥션 풀
- 데이터베이스: 최대 25 커넥션
- Redis: 최대 10 커넥션

## 트러블슈팅

### 일반적인 문제

**문제**: 데이터베이스 연결 실패
```bash
# 데이터베이스 상태 확인
docker-compose ps postgres

# 로그 확인
docker-compose logs postgres
```

**문제**: Redis 연결 실패
```bash
# Redis 상태 확인
redis-cli ping
```

**문제**: Stripe 웹훅 검증 실패
```bash
# 웹훅 시크릿 확인
stripe listen --forward-to localhost:8080/api/v1/payments/webhook
```

## 라이선스

MIT

## 지원

- GitHub Issues: https://github.com/modsynth/examples/issues
- Documentation: https://docs.modsynth.io
