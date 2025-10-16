<!--
Sync Impact Report:
Version change: 1.0.0 → 1.1.0
Modified principles: None
Added sections: VII. Naming Conventions (NON-NEGOTIABLE)
Removed sections: None
Templates requiring updates:
  ✅ .specify/templates/plan-template.md (no changes needed)
  ✅ .specify/templates/spec-template.md (no changes needed)  
  ✅ .specify/templates/tasks-template.md (no changes needed)
  ✅ .specify/templates/commands/*.md (no changes needed)
  ✅ README.md (no changes needed)
Follow-up TODOs: None
-->

# Acto Clean Architecture Constitution

## Core Principles

### I. Clean Architecture Layers (NON-NEGOTIABLE)
**The Dependency Rule**: Source code dependencies must only point inwards. Inner layers cannot know anything about outer layers.

**Four-Layer Structure**:
1. **Domain Layer (Entities)**: Enterprise business rules, models, and core business logic
   - Pure business objects with no external dependencies
   - Defines interfaces for repositories and use cases
   - Contains core domain models and business entities
   
2. **Use Case Layer (Business Logic)**: Application-specific business rules
   - Orchestrates data flow between entities and repositories
   - Contains application business logic
   - Defines interfaces that consuming layers will use
   - Independent of UI, database, or external frameworks

3. **Repository Layer (Data Access)**: Interface adapters for data sources
   - Implements repository interfaces defined in domain/use case layers
   - Handles database queries, external API calls, caching
   - Converts data between external formats and domain entities
   
4. **Delivery Layer (Controllers/Handlers)**: Interface adapters for input/output
   - HTTP handlers, CLI commands, gRPC services
   - Converts external requests to use case calls
   - Formats responses for external consumption

### II. Interface Declaration at Consumer Side
**Declare interfaces where they are used, not where they are implemented.**

- Use case layer defines repository interfaces it needs
- Delivery layer defines use case interfaces it needs
- Implementation details live in outer layers
- Promotes loose coupling and easier testing
- Enables dependency inversion principle

Example:
```go
// In usecase package
type ArticleRepository interface {
    GetByID(ctx context.Context, id int64) (*domain.Article, error)
}

// In repository package (implements the interface)
type mysqlArticleRepository struct { ... }
```

### III. Test-First Development (NON-NEGOTIABLE)
**TDD Mandatory**: Tests must be written and approved before implementation.

- **Unit Tests**: Every layer must have unit tests with mocks
  - Domain layer: Test business logic in isolation
  - Use case layer: Mock repositories and test business flows
  - Repository layer: Test data access logic (integration tests)
  - Delivery layer: Mock use cases and test request/response handling

- **Red-Green-Refactor Cycle**: 
  1. Write failing test
  2. Implement minimal code to pass
  3. Refactor while keeping tests green

- **Test Coverage**: Minimum 80% coverage for critical business logic
- **Mocking**: Use mockery or similar tools to generate mocks from interfaces

### IV. Package Organization & Internal Structure
**Service-focused packaging with internal isolation**

Based on the actual structure from [bxcodec/go-clean-arch](https://github.com/bxcodec/go-clean-arch):

```
project-root/
├── app/                    # Application entry point
│   └── main.go            # Bootstrap and dependency injection
│
├── domain/                # Domain layer (innermost, no dependencies)
│   ├── article.go         # Article entity and business rules
│   ├── author.go          # Author entity and business rules
│   └── errors.go          # Domain-specific error definitions
│
├── article/               # Use Case layer (business logic)
│   ├── service.go         # Article use case implementation
│   ├── service_test.go    # Use case unit tests
│   └── mocks/             # Generated mocks for testing
│
├── internal/              # Internal implementations (hidden from external imports)
│   ├── repository/        # Repository layer implementations
│   │   ├── mysql/         # MySQL-specific implementations
│   │   │   ├── article.go      # Article repository for MySQL
│   │   │   ├── article_test.go # Repository integration tests
│   │   │   ├── author.go       # Author repository for MySQL
│   │   │   └── author_test.go  # Repository integration tests
│   │   └── helper.go      # Repository helper functions
│   │
│   ├── rest/              # Delivery layer (HTTP handlers)
│   │   ├── article.go     # Article HTTP handlers
│   │   ├── article_test.go # Handler tests
│   │   ├── middleware/    # HTTP middlewares (CORS, auth, logging)
│   │   └── mocks/         # Generated mocks for handler testing
│   │
│   └── workers/           # Background workers (if needed)
│       └── ...            # Cron jobs, message consumers, etc.
│
├── misc/                  # Miscellaneous files
│   └── make/              # Makefile includes
│
├── .air.toml              # Hot reload configuration
├── .golangci.yaml         # Linter configuration
├── compose.yaml           # Docker Compose for local development
├── Dockerfile             # Container image definition
├── example.env            # Environment variable template
├── Makefile               # Build and development tasks
├── go.mod                 # Go module definition
└── article.sql            # Database schema
```

**Key Structural Principles**:

1. **`domain/` package**: Pure business entities and rules, no external dependencies
   - Contains only domain models (Article, Author)
   - Defines domain-specific errors
   - No framework imports, no database logic
   
2. **Service packages** (e.g., `article/`): Use case implementations
   - One package per domain service/aggregate
   - Contains business logic orchestration
   - Defines interfaces it needs (repository, external services)
   - Includes mocks for testing
   
3. **`internal/` package**: Implementation details hidden from external projects
   - Cannot be imported by external Go projects (enforced by Go compiler)
   - Contains all infrastructure code
   
4. **`internal/repository/`**: Data access implementations
   - Organized by data source type (`mysql/`, `postgres/`, `redis/`, etc.)
   - Each file implements interfaces defined in service layer
   - Integration tests alongside implementation
   
5. **`internal/rest/`**: HTTP delivery layer
   - HTTP handlers and routing
   - Request/response DTOs and validation
   - Middleware for cross-cutting concerns
   - Implements interfaces needed by use cases
   
6. **`app/`**: Application bootstrap
   - Single `main.go` file
   - Dependency injection and wiring
   - Configuration loading
   - Server initialization

**Package Naming Rules**:
- Use lowercase, singular nouns (article, not articles)
- Avoid generic names (util, common, helper) at root level
- Package name should match directory name
- Keep packages flat, avoid deep nesting

### V. Independence from Frameworks
**Frameworks are tools, not architectures**

- Business logic must not depend on web frameworks (Echo, Gin, etc.)
- Database logic must be swappable (MySQL ↔ PostgreSQL ↔ MongoDB)
- Delivery mechanisms can change without affecting business rules
- Use interfaces to abstract external dependencies
- Configuration should be injectable and environment-based

### VI. Testability Above All
**The architecture exists to enable testing**

- Every component must be testable in isolation
- UI can be tested without the backend
- Business rules can be tested without UI or database
- Database can be mocked for unit testing
- Integration tests for repository layer with real databases
- End-to-end tests for critical user journeys

### VII. Naming Conventions (NON-NEGOTIABLE)
**Consistent naming across all layers and interfaces**

**API Field Naming**:
- All API request and response fields MUST use camelCase (first letter lowercase)
- JSON tags MUST use camelCase format: `json:"fieldName"`
- Examples: `userId`, `pointTypeId`, `createdAt`, `displayName`

**Database Field Naming**:
- All database columns MUST use snake_case format
- Examples: `user_id`, `point_type_id`, `created_at`, `display_name`

**Rationale**: 
- API camelCase ensures consistency with JavaScript/frontend conventions
- Database snake_case follows SQL naming best practices and improves readability
- Clear separation between API and database naming prevents confusion
- JSON serialization/deserialization handles the conversion automatically

## Architecture Standards

### Dependency Management
- **Go Modules**: Strict version control with `go.mod`
- **No circular dependencies**: Enforce through design and code review
- **Minimal external dependencies**: Evaluate necessity before adding
- **Vendor when necessary**: For stability in production environments

### Error Handling
- **Explicit error returns**: No panics in business logic
- **Context propagation**: Pass context through all layers
- **Error wrapping**: Use `fmt.Errorf` with `%w` for error chains
- **Domain-specific errors**: Define custom error types in domain layer
- **HTTP error codes**: Map domain errors to appropriate HTTP status in delivery layer

### Data Flow
1. **Request** → Delivery Layer (validates input, calls use case)
2. **Use Case** → Orchestrates business logic, calls repositories
3. **Repository** → Fetches/persists data, returns domain models
4. **Domain Models** → Flow back through layers
5. **Response** → Delivery layer formats and returns to client

### API Design
- **RESTful conventions**: Follow REST principles for HTTP APIs
- **JSON by default**: Support other formats when needed
- **Versioning**: Use URL versioning (`/api/v1/`) for breaking changes
- **Pagination**: Implement for list endpoints (limit, offset, cursor-based)
- **Error responses**: Consistent error format across all endpoints

## Development Workflow

### Code Review Requirements
- **Architecture compliance**: Verify adherence to layer boundaries
- **Test coverage**: All new code must include tests
- **Interface design**: Review interface contracts before implementation
- **Documentation**: Public functions and types must be documented
- **No TODO comments**: Resolve or create issues instead

### Quality Gates
1. **Unit tests pass**: `make test`
2. **Linting passes**: `golangci-lint run` with no errors
3. **Coverage threshold**: Minimum 80% for new code
4. **Integration tests pass**: Repository layer tests with real database
5. **No security vulnerabilities**: Run `go list -json -m all | nancy sleuth`

### Development Tools
- **Air**: Hot reload for development (`make run` with air)
- **Docker Compose**: Local development environment
- **Mockery**: Generate mocks from interfaces
- **golangci-lint**: Comprehensive linting
- **Make**: Task automation and build scripts

### Git Workflow
- **Feature branches**: Create from `main` for new features
- **Conventional commits**: Use conventional commit format
  - `feat:` for new features
  - `fix:` for bug fixes
  - `refactor:` for code refactoring
  - `test:` for test additions/changes
  - `chore:` for maintenance tasks
- **PR reviews**: Minimum one approval required
- **Squash merges**: Keep main branch history clean

## Governance

### Constitution Authority
- This constitution supersedes all other development practices
- All code must comply with these principles
- Exceptions require explicit documentation and team approval

### Amendment Process
1. Propose amendment with rationale
2. Team discussion and consensus building
3. Update constitution with version bump
4. Create migration plan for existing code if needed
5. Communicate changes to all team members

### Compliance Verification
- All PRs must pass automated checks (tests, linting)
- Code reviews must verify architectural compliance
- Regular architecture reviews for major features
- Refactoring sessions to address technical debt

### Continuous Improvement
- Learn from Uncle Bob's Clean Architecture principles
- Study real-world implementations (like bxcodec/go-clean-arch)
- Adapt patterns to project needs, not blindly follow dogma
- Simplicity over complexity: YAGNI (You Aren't Gonna Need It)

**Version**: 1.1.0 | **Ratified**: 2024-12-19 | **Last Amended**: 2024-12-19