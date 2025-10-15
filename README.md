# Acto Points System

A Clean Architecture Go project providing a multi-type integer points system. Supports both HTTP service and embeddable library modes.

## Tech Stack
- Language: Go 1.20
- HTTP Router: gorilla/mux
- Storage: MySQL (balances, transactions, config)
- Cache/Ranking: Redis (ZSET per point type)
- API Docs: OpenAPI (YAML in `specs/001-/contracts/openapi.yaml`)

## Project Layout (key paths)
- `app/` – Application entrypoint (HTTP server, DI)
- `domain/points/` – Domain entities and errors (pure Go)
- `points/` – Use case layer (business logic, interfaces)
- `internal/repository/mysql/` – MySQL repositories and migrations
- `internal/repository/redis/` – Redis repositories (ranking)
- `internal/rest/` – HTTP handlers and middleware
- `specs/001-/` – Specification, plan, contracts, tasks, docs

## Getting Started (Service Mode)
1. Prerequisites: Go 1.20+, Docker (for MySQL/Redis)
2. Start dependencies:
   - `docker compose up -d mysql redis`
3. Environment variables (examples):
   - `export MYSQL_DSN="acto:acto@tcp(127.0.0.1:3306)/acto?parseTime=true&charset=utf8mb4&loc=Local"`
   - `export REDIS_ADDR="127.0.0.1:6379"`
   - `export HTTP_ADDR=":8080"`
4. Run the server:
   - `go run ./app`

Health check: `GET http://localhost:8080/health`

## Key Endpoints (examples)
- Point Types
  - `POST /api/v1/point-types`
  - `GET /api/v1/point-types`
  - `PATCH /api/v1/point-types`
- Balances (by point type name)
  - `POST /api/v1/users/balance/credit`
  - `POST /api/v1/users/balance/debit`
  - `GET  /api/v1/users/{userId}/transactions`
- Rankings
  - `GET  /api/v1/rankings?pointTypeId=...&limit=...&offset=...`
- Reward Distribution
  - `POST /api/v1/distributions`
- Redemption
  - `POST /api/v1/redeem`

OpenAPI: see `specs/001-/contracts/openapi.yaml`

## Library Mode (Use Cases)
Import the module and wire repository interfaces with your own implementations or the provided ones under `internal/`.

Example (pseudo):
```go
ptSvc := points.NewPointTypeService(ptRepo)
balSvc := points.NewBalanceService(balRepo, rankingRepo, ptRepo)
_ = balSvc.Credit(ctx, userID, "Gold", "Welcome bonus", 100)
```

## Development
- Build: `go build ./...`
- Test: `go test ./...`
- Lint: `golangci-lint run`
- Compose: `docker compose up -d` / `docker compose down`

## Notes
- All point values are integers (no decimals)
- Concurrency safety via row-level locks and transactions
- Error responses use `{ code, message, details? }`

## License
MIT (or project-specific)
