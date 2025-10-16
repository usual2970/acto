# Feature Specification: DI cleanup with uber-go/dig

**Feature Branch**: `003-previously-after-implementing`  
**Created**: 2025-10-16  
**Status**: Draft  
**Input**: User description: "Previously, after implementing the core features, we optimized the user experience. However, the project’s dependency injection during route initialization is somewhat disorganized. In this phase, we plan to use uber-go/dig to manage dependency injection more effectively."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Consistent DI container bootstrap (Priority: P1)

As a developer, I can initialize the application using a single DI container builder so that all dependencies (config, DB, Redis, repositories, services, handlers) are resolved in one place without scattered manual wiring.

**Why this priority**: Eliminates fragile manual wiring and supports maintainability and onboarding.

**Independent Test**: Boot the app by only editing container configuration; routes resolve successfully without touching route code.

**Acceptance Scenarios**:

1. Given an empty container, When I invoke the container builder, Then all core dependencies are provided.
2. Given the app starts, When routes are registered, Then handlers are resolved from the container via Invoke with no manual constructors.

---

### User Story 2 - Clear module boundaries for DI (Priority: P1)

As a developer, I can register modules (config, infra, repositories, services, delivery) with explicit Provide/Invoke boundaries so each module can be tested in isolation and reused.

**Why this priority**: Encourages separation of concerns and enables module-level tests.

**Independent Test**: Module unit test creates a minimal container with only that module and validates it resolves.

**Acceptance Scenarios**:

1. Given DI modules exist, When I include only the repository module, Then repository tests run without the HTTP module.
2. Given DI modules exist, When I include services + fake repos, Then service tests pass without DB.

---

### User Story 3 - Minimal surface for handlers (Priority: P2)

As a developer, I can register handlers through an adapter that only depends on interfaces so that route code is thin and framework-agnostic.

**Why this priority**: Keeps delivery layer free of infrastructure details.

**Independent Test**: Route registration compiles while swapping the HTTP framework adapter.

**Acceptance Scenarios**:

1. Given a handler interface, When container builds, Then handler is injected only with its required service interfaces.
2. Given a new framework adapter, When I register routes, Then no changes are needed to handlers/services.

---

### Edge Cases

- What happens when a dependency fails to construct?
  - Container build fails fast with a descriptive error surfaced at startup.
- How does the system handle optional dependencies (e.g., Redis not configured)?
  - Container supports optional providers with sensible fallbacks and clear error messages when features requiring them are invoked.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Provide a single container builder that registers config, DB, cache, repositories, services, and handlers.
- **FR-002**: Route registration MUST use Invoke to resolve handlers; no manual constructor calls in route files.
- **FR-003**: Define DI modules per layer (config/infra/repos/services/delivery) with clear provide boundaries.
- **FR-004**: Handlers MUST depend on service interfaces only (no repositories), resolved via container.
- **FR-005**: Container build MUST fail fast with a readable error if any provider fails.
- **FR-006**: Support test containers that replace infra providers with fakes/mocks.
- **FR-007**: Document container usage with examples for adding a new dependency.
- **FR-008**: Provide a migration guide from existing manual wiring to container wiring.
- **FR-009**: Expose a small adapter to register routes framework-agnostically via interface.

### Key Entities *(include if feature involves data)*

- **DI Module**: A group of providers for a layer (e.g., repositories).
- **Container Adapter**: Thin wrapper exposing Provide/Invoke and registration helpers.
- **Route Registrar Interface**: Minimal interface to bind method/path to handler.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Route files contain zero manual constructor calls (100% DI-resolved).
- **SC-002**: New feature wiring time reduced by 50% (baseline: current manual wiring time).
- **SC-003**: Container build fails with actionable error messages in 100% of misconfiguration test cases.
- **SC-004**: Module-level tests run without starting HTTP server or real DB.
