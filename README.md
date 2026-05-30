# FreakDream

Social network rewritten in Go + React. This repository is a monorepo with three top-level directories:

```
backend/    Go REST API + WebSocket server
frontend/   React SPA (Vite + TypeScript + Tailwind CSS)
```

The `class/`, `template/`, and `index.php` directories contain the original PHP codebase kept for reference.

---

## Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.24, chi router, sqlx, MySQL 8 |
| Auth | JWT (HS256), bcrypt |
| Real-time | gorilla/websocket |
| Frontend | React 19, Vite, TypeScript, Tailwind CSS |
| State | Zustand |
| HTTP client | Axios |
| Infra | Docker Compose, Nginx |

---

## Quick start (Docker)

```bash
cp backend/.env.example backend/.env   # edit secrets
docker compose up --build
```

App runs at **http://localhost**. API at **http://localhost/api**.

---

## Development

### Backend

```bash
cd backend
cp .env.example .env         # configure DB and secrets
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
npm run dev                  # proxies /api → localhost:8080
```

### Database

Apply the migration once against a running MySQL instance:

```bash
mysql -u root -p < backend/migrations/001_initial_schema.sql
```

---

## API endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | ✗ | Login, returns JWT |
| POST | `/api/auth/register` | ✗ | Register new account |
| POST | `/api/auth/logout` | ✓ | Logout (client clears token) |
| GET | `/api/users/me` | ✓ | Authenticated user profile |
| GET | `/api/users/:id` | ✓ | Public user profile |
| GET | `/api/users/search?q=` | ✓ | Search users by name/about |
| GET | `/api/users/:id/posts` | ✓ | List posts by user |
| GET | `/api/feed` | ✓ | Feed of followed users |
| POST | `/api/posts` | ✓ | Create text post |
| POST | `/api/posts/photo` | ✓ | Upload photo post (multipart) |
| POST | `/api/posts/share` | ✓ | Share/repost |
| DELETE | `/api/posts/:id` | ✓ | Delete own post |
| POST | `/api/follow/:id` | ✓ | Follow user |
| DELETE | `/api/follow/:id` | ✓ | Unfollow user |
| GET | `/api/follow/:id/count` | ✓ | Follower count |
| GET | `/api/follow/:id/status` | ✓ | Is following? |
| POST | `/api/friends/request/:id` | ✓ | Send friend request |
| POST | `/api/friends/accept/:id` | ✓ | Accept friend request |
| DELETE | `/api/friends/:id` | ✓ | Remove friend |
| GET | `/api/friends/requests` | ✓ | Pending friend requests |
| GET | `/api/notifications` | ✓ | List notifications |
| DELETE | `/api/notifications/:id` | ✓ | Remove notification |
| PATCH | `/api/notifications/:id/read` | ✓ | Mark as read |
| GET | `/api/chat/:id/history` | ✓ | Chat message history |
| GET | `/ws/chat` | ✓ (JWT query param) | WebSocket chat |

---

## Security improvements over original PHP code

- Passwords hashed with **bcrypt** (replaced MD5 without salt)
- Custom XOR encryption removed; data stored as plain text
- JWT-based stateless authentication (no sessions)
- Rate limiting on auth endpoints (5 req/min per IP)
- Parameterised SQL queries (no string interpolation)
- CORS configured to specific origins
