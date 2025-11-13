# Task Management App - Backend API

A production-ready Kanban-style task management system with real-time collaboration via WebSocket. Built with Go, Gin, PostgreSQL, and GORM.

## Features

### Core Functionality
- **User Authentication**: JWT-based auth with access and refresh tokens
- **Project Management**: Create and manage projects with team collaboration
- **Kanban Boards**: Multiple boards per project with customizable columns
- **Task Management**: Comprehensive task CRUD with drag-and-drop support
- **Real-time Collaboration**: WebSocket-based live updates for task changes
- **Role-Based Access Control**: Owner, Admin, Member, and Viewer roles
- **Task Features**:
  - Comments and discussions
  - Checklist items for subtasks
  - Labels and categorization
  - Due dates and priorities
  - File attachments
  - Task assignment

### Technical Features
- Clean architecture with layered design (Handler → Service → Repository → Domain)
- RESTful API with 30+ endpoints
- WebSocket Hub pattern for real-time updates
- Database migrations with GORM AutoMigrate
- Docker support with multi-stage builds
- Graceful shutdown handling
- CORS middleware
- Request logging
- Error handling with context wrapping

## Architecture

```
cmd/
  server/
    main.go           # Application entry point
internal/
  domain/            # Domain models and DTOs
    user.go
    project.go
    board.go
    task.go
  repository/        # Data access layer
    user_repo.go
    project_repo.go
    board_repo.go
    task_repo.go
  service/           # Business logic layer
    auth_service.go
    project_service.go
    board_service.go
    task_service.go
  handler/           # HTTP handlers
    auth.go
    project.go
    board.go
    task.go
  middleware/        # HTTP middleware
    auth.go
    cors.go
    logger.go
  websocket/         # WebSocket infrastructure
    hub.go
    client.go
  config/            # Configuration management
    config.go
migrations/          # SQL migrations
docker/              # Docker configuration
```

## Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **WebSocket**: gorilla/websocket
- **Environment**: godotenv
- **Password Hashing**: bcrypt

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14+
- Docker and Docker Compose (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd task-management-app/backend
```

2. Install dependencies:
```bash
make deps
```

3. Copy environment file:
```bash
cp .env.example .env
```

4. Update `.env` with your configuration:
```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=taskapp
DB_PASSWORD=your_password
DB_NAME=taskapp_db
DB_SSLMODE=disable

JWT_SECRET=your-secret-key-change-this-in-production
JWT_EXPIRATION=15  # minutes
```

### Running with Docker

```bash
# Start all services (PostgreSQL + Redis + App)
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
```

### Running Locally

1. Start PostgreSQL:
```bash
# Using Docker
docker run -d \
  --name taskapp-postgres \
  -e POSTGRES_USER=taskapp \
  -e POSTGRES_PASSWORD=taskapp_password \
  -e POSTGRES_DB=taskapp_db \
  -p 5432:5432 \
  postgres:14-alpine
```

2. Run the application:
```bash
make dev
```

The server will start on `http://localhost:8080`

### Building

```bash
# Build binary
make build

# Run binary
./bin/server
```

## API Endpoints

### Authentication
```
POST   /api/v1/auth/register      # Register new user
POST   /api/v1/auth/login         # Login user
POST   /api/v1/auth/refresh       # Refresh access token
GET    /api/v1/profile            # Get user profile (protected)
```

### Projects
```
GET    /api/v1/projects           # List user's projects
POST   /api/v1/projects           # Create project
GET    /api/v1/projects/:id       # Get project details
PUT    /api/v1/projects/:id       # Update project
DELETE /api/v1/projects/:id       # Delete project
POST   /api/v1/projects/:id/archive      # Archive project
POST   /api/v1/projects/:id/unarchive    # Unarchive project

# Project Members
GET    /api/v1/projects/:id/members               # List members
POST   /api/v1/projects/:id/members               # Add member
DELETE /api/v1/projects/:id/members/:memberID     # Remove member
PUT    /api/v1/projects/:id/members/:memberID/role  # Update member role
```

### Boards
```
POST   /api/v1/projects/:projectID/boards   # Create board
GET    /api/v1/projects/:projectID/boards   # List project boards
GET    /api/v1/boards/:id                   # Get board details
PUT    /api/v1/boards/:id                   # Update board
DELETE /api/v1/boards/:id                   # Delete board
```

### Tasks
```
POST   /api/v1/boards/:boardID/tasks        # Create task
GET    /api/v1/boards/:boardID/tasks        # List board tasks
GET    /api/v1/tasks/:id                    # Get task details
PUT    /api/v1/tasks/:id                    # Update task
DELETE /api/v1/tasks/:id                    # Delete task
POST   /api/v1/tasks/:id/move               # Move task to another board

# Task Comments
POST   /api/v1/tasks/:id/comments           # Add comment
DELETE /api/v1/tasks/:id/comments/:commentID  # Delete comment

# Task Checklist
POST   /api/v1/tasks/:id/checklist          # Add checklist item
PUT    /api/v1/tasks/:id/checklist/:itemID  # Update checklist item
DELETE /api/v1/tasks/:id/checklist/:itemID  # Delete checklist item

# Task Labels
POST   /api/v1/tasks/:id/labels             # Assign labels to task
```

### WebSocket
```
GET    /api/v1/ws/:projectId                # WebSocket connection (requires auth)
GET    /api/v1/projects/:projectId/online-users  # Get online users
```

### Health Check
```
GET    /health                              # Health check endpoint
```

## WebSocket Events

### Client → Server
```json
{
  "type": "PING",
  "project_id": 1,
  "user_id": 1
}
```

### Server → Client
```json
{
  "type": "TASK_CREATED",
  "project_id": 1,
  "user_id": 2,
  "data": {
    "id": 123,
    "title": "New Task",
    ...
  }
}
```

### Event Types
- `TASK_CREATED` - New task created
- `TASK_UPDATED` - Task updated
- `TASK_DELETED` - Task deleted
- `TASK_MOVED` - Task moved to another board
- `BOARD_CREATED` - New board created
- `BOARD_UPDATED` - Board updated
- `BOARD_DELETED` - Board deleted
- `COMMENT_ADDED` - Comment added to task
- `COMMENT_DELETED` - Comment deleted
- `CHECKLIST_ITEM_ADDED` - Checklist item added
- `CHECKLIST_ITEM_UPDATED` - Checklist item updated
- `CHECKLIST_ITEM_DELETED` - Checklist item deleted
- `TASK_LABELS_UPDATED` - Task labels changed

## Authentication

All protected endpoints require a JWT token in the Authorization header:

```bash
Authorization: Bearer <access_token>
```

### Token Flow
1. Register or login to receive `access_token` and `refresh_token`
2. Use `access_token` (expires in 15 minutes) for API requests
3. When `access_token` expires, use `refresh_token` to get a new one
4. `refresh_token` expires in 7 days

### Example
```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword",
    "username": "johndoe",
    "full_name": "John Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword"
  }'

# Use token
curl http://localhost:8080/api/v1/projects \
  -H "Authorization: Bearer <access_token>"
```

## Role-Based Access Control

### Project Roles (Hierarchical)
- **Owner**: Full control, can delete project, manage all members
- **Admin**: Manage members, create/delete boards, archive project
- **Member**: Create/edit/delete tasks, boards (can't delete boards), add comments
- **Viewer**: Read-only access

### Permission Matrix
| Action | Owner | Admin | Member | Viewer |
|--------|-------|-------|--------|--------|
| View project | ✓ | ✓ | ✓ | ✓ |
| Update project | ✓ | ✓ | ✗ | ✗ |
| Delete project | ✓ | ✗ | ✗ | ✗ |
| Archive project | ✓ | ✓ | ✗ | ✗ |
| Add members | ✓ | ✓ | ✗ | ✗ |
| Remove members | ✓ | ✓ | ✗ | ✗ |
| Change roles | ✓ | ✗ | ✗ | ✗ |
| Create boards | ✓ | ✓ | ✓ | ✗ |
| Delete boards | ✓ | ✓ | ✗ | ✗ |
| Create tasks | ✓ | ✓ | ✓ | ✗ |
| Edit tasks | ✓ | ✓ | ✓ | ✗ |
| Delete tasks | ✓ | ✓ | ✓ | ✗ |
| Add comments | ✓ | ✓ | ✓ | ✗ |

## Database Schema

### Users
- id, email (unique), password_hash, username (unique), full_name, avatar_url
- role (user/admin), is_active
- created_at, updated_at

### Projects
- id, name, description, icon, color
- owner_id (FK → users)
- is_archived, created_at, updated_at

### Project Members
- id, project_id (FK → projects), user_id (FK → users)
- role (owner/admin/member/viewer)
- created_at, updated_at

### Boards
- id, project_id (FK → projects), name, position
- created_at, updated_at

### Tasks
- id, board_id (FK → boards), title, description
- position, priority (low/medium/high/urgent)
- due_date, creator_id (FK → users), assignee_id (FK → users)
- is_completed, completed_at
- created_at, updated_at

### Labels
- id, project_id (FK → projects), name, color
- created_at, updated_at

### Comments
- id, task_id (FK → tasks), user_id (FK → users)
- content, created_at, updated_at

### Attachments
- id, task_id (FK → tasks), user_id (FK → users)
- filename, file_url, file_size, mime_type
- created_at

### Checklist Items
- id, task_id (FK → tasks), title
- is_completed, position
- created_at, updated_at

## Development

### Available Make Commands
```bash
make help          # Show available commands
make dev           # Run development server
make build         # Build the application
make test          # Run tests
make lint          # Run linter
make fmt           # Format code
make deps          # Download dependencies
make clean         # Clean build artifacts
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make docker-logs   # View application logs
```

### Project Structure Best Practices
- **Domain Layer**: Pure business entities, no external dependencies
- **Repository Layer**: Data access, GORM operations, error wrapping
- **Service Layer**: Business logic, validation, orchestration, WebSocket events
- **Handler Layer**: HTTP request/response handling, input validation
- **Middleware**: Cross-cutting concerns (auth, logging, CORS)

### Adding New Features
1. Define domain models in `internal/domain/`
2. Create repository interface and implementation in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Create HTTP handlers in `internal/handler/`
5. Register routes in `cmd/server/main.go`
6. Add WebSocket events if needed in service layer

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
go test -race ./...
```

## Deployment

### Docker Production Build

```bash
# Build image
docker build -f docker/Dockerfile -t task-management-app:latest .

# Run container
docker run -d \
  --name taskapp \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-db-password \
  -e JWT_SECRET=your-secret-key \
  task-management-app:latest
```

### Environment Variables

Required for production:
- `ENV=production`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET` (use strong random string)
- `PORT` (default: 8080)

Optional:
- `JWT_EXPIRATION` (default: 15 minutes)
- `REDIS_HOST`, `REDIS_PORT` (for future caching)

## Security Considerations

- **Password Hashing**: bcrypt with default cost (10)
- **JWT Tokens**:
  - Access tokens expire in 15 minutes
  - Refresh tokens expire in 7 days
  - Tokens include user ID, email, username
- **HTTPS**: Always use HTTPS in production
- **CORS**: Configure allowed origins in production
- **SQL Injection**: Protected by GORM's parameterized queries
- **Rate Limiting**: Implement rate limiting for production (TODO)

## Performance

- **Database Connection Pool**:
  - Max open connections: 25
  - Max idle connections: 5
  - Connection max lifetime: 5 minutes
- **WebSocket**: Hub pattern for efficient message broadcasting
- **Indexes**: Unique indexes on email, username, project-user pairs

## Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U taskapp -d taskapp_db
```

### WebSocket Connection Issues
- Ensure JWT token is valid and included in connection
- Check CORS settings for WebSocket upgrade
- Verify project ID exists and user has access

### Build Issues
```bash
# Clean and rebuild
make clean
make deps
make build
```

## License

MIT

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For issues and questions, please open an issue on GitHub.
