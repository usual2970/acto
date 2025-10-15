# Research: Points System

## Unknowns and Decisions

### HTTP Router
- Decision: Use `gorilla/mux`
- Rationale: Stable, feature-complete, familiar API; meets routing needs without heavy dependencies
- Alternatives considered: `chi` (lighter, also fine), `net/http` only (manual routing overhead)

### OpenAPI Tooling (NEEDS CLARIFICATION)
- Decision: Use `swaggo/swag` for annotation-based generation
- Rationale: Simple integration with Go handlers, generates Swagger UI quickly
- Alternatives considered: `ogen` (code-first strong typing, steeper adoption), manual YAML (error-prone)

### Redis Usage for Rankings
- Decision: Use Redis ZSET per point type for real-time rankings
- Rationale: O(log N) updates and range queries, sorted sets fit ranking semantics
- Alternatives considered: SQL window functions (heavy for real-time), custom in-memory heap (not shared across instances)

### Transaction Isolation Level / Engine
- Decision: MySQL InnoDB Read Committed (or REPEATABLE READ with `SELECT ... FOR UPDATE`), using explicit row-level locks
- Rationale: Matches spec intent (row-level locking) and MySQL defaults; ensures balance consistency under concurrency
- Alternatives considered: Serializable (contention, throughput impact), Read Uncommitted (dirty reads)

### Testing Stack
- Decision: `testify` + `mockery`, integration tests with Dockerized Postgres/Redis
- Rationale: Familiar, productive, supports constitution’s coverage goals
- Alternatives considered: Go mocks (ok), custom fakes (time-costly)

## Consolidated Choices

- Router: `gorilla/mux`
- OpenAPI: `swaggo/swag`
- DB: MySQL via `go-sql-driver/mysql`
- Cache/Ranking: Redis via `go-redis`
- Tests: `go test`, `testify`, `mockery`

## Impacts

- Clean Architecture preserved; router/DB/cache isolated in `internal/`
- Enables fast leaderboard updates; consistent error model across service/library


