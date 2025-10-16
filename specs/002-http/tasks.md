# Tasks: Enhanced Library Integration

**Branch**: `002-http` | **Spec**: specs/002-http/spec.md | **Plan**: specs/002-http/plan.md

## Dependencies (Story Order)

1. US1 - Simplified Library Initialization (P1)
2. US3 - Route Integration (P1)
3. US2 - Custom Repository Support (P2)
4. Polish & Cross-Cutting

## Phase 1: Setup

- [X] T001 Create `lib/` package directory at repository root
- [X] T002 Create `lib/config.go` with `LibraryConfig`, `RepositoryConfig`, `RouteConfig`, `CORSConfig`
- [X] T003 Create `lib/errors.go` with typed `LibraryError` and `ErrorType`
- [X] T004 Create `lib/factory.go` with `NewLibrary` and `NewLibraryWithContext`
- [X] T005 Create `lib/router.go` with router registration stubs (mux, gin, echo)

## Phase 2: Foundational

- [X] T006 Implement config validation in `lib/config.go`
- [X] T007 Implement error constructors in `lib/errors.go`
- [X] T008 Wire default repositories using provided DB/Redis in `lib/factory.go`
- [X] T009 Build `Services` container and health checks in `lib/factory.go`
- [X] T010 Add JSON tags (camelCase) to all new structs in `lib/*.go`

## Phase 3: US1 - Simplified Library Initialization (P1)

Goal: Initialize services with only DB and Redis. Independent test: create and list point types.

- [X] T011 [US1] Implement `NewLibrary(cfg LibraryConfig) (*Library, error)` in `lib/factory.go`
- [X] T012 [P] [US1] Inject MySQL repositories using `cfg.DB` in `lib/factory.go`
- [X] T013 [P] [US1] Inject Redis ranking repository using `cfg.Redis` in `lib/factory.go`
- [X] T014 [US1] Construct `Services` (point types, balance, distribution, redemption) in `lib/factory.go`
- [X] T015 [US1] Implement `(*Library).Health() error` in `lib/factory.go`
- [X] T016 [US1] Implement `(*Library).Close() error` resource cleanup in `lib/factory.go`

## Phase 4: US3 - Route Integration (P1)

Goal: Mount project routes inside host app router. Independent test: mount under `/api/v1` and call health/config/services.

- [X] T018 [US3] Implement Gorilla Mux registration in `lib/router.go`
- [ ] T019 [P] [US3] (cancelled) Implement Gin registration in `lib/router.go`
- [ ] T020 [P] [US3] (cancelled) Implement Echo registration in `lib/router.go`
- [X] T021 [US3] Support `RouteConfig.PathPrefix` and middleware chain in `lib/router.go`
- [X] T022 [US3] Implement `/health`, `/config`, `/services` handlers using `Library` in `internal/rest/handlers/library.go`
- [X] T023 [US3] Wire handlers into router registrations in `lib/router.go`
- [X] T024 [US3] Verify OpenAPI contracts in `specs/002-http/contracts/openapi.yaml`

## Phase 5: US2 - Custom Repository Support (P2)

Goal: Allow overrides for any repository. Independent test: swap PointTypeRepository and verify calls.

- [X] T025 [US2] Extend `RepositoryConfig` with optional overrides in `lib/config.go`
- [X] T026 [US2] In `NewLibrary`, prefer overrides over defaults in `lib/factory.go`
- [X] T027 [US2] Add validation for overrides implementing correct interfaces in `lib/config.go`
- [X] T028 [US2] Example usage in `examples/custom_repo_example.go`

## Phase 6: Polish & Cross-Cutting

- [X] T029 Add godoc comments to all exported types/functions in `lib/*.go`
- [X] T030 Ensure error wrapping and context on failures in `lib/*.go`
- [X] T031 Add linters for new package paths in `.golangci.yaml` if needed
- [X] T032 Ensure naming conventions (camelCase JSON, snake_case DB) across new code
- [X] T033 Update `README.md` with library usage section and examples links
- [ ] T035 (cancelled) Add integration health test using sqlmock and miniredis under `specs/002-http/tests/`
 
### Phase 6 Progress
- [X] T034 Add basic unit tests for `lib/factory.go` and `lib/router.go`

## Parallel Execution Examples

- T012 and T013 can run in parallel (separate repository wiring)
- T029–T033 can run in parallel after core features complete

## MVP Scope

- Complete Phases 1–3 (US1) to deliver a usable library with simple initialization.
