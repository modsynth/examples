# Real-Time Chat - Production Template

프로덕션 준비가 완료된 실시간 채팅 애플리케이션입니다. (Slack/Discord 스타일)

## 사용된 Modsynth 모듈

### Backend
- **auth-module** - JWT 인증
- **db-module** - PostgreSQL
- **cache-module** - Redis (메시지 캐싱, 온라인 상태)
- **logging-module** - 구조화된 로깅
- **websocket-client** - WebSocket 통신
- **messaging-module** - RabbitMQ (메시지 큐)
- **file-storage-module** - S3 (파일 첨부, 프로필 이미지)
- **notification-module** - 푸시 알림
- **search-module** - Elasticsearch (메시지 검색)
- **monitoring-module** - Prometheus

### Frontend
- **ui-components** - React UI 컴포넌트
- **auth-client** - 인증
- **api-client** - REST API
- **state-management** - Redux Toolkit
- **routing** - React Router
- **form-validation** - 메시지 입력 검증
- **websocket-client** - 실시간 통신
- **error-handling** - Error Boundary

## 기능

### 채널 및 DM
- 공개/비공개 채널
- 다이렉트 메시지 (1:1, 그룹)
- 채널 생성/수정/삭제
- 채널 초대/퇴장
- 스레드 답글

### 메시지
- 텍스트 메시지
- 파일 첨부 (이미지, 동영상, 문서)
- 이모지 반응
- 메시지 편집/삭제
- 메시지 고정
- 메시지 검색
- 멘션 (@username)
- 코드 블록 하이라이팅

### 실시간 기능
- 실시간 메시지 수신
- 타이핑 인디케이터
- 온라인/오프라인 상태
- 읽음 표시
- 실시간 알림

### 사용자 기능
- 프로필 관리
- 상태 메시지
- 커스텀 상태
- 알림 설정
- 테마 설정 (라이트/다크)

### 고급 기능
- 음성 통화 (WebRTC)
- 화면 공유
- 봇 통합
- 웹훅
- API 키 관리

## 아키텍처

```
real-time-chat/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go
│   │   │   │   ├── channels.go
│   │   │   │   ├── messages.go
│   │   │   │   └── users.go
│   │   │   ├── middleware/
│   │   │   └── routes.go
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   ├── channel.go
│   │   │   ├── message.go
│   │   │   └── workspace.go
│   │   ├── repository/
│   │   ├── service/
│   │   │   ├── chat_service.go
│   │   │   ├── presence_service.go
│   │   │   └── search_service.go
│   │   ├── websocket/
│   │   │   ├── hub.go
│   │   │   ├── client.go
│   │   │   └── message.go
│   │   └── queue/
│   │       └── consumer.go
│   ├── migrations/
│   ├── docker/
│   └── README.md
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── Chat/
│   │   │   │   ├── MessageList.tsx
│   │   │   │   ├── MessageInput.tsx
│   │   │   │   ├── Message.tsx
│   │   │   │   └── Thread.tsx
│   │   │   ├── Sidebar/
│   │   │   │   ├── ChannelList.tsx
│   │   │   │   ├── DirectMessages.tsx
│   │   │   │   └── UserList.tsx
│   │   │   ├── Channel/
│   │   │   │   ├── ChannelHeader.tsx
│   │   │   │   └── ChannelSettings.tsx
│   │   │   └── Voice/
│   │   │       └── VoiceCall.tsx
│   │   ├── pages/
│   │   │   ├── Chat.tsx
│   │   │   ├── DirectMessage.tsx
│   │   │   └── Settings.tsx
│   │   ├── store/
│   │   │   └── slices/
│   │   │       ├── channelsSlice.ts
│   │   │       ├── messagesSlice.ts
│   │   │       └── usersSlice.ts
│   │   ├── hooks/
│   │   │   ├── useWebSocket.ts
│   │   │   ├── usePresence.ts
│   │   │   └── useVoice.ts
│   │   └── services/
│   │       ├── websocket.ts
│   │       └── webrtc.ts
│   └── README.md
│
└── docker-compose.yml
```

## 설치 및 실행

### 1. 환경 설정

**Backend** (`backend/.env`):
```env
PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=chatapp
DB_PASSWORD=your-password
DB_NAME=chatapp_db

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Elasticsearch
ES_ADDRESSES=http://localhost:9200

# S3
S3_ENDPOINT=
S3_ACCESS_KEY=
S3_SECRET_KEY=
S3_BUCKET=chatapp-files

# JWT
JWT_SECRET=your-secret-key
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=7d

# WebRTC (선택사항)
TURN_SERVER_URL=turn:turnserver.com:3478
TURN_USERNAME=
TURN_PASSWORD=
```

**Frontend** (`frontend/.env`):
```env
VITE_API_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_TURN_SERVER_URL=turn:turnserver.com:3478
```

### 2. Docker로 실행

```bash
# 전체 스택 실행
docker-compose up -d

# 접속
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
```

### 3. 로컬 개발

**Backend:**
```bash
cd backend
go mod download
make migrate-up
make dev
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## 핵심 기능 구현

### 1. WebSocket 메시지 처리

**Backend:**
```go
// internal/websocket/hub.go
type Hub struct {
    clients      map[string]*Client
    channels     map[string]map[string]bool
    broadcast    chan *Message
    register     chan *Client
    unregister   chan *Client
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client.ID] = client
            h.sendPresenceUpdate(client.UserID, "online")

        case client := <-h.unregister:
            delete(h.clients, client.ID)
            h.sendPresenceUpdate(client.UserID, "offline")

        case message := <-h.broadcast:
            h.broadcastToChannel(message)
        }
    }
}

func (h *Hub) broadcastToChannel(msg *Message) {
    channel := h.channels[msg.ChannelID]
    for clientID := range channel {
        if client, ok := h.clients[clientID]; ok {
            client.send <- msg
        }
    }
}
```

**Frontend:**
```typescript
// hooks/useWebSocket.ts
export const useWebSocket = () => {
  const dispatch = useDispatch();
  const [socket, setSocket] = useState<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/ws`);

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);

      switch (message.type) {
        case 'MESSAGE':
          dispatch(addMessage(message.payload));
          break;
        case 'TYPING':
          dispatch(setTyping(message.payload));
          break;
        case 'PRESENCE':
          dispatch(updatePresence(message.payload));
          break;
        case 'READ_RECEIPT':
          dispatch(markAsRead(message.payload));
          break;
      }
    };

    setSocket(ws);
    return () => ws.close();
  }, []);

  return { socket, sendMessage, sendTyping };
};
```

### 2. 타이핑 인디케이터

```typescript
// components/Chat/MessageInput.tsx
export const MessageInput: React.FC = () => {
  const { sendTyping } = useWebSocket();
  const [isTyping, setIsTyping] = useState(false);
  const typingTimeout = useRef<NodeJS.Timeout>();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);

    if (!isTyping) {
      setIsTyping(true);
      sendTyping(channelId, true);
    }

    clearTimeout(typingTimeout.current);
    typingTimeout.current = setTimeout(() => {
      setIsTyping(false);
      sendTyping(channelId, false);
    }, 2000);
  };

  return <Input value={value} onChange={handleChange} />;
};
```

### 3. 메시지 검색 (Elasticsearch)

```go
// internal/service/search_service.go
func (s *SearchService) SearchMessages(query string, channelID string) ([]Message, error) {
    results, err := s.searchModule.Search(context.Background(), search.SearchQuery{
        Index: "messages",
        Query: search.MatchQuery{
            Field: "content",
            Value: query,
        },
        Filters: []search.Filter{
            {Field: "channel_id", Value: channelID},
        },
        Size: 50,
    })

    return results, err
}
```

### 4. 음성 통화 (WebRTC)

```typescript
// hooks/useVoice.ts
export const useVoice = (channelId: string) => {
  const [localStream, setLocalStream] = useState<MediaStream | null>(null);
  const [remoteStream, setRemoteStream] = useState<MediaStream | null>(null);
  const peerConnection = useRef<RTCPeerConnection | null>(null);

  const startCall = async () => {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true, video: true });
    setLocalStream(stream);

    const pc = new RTCPeerConnection({
      iceServers: [
        { urls: 'stun:stun.l.google.com:19302' },
        { urls: TURN_SERVER_URL, username: TURN_USERNAME, credential: TURN_PASSWORD }
      ]
    });

    stream.getTracks().forEach(track => pc.addTrack(track, stream));

    pc.ontrack = (event) => {
      setRemoteStream(event.streams[0]);
    };

    pc.onicecandidate = (event) => {
      if (event.candidate) {
        socket.send(JSON.stringify({
          type: 'ICE_CANDIDATE',
          candidate: event.candidate
        }));
      }
    };

    peerConnection.current = pc;
  };

  return { localStream, remoteStream, startCall, endCall };
};
```

### 5. 읽음 표시

```go
// internal/service/chat_service.go
func (s *ChatService) MarkAsRead(userID string, channelID string, messageID string) error {
    // 읽음 상태 저장
    err := s.repo.SaveReadReceipt(userID, channelID, messageID)
    if err != nil {
        return err
    }

    // Redis에 캐싱
    key := fmt.Sprintf("read:%s:%s", userID, channelID)
    s.cache.Set(context.Background(), key, messageID, 24*time.Hour)

    // WebSocket으로 다른 사용자에게 알림
    s.hub.BroadcastToChannel(channelID, &Message{
        Type: "READ_RECEIPT",
        Payload: map[string]interface{}{
            "user_id": userID,
            "message_id": messageID,
        },
    })

    return nil
}
```

## API 엔드포인트

### 인증
```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
```

### 워크스페이스
```
GET    /api/v1/workspaces
POST   /api/v1/workspaces
GET    /api/v1/workspaces/:id
PUT    /api/v1/workspaces/:id
```

### 채널
```
GET    /api/v1/workspaces/:workspaceId/channels
POST   /api/v1/workspaces/:workspaceId/channels
GET    /api/v1/channels/:id
PUT    /api/v1/channels/:id
DELETE /api/v1/channels/:id
POST   /api/v1/channels/:id/members
DELETE /api/v1/channels/:id/members/:userId
```

### 메시지
```
GET    /api/v1/channels/:channelId/messages
POST   /api/v1/channels/:channelId/messages
PUT    /api/v1/messages/:id
DELETE /api/v1/messages/:id
POST   /api/v1/messages/:id/reactions
POST   /api/v1/messages/:id/thread
GET    /api/v1/messages/search
```

### 다이렉트 메시지
```
GET    /api/v1/dm
POST   /api/v1/dm
GET    /api/v1/dm/:id/messages
```

### 사용자
```
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
GET    /api/v1/users/:id/presence
PUT    /api/v1/users/:id/status
```

### WebSocket
```
WS     /api/v1/ws
```

## 메시지 타입

### 클라이언트 → 서버
```typescript
{
  type: 'MESSAGE' | 'TYPING' | 'READ_RECEIPT' | 'REACTION' | 'JOIN_CHANNEL' | 'LEAVE_CHANNEL',
  payload: {
    channel_id: string,
    content?: string,
    user_id: string,
    // ...
  }
}
```

### 서버 → 클라이언트
```typescript
{
  type: 'MESSAGE' | 'TYPING' | 'PRESENCE' | 'READ_RECEIPT' | 'REACTION' | 'NOTIFICATION',
  payload: {
    id: string,
    channel_id: string,
    user: User,
    content: string,
    created_at: string,
    // ...
  }
}
```

## 성능 최적화

### 메시지 로딩
- 무한 스크롤 (가상 스크롤)
- 메시지 페이지네이션 (50개씩)
- Redis 캐싱 (최근 메시지)

### WebSocket 최적화
- 메시지 배칭
- 압축 (gzip)
- Heartbeat (30초)

### 데이터베이스
- 인덱스: channel_id, user_id, created_at
- 파티셔닝: 월별 메시지 테이블
- 아카이빙: 6개월 이상 오래된 메시지

## 모니터링

### Metrics
- 활성 WebSocket 연결 수
- 초당 메시지 수
- 채널별 활성 사용자
- API 응답 시간
- 메시지 전송 지연 시간

### Alerting
- WebSocket 연결 실패율 > 5%
- 메시지 전송 지연 > 500ms
- 데이터베이스 커넥션 풀 > 80%

## 보안

- End-to-End 암호화 (선택사항)
- HTTPS/WSS 필수
- JWT 토큰 인증
- Rate limiting
- CORS 설정
- XSS 방지
- 파일 업로드 검증

## 확장성

### Horizontal Scaling
- Redis Pub/Sub로 여러 서버 간 메시지 동기화
- RabbitMQ로 메시지 큐잉
- 로드 밸런서 (Nginx, HAProxy)

### Database Sharding
- 워크스페이스별 샤딩
- 읽기 복제본

## 라이선스

MIT

## 데모

(프로덕션 배포 후 추가)

## 기여

이슈 및 PR은 환영합니다!
