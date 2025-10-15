# Quickstart: Points System

## Modes
- Service (HTTP): Run a REST API server
- Library (Go): Import and call use cases directly

## Service Mode
1. Requirements: Go 1.20+, Docker (for MySQL/Redis)
2. Start dependencies:
   - `docker compose up -d mysql redis`
3. Environment (examples):
   - `export MYSQL_DSN="acto:acto@tcp(127.0.0.1:3306)/acto?parseTime=true&charset=utf8mb4&loc=Local"`
   - `export REDIS_ADDR="127.0.0.1:6379"`
   - `export HTTP_ADDR=":8080"`
4. Run server:
   - `make run` (hot reload with Air) or `go run ./app`
5. Swagger UI:
   - Generated via `swag init`; visit `/swagger/index.html`

## Library Mode
1. Import module: `go get github.com/yourorg/acto`
2. Construct services by providing repository interfaces (MySQL/Redis impls in `internal/`)
3. Call use cases:
   - `CreditPoints(ctx, userId, pointTypeName, amount, reason)`
   - `DebitPoints(ctx, userId, pointTypeName, amount, reason)`
   - `ListTransactions(ctx, userId, filter)`

## Error Model
- Service: JSON `{ code, message, details? }`
- Library: typed errors for programmatic checks

## Authorization
- Roles: admin, user
- Admin: manage point types, reward rules, distributions
- User: view balances/history, redeem
