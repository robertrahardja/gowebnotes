## External Resources

## Development and Deployment Pipeline

```mermaid
flowchart LR
    %% Developer actions
    DEV1[devs] --> PUSH[push code]
    DEV2[devs] --> PUSH
    DEV3[devs] --> PUSH

    %% Source code and trigger
    PUSH --> SOURCE[source code]
    SOURCE --> TRIGGER[trigger]

    %% CI Pipeline
    TRIGGER --> TESTS["Tests<br/>• Static check<br/>• Linter<br/>• Build"]

    %% Branch check and build trigger
    TESTS --> BRANCH{"If main<br/>branch<br/>merged"}
    BRANCH -->|"Trigger a new build"| BUILD[Build]

    %% Build outcomes
    BUILD -->|passed| ARTIFACT[Artifact Repository]
    BUILD -.->|failed| SOURCE

    %% API versions and deployment
    ARTIFACT --> API1[API v1.2]
    ARTIFACT --> API2[API v1.1]
    ARTIFACT --> API3[...]

    %% Cloud deployment
    API1 --> CLOUD[Cloud Run]
    API2 --> CLOUD
    API3 --> CLOUD

    %% Final endpoint
    CLOUD --> HTTPS[HTTPS]

```

## 1. Initial Setup and Project Structure

**Create project structure:**

```bash
mkdir gopher-social
cd gopher-social
go mod init social
```

**Create folder structure:**

```
├── bin/
├── cmd/
│   ├── api/
│   └── migrate/
│       └── migrations/
├── internal/
├── docs/
└── scripts/
```

## 2. Basic HTTP Server Setup

**Create main files:**

`cmd/api/main.go`:

```go
package main

import (
    "log"
    "net/http"
    "time"
)

type application struct {
    config config
}

type config struct {
    addr string
}

func main() {
    cfg := config{
        addr: ":8080",
    }

    app := &application{
        config: cfg,
    }

    log.Fatal(app.run())
}

func (app *application) run() error {
    srv := &http.Server{
        Addr:         app.config.addr,
        Handler:      app.mount(),
        WriteTimeout: time.Second * 30,
        ReadTimeout:  time.Second * 10,
        IdleTimeout:  time.Minute,
    }

    log.Printf("server has started %s", app.config.addr)
    return srv.ListenAndServe()
}
```

`cmd/api/api.go`:

```go
package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func (app *application) mount() http.Handler {
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    r.Route("/v1", func(r chi.Router) {
        r.Get("/health", app.healthCheckHandler)
    })

    return r
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("OK"))
}
```

**Install Chi router:**

```bash
go get github.com/go-chi/chi/v5
```

## 3. Database Setup with Docker

**Create docker-compose.yml:**

```yaml
version: "3.8"
services:
  postgres:
    image: postgres:13
    container_name: gopher_social_db
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: adminpassword
      POSTGRES_DB: socialnetwork
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Start database:**

```bash
docker-compose up -d
```

**Install database dependencies:**

```bash
go get github.com/lib/pq
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
```

## 4. Environment Configuration

**Install godotenv:**

```bash
go get github.com/joho/godotenv
```

**Create .env file:**

```env
DB_ADDR=postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable
ENV=development
```

**Create config with environment variables:**

`cmd/api/env.go`:

```go
package main

import (
    "os"
    "strconv"
)

func envString(key, fallback string) string {
    val, ok := os.LookupEnv(key)
    if !ok {
        return fallback
    }
    return val
}

func envInt(key string, fallback int) int {
    val, ok := os.LookupEnv(key)
    if !ok {
        return fallback
    }
    valAsInt, err := strconv.Atoi(val)
    if err != nil {
        return fallback
    }
    return valAsInt
}
```

## 5. Database Migrations

**Create Makefile:**

```makefile
include .env

migration:
 migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
 migrate -path=cmd/migrate/migrations -database=${DB_ADDR} up

migrate-down:
 migrate -path=cmd/migrate/migrations -database=${DB_ADDR} down

.PHONY: migration migrate-up migrate-down
```

**Create first migration:**

```bash
make migration create_users_table
```

**User migration file (`cmd/migrate/migrations/001_create_users_table.up.sql`):**

```sql
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username varchar(100) NOT NULL UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    is_active boolean NOT NULL DEFAULT false
);
```

**Run migration:**

```bash
make migrate-up
```

## 6. Repository Pattern Implementation

**Create storage interface (`internal/store/storage.go`):**

```go
package store

import (
    "context"
    "database/sql"
)

type Storage struct {
    Users interface {
        Create(context.Context, *User) error
        GetByID(context.Context, int64) (*User, error)
        GetByEmail(context.Context, string) (*User, error)
    }
}

func NewStorage(db *sql.DB) Storage {
    return Storage{
        Users: &UserStore{db: db},
    }
}
```

**Create user model (`internal/store/users.go`):**

```go
package store

import (
    "context"
    "crypto/sha256"
    "database/sql"
    "encoding/hex"
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"
)

var (
    ErrNotFound = errors.New("record not found")
    ErrDuplicateEmail = errors.New("user already exists with this email")
    ErrDuplicateUsername = errors.New("user already exists with this username")
)

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  password  `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    IsActive  bool      `json:"is_active"`
}

type password struct {
    text *string
    hash []byte
}

func (p *password) Set(plaintextPassword string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    p.text = &plaintextPassword
    p.hash = hash
    return nil
}

type UserStore struct {
    db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (username, email, password)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`

    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.hash).Scan(
        &user.ID,
        &user.CreatedAt,
    )

    if err != nil {
        // Handle unique constraint violations
        return err
    }

    return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
    query := `
        SELECT id, username, email, created_at, is_active
        FROM users
        WHERE id = $1 AND is_active = true`

    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    user := &User{}
    err := s.db.QueryRowContext(ctx, query, userID).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
        &user.IsActive,
    )

    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrNotFound
        default:
            return nil, err
        }
    }

    return user, nil
}
```

## 7. JSON Helpers

**Create JSON utilities (`cmd/api/json.go`):**

```go
package main

import (
    "encoding/json"
    "net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
    maxBytes := 1_048_576 // 1MB
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()

    return dec.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
    type envelope struct {
        Error string `json:"error"`
    }

    return writeJSON(w, status, &envelope{Error: message})
}
```

## 8. User Registration

**Install validator:**

```bash
go get github.com/go-playground/validator/v10
```

**Create user registration handler:**

```go
type CreateUserPayload struct {
    Username string `json:"username" validate:"required,max=100"`
    Email    string `json:"email" validate:"required,email,max=255"`
    Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
    var payload CreateUserPayload
    if err := readJSON(w, r, &payload); err != nil {
        writeJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    if err := Validate.Struct(payload); err != nil {
        writeJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    user := &store.User{
        Username: payload.Username,
        Email:    payload.Email,
    }

    if err := user.Password.Set(payload.Password); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    ctx := r.Context()
    if err := app.store.Users.Create(ctx, user); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    if err := writeJSON(w, http.StatusCreated, user); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }
}
```

## 9. JWT Authentication

**Install JWT library:**

```bash
go get github.com/golang-jwt/jwt/v5
```

**Create auth package (`internal/auth/auth.go`):**

```go
package auth

import (
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
    GenerateToken(claims jwt.Claims) (string, error)
    ValidateToken(token string) (*jwt.Token, error)
}

type JWTAuthenticator struct {
    secret   []byte
    aud      string
    iss      string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
    return &JWTAuthenticator{
        secret: []byte(secret),
        aud:    aud,
        iss:    iss,
    }
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(a.secret)
    if err != nil {
        return "", err
    }
    return tokenString, nil
}
```

**Create token endpoint:**

```go
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
    type TokenPayload struct {
        Email    string `json:"email" validate:"required,email"`
        Password string `json:"password" validate:"required"`
    }

    var payload TokenPayload
    if err := readJSON(w, r, &payload); err != nil {
        writeJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    // Validate and authenticate user
    user, err := app.store.Users.GetByEmail(r.Context(), payload.Email)
    if err != nil {
        writeJSONError(w, http.StatusUnauthorized, "Invalid credentials")
        return
    }

    // Generate JWT token
    claims := jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
        "iat": time.Now().Unix(),
        "nbf": time.Now().Unix(),
        "iss": "gopher-social",
        "aud": "gopher-social",
    }

    token, err := app.authenticator.GenerateToken(claims)
    if err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    if err := writeJSON(w, http.StatusCreated, map[string]string{"token": token}); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }
}
```

## 10. Middleware for Authentication

**Create authentication middleware:**

```go
func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            writeJSONError(w, http.StatusUnauthorized, "Authorization header missing")
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            writeJSONError(w, http.StatusUnauthorized, "Authorization header malformed")
            return
        }

        token, err := app.authenticator.ValidateToken(parts[1])
        if err != nil {
            writeJSONError(w, http.StatusUnauthorized, "Invalid token")
            return
        }

        claims := token.Claims.(jwt.MapClaims)
        userID := int64(claims["sub"].(float64))

        user, err := app.store.Users.GetByID(r.Context(), userID)
        if err != nil {
            writeJSONError(w, http.StatusUnauthorized, "User not found")
            return
        }

        ctx := context.WithValue(r.Context(), userKey, user)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## 11. Posts CRUD Implementation

**Create posts migration:**

```bash
make migration create_posts_table
```

**Posts migration (`cmd/migrate/migrations/002_create_posts_table.up.sql`):**

```sql
CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    content text NOT NULL,
    title varchar(100) NOT NULL,
    user_id bigint NOT NULL,
    tags varchar(100) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
```

**Create posts handlers:**

```go
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
    type CreatePostPayload struct {
        Title   string `json:"title" validate:"required,max=100"`
        Content string `json:"content" validate:"required"`
        Tags    string `json:"tags"`
    }

    var payload CreatePostPayload
    if err := readJSON(w, r, &payload); err != nil {
        writeJSONError(w, http.StatusBadRequest, err.Error())
        return
    }

    user := getUserFromContext(r)

    post := &store.Post{
        Title:   payload.Title,
        Content: payload.Content,
        Tags:    payload.Tags,
        UserID:  user.ID,
    }

    ctx := r.Context()
    if err := app.store.Posts.Create(ctx, post); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }

    if err := writeJSON(w, http.StatusCreated, post); err != nil {
        writeJSONError(w, http.StatusInternalServerError, "Internal server error")
        return
    }
}
```

## 12. Swagger Documentation

**Install Swagger:**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/http-swagger/v2
```

**Add to Makefile:**

```makefile
gen-docs:
 swag init -g ./cmd/api --parseDependency --parseInternal

.PHONY: gen-docs
```

**Generate docs:**

```bash
make gen-docs
```

## 13. Testing Setup

**Create test utilities (`cmd/api/test_utils.go`):**

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func newTestApp(t *testing.T) *application {
    t.Helper()

    return &application{
        store: store.NewMockStore(),
        cache: cache.NewMockStore(),
    }
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    mux.ServeHTTP(rr, req)
    return rr
}
```

**Create mock stores (`internal/store/mocks.go`):**

```go
package store

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, user *User) error {
    return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
    return &User{ID: userID}, nil
}

func NewMockStore() Storage {
    return Storage{
        Users: &MockUserStore{},
    }
}
```

## 14. Rate Limiting

**Create rate limiter (`internal/ratelimit/fixed_window.go`):**

```go
package ratelimit

import (
    "net"
    "sync"
    "time"
)

type Limiter interface {
    Allow(ip string) (bool, time.Duration)
}

type FixedWindowLimiter struct {
    mutex    sync.RWMutex
    clients  map[string]int
    limit    int
    window   time.Duration
}

func NewFixedWindow(limit int, window time.Duration) *FixedWindowLimiter {
    return &FixedWindowLimiter{
        clients: make(map[string]int),
        limit:   limit,
        window:  window,
    }
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
    rl.mutex.Lock()
    defer rl.mutex.Unlock()

    count, exists := rl.clients[ip]
    if !exists {
        rl.clients[ip] = 1
        go rl.resetClient(ip)
        return true, 0
    }

    if count >= rl.limit {
        return false, rl.window
    }

    rl.clients[ip]++
    return true, 0
}

func (rl *FixedWindowLimiter) resetClient(ip string) {
    time.Sleep(rl.window)
    rl.mutex.Lock()
    delete(rl.clients, ip)
    rl.mutex.Unlock()
}
```

## 15. CI/CD Setup

**Create GitHub Actions (`.github/workflows/audit.yml`):**

```yaml
name: Audit

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.22

      - name: Verify dependencies
        run: go mod verify

      - name: Build
        run: go build -v ./...

      - name: Run go vet
        run: go vet ./...

      - name: Install staticcheck
        run: go install honnef.co/go/tools/cmd/staticcheck@latest

      - name: Run staticcheck
        run: staticcheck ./...

      - name: Run tests
        run: go test -race -vet=off ./...
```

## 16. Deployment to Google Cloud

**Create Dockerfile:**

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/api/*.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

**Deploy to Cloud Run:**

```bash
# Build and push image
gcloud builds submit --tag gcr.io/PROJECT-ID/gopher-social

# Deploy to Cloud Run
gcloud run deploy gopher-social \
  --image gcr.io/PROJECT-ID/gopher-social \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

## 17. Running the Complete Application

**Final commands to run everything:**

```bash
# Start database
docker-compose up -d

# Load environment variables
source .env

# Run migrations
make migrate-up

# Generate API documentation
make gen-docs

# Run tests
make test

# Start the server
go run cmd/api/*.go
```

---

## 1. Environment Variables Setup

**Create environment variable handler:**

```bash
# Create internal/env package
mkdir -p internal/env
```

```go
// internal/env/env.go
package env

import "os"

func GetString(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}

func GetInt(key string, fallback int) int {
    if value, ok := os.LookupEnv(key); ok {
        if val, err := strconv.Atoi(value); err == nil {
            return val
        }
    }
    return fallback
}
```

**Setup environment loading:**

```bash
# Create .envrc file (for direnv)
echo 'export SERVER_ADDR=":3000"' > .envrc
direnv allow
```

## 2. Database Connection Pool Configuration

**Configure database in main application:**

```go
// cmd/api/main.go
type config struct {
    addr string
    db   dbConfig
}

type dbConfig struct {
    addr         string
    maxOpenConns int
    maxIdleConns int
    maxIdleTime  string
}

func main() {
    cfg := config{
        addr: env.GetString("SERVER_ADDR", ":8080"),
        db: dbConfig{
            addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/socialnetwork?sslmode=disable"),
            maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
            maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
            maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
        },
    }
}
```

**Create database connection pool:**

```go
// internal/db/db.go
package db

import (
    "context"
    "database/sql"
    "time"
    _ "github.com/lib/pq"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
    db, err := sql.Open("postgres", addr)
    if err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)

    duration, err := time.ParseDuration(maxIdleTime)
    if err != nil {
        return nil, err
    }
    db.SetConnMaxIdleTime(duration)

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, err
    }

    return db, nil
}
```

## 3. Docker Compose for Database

**Create docker-compose.yaml:**

```yaml
services:
  postgres:
    image: postgres:16.3
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: adminpassword
      POSTGRES_DB: socialnetwork
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

**Start database:**

```bash
docker-compose up -d
```

## 4. Hot Reloading Setup

**Install and configure Air:**

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Initialize air configuration
air init
```

**Configure .air.toml:**

```toml
[build]
  cmd = "go build -o ./bin/api ./cmd/api"
  bin = "bin/api"
  full_bin = "./bin/api"
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["bin", "web", "docs", "scripts"]
```

**Start with hot reload:**

```bash
air
```

## 5. SQL Migrations

**Install migrate tool:**

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Create Makefile:**

```makefile
include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migrate-create:
 migrate create -seq -ext=.sql -dir=$(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
 migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
 migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down
```

**Create migrations:**

```bash
# Create users table
make migrate-create create_users

# Create posts table
make migrate-create create_posts

# Create comments table
make migrate-create add_comments
```

**Users migration (up):**

```sql
-- cmd/migrate/migrations/001_create_users.up.sql
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    username varchar(255) UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
```

**Posts migration (up):**

```sql
-- cmd/migrate/migrations/002_create_posts.up.sql
CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    title varchar(100) NOT NULL,
    content text NOT NULL,
    tags varchar(100)[],
    user_id bigint NOT NULL,
    version integer NOT NULL DEFAULT 0,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

ALTER TABLE posts ADD CONSTRAINT fk_posts_user_id
FOREIGN KEY (user_id) REFERENCES users (id);
```

**Run migrations:**

```bash
make migrate-up
```

## 6. JSON Encoding/Decoding and Error Handling

**Create internal errors package:**

```go
// internal/errors/errors.go
package errors

import (
    "log"
    "net/http"
)

type application struct {
    // dependencies
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
    log.Printf("internal server error: %s path: %s error: %v", r.Method, r.URL.Path, err)

    response := map[string]string{"error": "the server encountered a problem"}
    app.writeJSON(w, http.StatusInternalServerError, response)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
    log.Printf("bad request error: %s path: %s error: %v", r.Method, r.URL.Path, err)

    response := map[string]string{"error": err.Error()}
    app.writeJSON(w, http.StatusBadRequest, response)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
    log.Printf("not found error: %s path: %s error: %v", r.Method, r.URL.Path, err)

    response := map[string]string{"error": "resource not found"}
    app.writeJSON(w, http.StatusNotFound, response)
}
```

**JSON utilities:**

```go
// internal/json/json.go
package json

import (
    "encoding/json"
    "net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(data)
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
    return json.NewDecoder(r.Body).Decode(data)
}
```

## 7. HTTP Payload Validation

**Install validation package:**

```bash
go get github.com/go-playground/validator/v10
```

**Setup validation:**

```go
// internal/json/json.go (add to existing file)
import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
    validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(data interface{}) error {
    return validate.Struct(data)
}
```

**Example handler with validation:**

```go
// API handler example
type CreatePostPayload struct {
    Title   string   `json:"title" validate:"required,max=100"`
    Content string   `json:"content" validate:"required,max=1000"`
    Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
    var payload CreatePostPayload

    if err := json.ReadJSON(w, r, &payload); err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    if err := json.ValidateStruct(payload); err != nil {
        app.badRequestResponse(w, r, err)
        return
    }

    // Create post logic...

    app.writeJSONResponse(w, http.StatusCreated, payload)
}
```

## 8. Standardized JSON Responses

**Create response wrapper:**

```go
// Update internal/json/json.go
type Envelope map[string]interface{}

func (app *application) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
    envelope := Envelope{"data": data}
    WriteJSON(w, status, envelope)
}
```

## 9. Database Query Timeouts

**Add context timeouts to store methods:**

```go
// internal/store/posts.go
const QueryTimeoutDuration = time.Second * 5

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
    ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
    defer cancel()

    query := `SELECT id, title, content, user_id, created_at FROM posts WHERE id = $1`

    var post Post
    err := s.db.QueryRowContext(ctx, query, id).Scan(
        &post.ID,
        &post.Title,
        &post.Content,
        &post.UserID,
        &post.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &post, nil
}
```

## 10. Database Seeding

**Create seed script:**

```go
// cmd/migrate/seed/main.go
package main

import (
    "log"
    "context"
    "internal/db"
    "internal/store"
)

func main() {
    addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable")

    conn, err := db.New(addr, 3, 3, "15m")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    store := store.NewStorage(conn)

    if err := seed(store); err != nil {
        log.Fatal(err)
    }

    log.Println("Seeding completed")
}

func seed(store *store.Storage) error {
    ctx := context.Background()

    // Generate users
    users := generateUsers(100)
    for _, user := range users {
        if _, err := store.Users.Create(ctx, user); err != nil {
            return err
        }
    }

    // Generate posts and comments...
    return nil
}
```

**Add to Makefile:**

```makefile
.PHONY: seed
seed:
 go run cmd/migrate/seed/main.go
```

**Run seed:**

```bash
make seed
```

## Commands Summary

```bash
# Setup
mkdir -p internal/{env,db,json,errors}
go mod init social-api
go get github.com/go-playground/validator/v10
go get github.com/lib/pq

# Hot reload
go install github.com/cosmtrek/air@latest
air init
air

# Database
docker-compose up -d
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
make migrate-create create_users
make migrate-up

# Seeding
make seed

# Development
direnv allow  # Load environment variables
air          # Start with hot reload
```
