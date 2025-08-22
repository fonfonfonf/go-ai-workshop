## Cursor Backend API

Minimal Go Fiber API using Clean Architecture.

### Architecture
- Domain: `internal/models`
- Use cases: `internal/usecases`
- Interfaces (HTTP): `internal/handlers`, `internal/routes`
- Infrastructure (adapters): `internal/storage/memory`
- Wiring: `cmd/main.go`

Flow: HTTP -> handlers -> usecases -> repositories (interface) -> storage (implementation)

### Endpoints (base `/api/v1`)
- GET `/healthz`
- POST `/api/v1/login` { "email": "...", "password": "..." }
- GET `/api/v1/me` (Bearer token)
- PUT `/api/v1/me` { "name": "..." } (Bearer token)

### Sequence diagrams (Mermaid)
```mermaid
sequenceDiagram
  autonumber
  participant C as Client
  participant H as Fiber Handler
  participant M as JWT Middleware
  participant A as AuthService
  participant U as UserRepository (SQLite)

  rect rgb(245,245,245)
    note over C,H: Login
    C->>H: POST /api/v1/login {email,password}
    H->>A: Login(email,password)
    A->>U: GetByEmail(email)
    U-->>A: User(email, passwordHash)
    A-->>A: bcrypt compare
    A-->>H: JWT token + profile
    H-->>C: 200 {token,email,name}
  end

  rect rgb(245,245,245)
    note over C,H: Get profile
    C->>H: GET /api/v1/me (Authorization: Bearer)
    H->>M: validate JWT
    M-->>H: userID
    H->>A: GetProfile(userID)
    A->>U: GetByID(userID)
    U-->>A: User
    A-->>H: User
    H-->>C: 200 {id,email,name}
  end

  rect rgb(245,245,245)
    note over C,H: Update profile
    C->>H: PUT /api/v1/me {name} (Bearer)
    H->>M: validate JWT
    M-->>H: userID
    H->>A: UpdateProfile(userID,name)
    A->>U: GetByID(userID)
    U-->>A: User
    A-->>U: Update(User{name})
    U-->>A: Updated User
    A-->>H: Updated User
    H-->>C: 200 {id,email,name}
  end
```

### Run locally
```bash
export PORT=8080
export SQLITE_DSN="file:aiworkshop.db?cache=shared&mode=rwc"
export JWT_SECRET="change-me"
go mod tidy
go run ./cmd
```

### Example
```bash
curl -s -X POST http://localhost:8001/api/v1/todos/ \
  -H 'Content-Type: application/json' \
  -d '{"title":"first"}'
```

### Layout
```
cmd/
  main.go
internal/
  handlers/
    auth_handler.go
  middleware/
    jwt_middleware.go
  models/
    user.go
  repositories/
    user_repository.go
  routes/
    routes.go
  storage/
    sqlite/
      user_sqlite.go
  usecases/
    auth_service.go
```

