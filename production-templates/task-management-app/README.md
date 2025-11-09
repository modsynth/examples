# Task Management App - Production Template

프로덕션 준비가 완료된 풀스택 작업 관리 애플리케이션입니다. (Trello/Asana 스타일)

## 사용된 Modsynth 모듈

### Backend
- **auth-module** - JWT 인증
- **db-module** - PostgreSQL
- **cache-module** - Redis
- **logging-module** - 구조화된 로깅
- **websocket-client** - 실시간 업데이트
- **notification-module** - 이메일 알림
- **file-storage-module** - 첨부 파일
- **monitoring-module** - 성능 모니터링

### Frontend
- **ui-components** - React UI 컴포넌트
- **auth-client** - 인증 클라이언트
- **api-client** - REST API 클라이언트
- **state-management** - Redux Toolkit
- **routing** - React Router
- **form-validation** - React Hook Form
- **websocket-client** - 실시간 통신
- **error-handling** - Error Boundary
- **i18n** - 다국어 지원

## 기능

### 사용자 관리
- 회원가입/로그인
- 소셜 로그인 (Google, GitHub)
- 프로필 관리
- 팀 초대 및 관리

### 프로젝트 관리
- 프로젝트 생성/수정/삭제
- 프로젝트 멤버 관리
- 권한 관리 (Owner, Admin, Member)
- 프로젝트 템플릿

### 작업 관리
- 칸반 보드 (Todo, In Progress, Done)
- 작업 생성/수정/삭제
- 드래그 앤 드롭
- 작업 할당
- 마감일 설정
- 라벨 및 태그
- 체크리스트
- 첨부 파일
- 댓글

### 실시간 협업
- WebSocket 기반 실시간 동기화
- 다른 사용자의 작업 실시간 반영
- 온라인 사용자 표시
- 실시간 알림

### 알림
- 이메일 알림
- 브라우저 푸시 알림
- 작업 할당 알림
- 마감일 알림
- 댓글 알림

## 아키텍처

```
task-management-app/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── domain/
│   │   ├── repository/
│   │   ├── service/
│   │   └── websocket/
│   │       └── hub.go          # WebSocket hub
│   ├── migrations/
│   ├── docker/
│   ├── go.mod
│   └── README.md
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Board/
│   │   │   │   ├── Board.tsx
│   │   │   │   ├── Column.tsx
│   │   │   │   └── Card.tsx
│   │   │   ├── Task/
│   │   │   │   ├── TaskModal.tsx
│   │   │   │   ├── TaskForm.tsx
│   │   │   │   └── CommentList.tsx
│   │   │   ├── Project/
│   │   │   └── Layout/
│   │   ├── pages/
│   │   │   ├── Dashboard.tsx
│   │   │   ├── Projects.tsx
│   │   │   ├── Board.tsx
│   │   │   └── Settings.tsx
│   │   ├── store/
│   │   │   ├── slices/
│   │   │   │   ├── authSlice.ts
│   │   │   │   ├── projectsSlice.ts
│   │   │   │   └── tasksSlice.ts
│   │   │   └── store.ts
│   │   ├── hooks/
│   │   │   ├── useWebSocket.ts
│   │   │   └── useAuth.ts
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   └── websocket.ts
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   ├── vite.config.ts
│   └── README.md
│
└── docker-compose.yml
```

## 설치 및 실행

### 1. 환경 설정

```bash
# 레포지토리 클론
git clone https://github.com/modsynth/examples.git
cd examples/production-templates/task-management-app

# 환경 변수 설정
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

**Backend 환경 변수** (`backend/.env`):
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskapp
DB_PASSWORD=your-password
DB_NAME=taskapp_db

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your-secret-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d

# OAuth (선택사항)
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GITHUB_CLIENT_ID=
GITHUB_CLIENT_SECRET=

# Email
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# S3
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=taskapp-files
```

**Frontend 환경 변수** (`frontend/.env`):
```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_GOOGLE_CLIENT_ID=
VITE_GITHUB_CLIENT_ID=
```

### 2. Docker로 실행

```bash
# 전체 스택 실행
docker-compose up -d

# 로그 확인
docker-compose logs -f

# 접속
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080
# API Docs: http://localhost:8080/swagger
```

### 3. 로컬 개발

**Backend:**
```bash
cd backend

# 의존성 설치
go mod download

# 데이터베이스 마이그레이션
make migrate-up

# 개발 서버 실행
make dev
```

**Frontend:**
```bash
cd frontend

# 의존성 설치
npm install

# 개발 서버 실행
npm run dev
```

## 주요 기능 구현

### 1. 실시간 동기화 (WebSocket)

**Backend (Go):**
```go
// internal/websocket/hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            delete(h.clients, client)
        case message := <-h.broadcast:
            for client := range h.clients {
                client.send <- message
            }
        }
    }
}
```

**Frontend (React):**
```typescript
// hooks/useWebSocket.ts
export const useWebSocket = (projectId: string) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const dispatch = useDispatch();

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/ws/${projectId}`);

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);

      switch (data.type) {
        case 'TASK_CREATED':
          dispatch(addTask(data.payload));
          break;
        case 'TASK_UPDATED':
          dispatch(updateTask(data.payload));
          break;
        case 'TASK_MOVED':
          dispatch(moveTask(data.payload));
          break;
      }
    };

    setSocket(ws);
    return () => ws.close();
  }, [projectId]);

  return socket;
};
```

### 2. 드래그 앤 드롭 (React DnD)

```typescript
// components/Board/Card.tsx
import { useDrag, useDrop } from 'react-dnd';

export const Card: React.FC<CardProps> = ({ task, columnId }) => {
  const [{ isDragging }, drag] = useDrag({
    type: 'TASK',
    item: { id: task.id, columnId },
    collect: (monitor) => ({
      isDragging: monitor.isDragging(),
    }),
  });

  return (
    <div ref={drag} className={isDragging ? 'opacity-50' : ''}>
      <TaskCard task={task} />
    </div>
  );
};
```

### 3. 권한 관리

```go
// internal/middleware/auth.go
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(*domain.User)

        for _, role := range roles {
            if user.Role == role {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
```

## API 엔드포인트

### 인증
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/google
POST   /api/v1/auth/github
POST   /api/v1/auth/refresh
```

### 프로젝트
```
GET    /api/v1/projects
POST   /api/v1/projects
GET    /api/v1/projects/:id
PUT    /api/v1/projects/:id
DELETE /api/v1/projects/:id
POST   /api/v1/projects/:id/members
DELETE /api/v1/projects/:id/members/:userId
```

### 작업
```
GET    /api/v1/projects/:projectId/tasks
POST   /api/v1/projects/:projectId/tasks
GET    /api/v1/tasks/:id
PUT    /api/v1/tasks/:id
DELETE /api/v1/tasks/:id
PUT    /api/v1/tasks/:id/move
POST   /api/v1/tasks/:id/comments
POST   /api/v1/tasks/:id/attachments
```

### WebSocket
```
WS     /api/v1/ws/:projectId
```

## 테스트

```bash
# Backend 테스트
cd backend
make test

# Frontend 테스트
cd frontend
npm test

# E2E 테스트
npm run test:e2e
```

## 배포

### Kubernetes 배포

```bash
# 이미지 빌드
docker build -t taskapp-backend:v1.0.0 ./backend
docker build -t taskapp-frontend:v1.0.0 ./frontend

# 배포
kubectl apply -f k8s/
```

### Vercel + Railway 배포

**Frontend (Vercel):**
```bash
npm install -g vercel
cd frontend
vercel --prod
```

**Backend (Railway):**
```bash
railway login
railway up
```

## 성능 최적화

### Frontend
- React.memo로 불필요한 리렌더 방지
- 가상 스크롤 (react-window)
- 이미지 최적화
- 코드 스플리팅
- Service Worker 캐싱

### Backend
- Redis 캐싱 (프로젝트, 작업 목록)
- 데이터베이스 인덱스
- 커넥션 풀 최적화
- WebSocket 메시지 배칭

## 모니터링

### Metrics
- 활성 사용자 수
- WebSocket 연결 수
- API 응답 시간
- 데이터베이스 쿼리 성능
- 작업 생성/완료 통계

### Logging
- 구조화된 로그 (JSON)
- 로그 레벨: DEBUG, INFO, WARN, ERROR
- ELK Stack 통합 가능

## 보안

- JWT 토큰 인증
- HTTPS 필수
- CORS 설정
- Rate limiting
- XSS 방지
- CSRF 토큰
- SQL Injection 방지

## 라이선스

MIT

## 스크린샷

(프로덕션 배포 후 추가)

## 기여

버그 리포트 및 기능 요청은 GitHub Issues로 제출해주세요.
