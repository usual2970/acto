# Feature Specification: Enhanced Library Integration

**Feature Branch**: `002-http`  
**Created**: 2024-12-19  
**Status**: Draft  
**Input**: User description: "当前本项目已经基本可用了，可以作为独立的http服务，也可以作为类库使用。"

## Clarifications

### Session 2024-12-19

- Q: What level of security validation and error handling should be included in the library initialization? → A: Include basic security validation (connection validation, input sanitization) and comprehensive error handling with specific error types

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Simplified Library Initialization (Priority: P1)

Developers can integrate the Acto points system into their applications with minimal setup by providing only database and Redis connections, without needing to implement repository interfaces themselves.

**Why this priority**: This addresses the primary pain point of complex library usage and enables rapid adoption by reducing integration complexity from hours to minutes.

**Independent Test**: Can be fully tested by creating a new Go project, importing the library, and initializing services with just database and Redis connections, then performing basic operations like creating point types and managing balances.

**Acceptance Scenarios**:

1. **Given** a developer wants to use Acto as a library, **When** they provide database and Redis connections, **Then** they can immediately start using all point management features
2. **Given** a developer has initialized the library, **When** they call service methods, **Then** operations work without requiring custom repository implementations

---

### User Story 2 - Custom Repository Support (Priority: P2)

Developers can still implement custom repository interfaces when they need specialized data storage logic or want to use different database systems.

**Why this priority**: Maintains flexibility for advanced use cases while providing simple defaults for common scenarios.

**Independent Test**: Can be tested by implementing a custom repository that logs all operations, then verifying that the custom repository is used instead of the default implementation.

**Acceptance Scenarios**:

1. **Given** a developer wants custom data storage behavior, **When** they implement repository interfaces, **Then** their custom implementations are used instead of defaults
2. **Given** a developer has custom repositories, **When** they initialize services, **Then** the system uses their implementations without conflicts

---

### User Story 3 - Route Integration (Priority: P1)

Developers can integrate Acto's HTTP routes into their existing web applications as a subset of their routing system.

**Why this priority**: Enables seamless integration with existing web applications without requiring separate HTTP servers or complex routing configurations.

**Independent Test**: Can be tested by creating a web application with existing routes, then adding Acto routes as a sub-router, and verifying that both application routes and Acto routes work correctly.

**Acceptance Scenarios**:

1. **Given** a developer has an existing web application, **When** they integrate Acto routes, **Then** both their routes and Acto routes work without conflicts
2. **Given** a developer has integrated Acto routes, **When** they make HTTP requests to Acto endpoints, **Then** they receive proper responses with correct data formats

---

### Edge Cases

- What happens when database connection fails during initialization?
- How does the system handle Redis connection failures when ranking features are used?
- What happens when custom repositories are partially implemented?
- How does route integration handle path conflicts with existing routes?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a simple initialization function that accepts database and Redis connections and returns configured services
- **FR-002**: System MUST allow developers to override default repository implementations with custom ones
- **FR-003**: System MUST provide a route registration function that integrates Acto routes into existing HTTP routers
- **FR-004**: System MUST maintain backward compatibility with existing direct repository usage patterns
- **FR-005**: System MUST handle missing dependencies gracefully with clear error messages
- **FR-006**: System MUST support both default and custom repository implementations simultaneously
- **FR-007**: System MUST provide route mounting that doesn't interfere with existing application routes
- **FR-008**: System MUST validate database and Redis connections during initialization with specific error types
- **FR-009**: System MUST sanitize all input parameters to prevent injection attacks
- **FR-010**: System MUST provide comprehensive error handling with typed error responses for all failure scenarios

### Key Entities

- **LibraryConfig**: Configuration object containing database connection, Redis client, and optional custom repositories
- **RouteConfig**: Configuration for HTTP route integration including path prefixes and middleware options
- **ServiceContainer**: Container object that holds all initialized services (point types, balances, rankings, etc.)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can integrate Acto library into a new project in under 5 minutes
- **SC-002**: Library initialization requires no more than 3 lines of code for basic usage
- **SC-003**: Route integration adds no more than 2 lines of code to existing applications
- **SC-004**: 95% of common use cases can be handled without custom repository implementations
- **SC-005**: Custom repository implementations can be added without modifying existing code
- **SC-006**: Route integration works with at least 3 popular Go HTTP router libraries (gorilla/mux, gin, echo)