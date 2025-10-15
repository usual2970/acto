# Tasks: Points System (Feature `001-`)

This file lists implementation tasks in execution order, organized by phases and user stories. Follow Clean Architecture structure documented in `plan.md`.

## Phase 1 — Setup

- [X] T001 Create project structure at repo root (`app/`, `domain/points/`, `points/`, `internal/repository/mysql/`, `internal/repository/redis/`, `internal/rest/handlers/`, `internal/rest/middleware/`, `tests/`)
- [X] T002 Add dependencies in `go.mod`: `github.com/gorilla/mux`, `github.com/go-sql-driver/mysql`, `github.com/redis/go-redis/v9`, `github.com/stretchr/testify`, `github.com/vektra/mockery/v2`, `github.com/swaggo/swag`
- [X] T003 Create `compose.yaml` with MySQL + Redis services and default credentials at repo root
- [X] T004 Add `Makefile` targets: run, test, lint, swag, mocks, compose-up/down at repo root
- [X] T005 Add `.golangci.yaml` and initial rules at repo root
- [X] T006 Add `.air.toml` and `app/main.go` bootstrap using `mux.Router`

## Phase 2 — Foundational (cross-cutting)

- [X] T007 Create domain errors in `domain/points/errors.go` (insufficient balance, cannot delete active type, etc.)
- [X] T008 Define core domain entities in `domain/points/entities.go` (PointType, UserBalance, Transaction, RewardRule, RewardDistribution, RedemptionReward, RedemptionRecord)
- [X] T009 Define use case interfaces in `points/service.go` (repositories declared at consumer side)
- [X] T010 [P] Add config loader in `internal/config/config.go` (DB DSN, Redis, server)
- [X] T011 [P] Add logging setup in `internal/log/log.go`
- [X] T012 [P] Add error response model and mapper in `internal/rest/handlers/errors.go`
- [X] T013 [P] Add auth role middleware in `internal/rest/middleware/auth.go`
- [X] T014 Wire DI in `app/main.go` (DB pool, Redis client, repositories, services, mux routes)

## Phase 3 — User Story 1 (P1): Points Configuration Management

Goal: Operations staff can create, list, update, delete point types (no delete if balances exist). Independent test: see `spec.md` US1.

- [X] T015 [US1] Create MySQL migrations for point types in `internal/repository/mysql/migrations/001_point_types.sql`
- [X] T016 [US1] Implement PointType repository in `internal/repository/mysql/point_type_repo.go`
- [X] T017 [US1] Implement PointType use cases in `points/point_type_service.go`
- [X] T018 [P] [US1] Implement REST handlers in `internal/rest/handlers/point_types.go`
- [X] T019 [P] [US1] Register mux routes in `app/main.go` (`/api/v1/point-types` CRUD)

## Phase 4 — User Story 2 (P1): Points Balance Management

Goal: Credit/debit balances with history, prevent negative, per-type separation. Independent test: see `spec.md` US2.

- [X] T020 [US2] Create MySQL migrations for balances and transactions in `internal/repository/mysql/migrations/002_balances_transactions.sql`
- [X] T021 [US2] Implement Balance/Transaction repository in `internal/repository/mysql/balance_tx_repo.go`
- [X] T022 [US2] Implement Balance service methods (Credit, Debit) with row locks in `points/balance_service.go`
- [X] T023 [P] [US2] Implement REST handlers in `internal/rest/handlers/balances.go`
- [X] T024 [P] [US2] Register mux routes in `app/main.go` (`/users/{userId}/balances/...`)

## Phase 5 — User Story 3 (P2): Points Ranking System

Goal: Rankings per point type, descending by balance, ties same rank, real-time updates. Independent test: see `spec.md` US3.

- [X] T025 [US3] Implement Redis ranking repo in `internal/repository/redis/ranking_repo.go` (ZSET per point type)
- [X] T026 [US3] Update balance service to publish ranking updates in `points/balance_service.go`
- [X] T027 [P] [US3] Implement ranking handlers in `internal/rest/handlers/rankings.go`
- [X] T028 [P] [US3] Register mux routes in `app/main.go` (`/rankings/{pointTypeId}`)

## Phase 6 — User Story 4 (P2): Rank-Based Reward Distribution

Goal: Configure rank-based reward rules and execute distributions per snapshot; prevent duplicates. Independent test: see `spec.md` US4.

- [X] T029 [US4] Create MySQL migrations for reward rules and distributions in `internal/repository/mysql/migrations/003_rewards.sql`
- [X] T030 [US4] Implement RewardRule/Distribution repository in `internal/repository/mysql/rewards_repo.go`
- [X] T031 [US4] Implement Distribution service in `points/distribution_service.go`
- [X] T032 [P] [US4] Implement REST handlers in `internal/rest/handlers/distributions.go`
- [X] T033 [P] [US4] Register mux routes in `app/main.go` (`/distributions`)

## Phase 7 — User Story 5 (P1): Points Redemption System

Goal: Catalog browse, validate multi-type costs, atomic deduction, record redemption. Independent test: see `spec.md` US5.

- [X] T034 [US5] Create MySQL migrations for redemption rewards/records in `internal/repository/mysql/migrations/004_redemptions.sql`
- [X] T035 [US5] Implement Redemption repository in `internal/repository/mysql/redemption_repo.go`
- [X] T036 [US5] Implement Redemption service in `points/redemption_service.go`
- [X] T037 [P] [US5] Implement REST handlers in `internal/rest/handlers/redemptions.go`
- [X] T038 [P] [US5] Register mux routes in `app/main.go` (`/rewards`, `/redeem`)

## Final Phase — Polish & Cross-Cutting

- [X] T039 Update pagination helpers and apply to list endpoints in `internal/rest/handlers/*`
- [X] T040 Add filtering for transaction history in `internal/rest/handlers/balances.go`
- [X] T041 Validate error code mapping and update `internal/rest/handlers/errors.go`
- [X] T042 Update `specs/001-/contracts/openapi.yaml` to reflect final handlers' request/response schemas
- [X] T043 Update `specs/001-/quickstart.md` with run commands and environment examples

---

## Dependencies (Story Order)

1. Setup → 2. Foundational → 3. US1 (P1) → 4. US2 (P1) → 5. US3 (P2) → 6. US4 (P2) → 7. US5 (P1) → Final Polish

## Parallel Execution Examples

- In Foundational: T010/T011/T012/T013 can run in parallel once T009 exists
- In US1: T018 and T019 can run in parallel after T017 is defined
- In US2: T023 and T024 can run in parallel after T022 is defined
- In US3/US4/US5: handler + route tasks can run in parallel after service methods exist

## Implementation Strategy

- MVP scope: Deliver US1 fully (CRUD for point types) enabling configuration
- Incremental delivery: US2 to enable core earn/spend; US3 for rankings; US4 for automation; US5 for redemption
- Keep domain pure; use interfaces in `points/` and implementations in `internal/`

## Independent Test Criteria (from spec)

- US1: Create/list/update/delete point types as specified (no delete with active balances)
- US2: Credit/debit with history; prevent negative; per-type separation
- US3: Ranking order correctness, real-time updates, tie handling
- US4: Rank-based rewards distribution with duplicate prevention and scheduling
- US5: Redemption with multi-type costs, atomic deduction, inventory handling

## Format Validation

All tasks follow the required checklist format: `- [ ] T### [P] [US#] Description with file path`.
