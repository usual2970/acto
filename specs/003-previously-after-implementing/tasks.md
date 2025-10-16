# Tasks: DI cleanup with uber-go/dig (Feature `003-previously-after-implementing`)

This file lists implementation tasks in execution order, organized by phases and user stories. Follow Clean Architecture and the spec/plan.

## Phase 1 — Setup

- [X] T001 Add uber-go/dig dependency in `go.mod`
- [X] T002 Create DI container package directory `lib/container/`
- [X] T003 Create route registrar interfaces/adapters file `lib/router.go`
- [X] T004 Create library-level config contract file `lib/config.go`

## Phase 2 — Foundational (DI modules and builder)

- [X] T005 Create container builder `lib/container/builder.go` (new dig.Container, Provide modules, Bootstrap)
- [X] T006 Create modules definition `lib/container/modules.go` (config, infra, repositories, services, delivery)
- [X] T007 [P] Add testing helpers `lib/container/testing.go` (replace providers for tests)
- [X] T008 Define registrar interface and mux adapter in `lib/router.go` (Handle method, path, handler)
- [X] T009 Replace manual wiring in `app/main.go` to build container and use registrar (no manual constructors)

## Phase 3 — [US1] Consistent DI container bootstrap (P1)

Goal: Single DI builder initializes config, DB, Redis, repos, services, handlers; route init uses container.

- [X] T010 [US1] Provide config provider in `lib/container/modules.go`
- [X] T011 [US1] Provide DB and Redis providers in `lib/container/modules.go`
- [X] T012 [US1] Provide repository providers in `lib/container/modules.go`
- [X] T013 [US1] Provide service providers in `lib/container/modules.go`
- [X] T014 [US1] Provide handler providers in `lib/container/modules.go`
- [X] T015 [P] [US1] Invoke route registration using resolved handlers in `app/main.go`

## Phase 4 — [US2] Clear module boundaries for DI (P1)

Goal: Modules per layer with clear Provide/Invoke boundaries; isolated tests.

- [X] T016 [US2] Split modules by layer functions in `lib/container/modules.go` (ConfigModule, InfraModule, RepoModule, ServiceModule, DeliveryModule)
- [X] T017 [US2] Add minimal module-level test using dig in `lib/container/testing.go`
- [X] T018 [P] [US2] Document module usage in `specs/003-previously-after-implementing/quickstart.md`

## Phase 5 — [US3] Minimal surface for handlers (P2)

Goal: Handlers depend only on service interfaces; delivery is framework-agnostic.

- [X] T019 [US3] Audit handlers to ensure no repository is injected (services only) in `internal/rest/handlers/*`
- [X] T020 [US3] Ensure registrar adapter used for route registration in `app/main.go`
- [X] T021 [P] [US3] Add example of swapping registrar (doc snippet) in `quickstart.md`

## Final Phase — Polish & Cross-Cutting

- [X] T022 Add failure diagnostics on container build errors in `lib/container/builder.go`
- [X] T023 Update `README.md` with DI overview and getting started
- [X] T024 Update `specs/003-previously-after-implementing/contracts/openapi.yaml` note (no public endpoints)
- [X] T025 Validate that route files contain zero manual constructors (grep check) and adjust `app/main.go`

---

## Dependencies (Story Order)

1. Setup → 2. Foundational → 3. US1 (P1) → 4. US2 (P1) → 5. US3 (P2) → Final Polish

## Parallel Execution Examples

- In Foundational: T007 and T008 can run in parallel after T006
- In US1: T010–T014 can be implemented in sequence, T015 in parallel after T014
- In US2: T017 and T018 can run in parallel after T016
- In US3: T020 and T021 can run in parallel after T019

## Implementation Strategy

- MVP scope: US1 (container builder + route resolution via container)
- Incremental: US2 (modules and tests), US3 (registrar and handler surface), then Polish

## Independent Test Criteria (from spec)

- US1: App boots with DI builder; routes resolve via Invoke; no manual constructors
- US2: Module tests pass with isolated containers and fake providers
- US3: Handlers only depend on service interfaces; registrar swap requires no handler/service code changes

## Format Validation

All tasks follow the required checklist format: `- [ ] T### [P] [US#] Description with file path`.
