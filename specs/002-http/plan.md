# Implementation Plan: Enhanced Library Integration

**Branch**: `002-http` | **Date**: 2024-12-19 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-http/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Create a simplified library initialization system that allows developers to integrate Acto points system with minimal setup (database + Redis connections), while maintaining support for custom repository implementations and HTTP route integration. The approach involves creating factory functions and configuration objects that abstract away the complexity of manual dependency injection.

## Technical Context

**Language/Version**: Go 1.21+ (current project standard)  
**Primary Dependencies**: 
- `database/sql` for database abstraction
- `github.com/go-sql-driver/mysql` for MySQL driver
- `github.com/redis/go-redis/v9` for Redis client
- `github.com/gorilla/mux` for HTTP routing
- `github.com/gin-gonic/gin` for Gin router support
- `github.com/labstack/echo/v4` for Echo router support

**Storage**: MySQL database + Redis for caching/ranking  
**Testing**: Go's built-in testing framework with `testify` for assertions  
**Target Platform**: Cross-platform Go applications (Linux, macOS, Windows)  
**Project Type**: Library package for Go applications  
**Performance Goals**: 
- Library initialization < 100ms
- Support 1000+ concurrent operations
- Memory usage < 50MB for basic usage

**Constraints**: 
- Must maintain backward compatibility with existing direct repository usage
- Must support at least 3 popular HTTP routers (gorilla/mux, gin, echo)
- Input validation and sanitization for security
- Comprehensive error handling with typed errors

**Scale/Scope**: 
- Support 95% of common use cases without custom repositories
- Integration complexity: 3 lines of code for basic usage
- Route integration: 2 lines of code for existing applications

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Clean Architecture Compliance ✅
- **Layer Dependencies**: Library initialization will maintain proper layer separation
  - Domain layer remains pure (no external dependencies)
  - Use case layer defines repository interfaces
  - Repository implementations in `internal/repository/`
  - HTTP handlers in `internal/rest/`
- **Interface Declaration**: Repository interfaces declared in use case layer, implemented in repository layer
- **Dependency Inversion**: Library factory will inject dependencies from outer layers to inner layers

### Test-First Development ✅
- **Unit Tests**: Each new component will have unit tests with mocks
- **Integration Tests**: Repository layer tests with real database connections
- **Test Coverage**: Minimum 80% coverage for new library initialization code
- **TDD Cycle**: Red-Green-Refactor for all new functionality

### Package Organization ✅
- **Service-focused packaging**: Library initialization will be in a new `lib/` package
- **Internal isolation**: Repository implementations remain in `internal/`
- **Clean separation**: Library API separate from internal implementation details

### Independence from Frameworks ✅
- **Framework agnostic**: Library initialization will not depend on specific HTTP frameworks
- **Database abstraction**: Uses `database/sql` interface, not specific drivers
- **Router compatibility**: Supports multiple HTTP routers without framework lock-in

### Naming Conventions ✅
- **API Fields**: All library configuration and response fields will use camelCase
- **Database Fields**: All database operations will use snake_case
- **JSON Tags**: Proper camelCase JSON tags for all API responses

### Quality Gates ✅
- **Linting**: All new code will pass `golangci-lint`
- **Security**: Input validation and sanitization implemented
- **Error Handling**: Comprehensive error handling with typed errors
- **Documentation**: All public functions documented

## Constitution Check (Post-Design)

*Re-evaluated after Phase 1 design completion*

### Clean Architecture Compliance ✅
- **Layer Dependencies**: Library initialization maintains proper layer separation
  - New `lib/` package acts as facade layer
  - Domain layer remains pure (no external dependencies)
  - Use case layer defines repository interfaces
  - Repository implementations in `internal/repository/`
  - HTTP handlers in `internal/rest/`
- **Interface Declaration**: Repository interfaces declared in use case layer, implemented in repository layer
- **Dependency Inversion**: Library factory injects dependencies from outer layers to inner layers

### Test-First Development ✅
- **Unit Tests**: Each new component will have unit tests with mocks
- **Integration Tests**: Repository layer tests with real database connections
- **Test Coverage**: Minimum 80% coverage for new library initialization code
- **TDD Cycle**: Red-Green-Refactor for all new functionality

### Package Organization ✅
- **Service-focused packaging**: Library initialization in new `lib/` package
- **Internal isolation**: Repository implementations remain in `internal/`
- **Clean separation**: Library API separate from internal implementation details
- **Backward compatibility**: Existing packages unchanged

### Independence from Frameworks ✅
- **Framework agnostic**: Library initialization does not depend on specific HTTP frameworks
- **Database abstraction**: Uses `database/sql` interface, not specific drivers
- **Router compatibility**: Supports multiple HTTP routers without framework lock-in
- **Configuration injection**: All dependencies injected via configuration

### Naming Conventions ✅
- **API Fields**: All library configuration and response fields use camelCase
- **Database Fields**: All database operations use snake_case
- **JSON Tags**: Proper camelCase JSON tags for all API responses
- **Consistent naming**: All new entities follow established patterns

### Quality Gates ✅
- **Linting**: All new code will pass `golangci-lint`
- **Security**: Input validation and sanitization implemented
- **Error Handling**: Comprehensive error handling with typed errors
- **Documentation**: All public functions documented
- **Performance**: Library initialization < 100ms, minimal memory overhead

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```
lib/                          # New library package for simplified initialization
├── config.go                 # Library configuration structures
├── factory.go                # Factory functions for service initialization
├── errors.go                 # Library-specific error types
└── router.go                 # HTTP router integration utilities

internal/                     # Existing internal implementations (unchanged)
├── repository/
│   ├── mysql/               # MySQL repository implementations
│   └── redis/               # Redis repository implementations
└── rest/                    # HTTP handlers and routing
    ├── handlers/            # HTTP request handlers
    └── middleware/          # HTTP middleware

domain/                      # Domain layer (unchanged)
├── points/
│   ├── entities.go         # Domain entities
│   └── errors.go           # Domain errors

points/                      # Use case layer (unchanged)
├── service.go              # Service interfaces
├── point_type_service.go   # Point type business logic
└── balance_service.go      # Balance business logic

examples/                    # Library usage examples
├── basic_usage.go          # Simple integration example
├── advanced_usage.go       # Full feature example
└── custom_repo_example.go  # Custom repository example
```

**Structure Decision**: Single library package approach. The new `lib/` package will contain the simplified initialization API, while existing `internal/`, `domain/`, and `points/` packages remain unchanged to maintain backward compatibility. This approach allows both simple library usage and advanced custom implementations.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
