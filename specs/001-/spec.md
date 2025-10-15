# Feature Specification: Points System

**Feature Branch**: `001-`  
**Created**: 2025-10-15  
**Status**: Draft  
**Input**: User description: "积分系统 - 支持多种积分类型、增减积分、排名和奖励兑换"

## Architecture Overview

This points system is designed with dual deployment modes to maximize flexibility:

1. **Standalone HTTP Service**: Can be deployed as an independent microservice with RESTful APIs
2. **Embeddable Library**: Can be imported and used directly as a Go library in other applications

This dual-mode design follows Clean Architecture principles where business logic (domain and use case layers) is completely independent of delivery mechanisms. The same core functionality can be accessed via HTTP endpoints or direct function calls.

## Clarifications

### Session 2025-10-15

- Q: Should the system support decimal points (e.g., 12.5 coins) or only integers for point values? → A: Integer values only - all point types use whole numbers (no decimals)
- Q: How should rankings be updated when point balances change? → A: Real-time updates - rankings recalculate immediately after every point transaction (can use Redis ZSET for efficient implementation)
- Q: What authorization model should the system use for operations? → A: Role-based with two roles - operations staff (admin) and end users, system checks role for privileged operations
- Q: What error response format should the system use? → A: Structured JSON with error codes - errors include code, message, and optional details field
- Q: What database transaction isolation level should be used for balance operations? → A: Read Committed with row-level locking - prevents dirty reads, good concurrency, explicit locks for updates

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Points Configuration Management (Priority: P1)

Operations staff can configure multiple types of points (e.g., Gold Coins, Diamonds, Star Points) with custom names, descriptions, and properties. Each point type can be independently managed and has its own balance tracking.

**Why this priority**: Foundation requirement - without point types configured, no other functionality can work. This is the core data structure that everything else depends on.

**Independent Test**: Operations staff can create a new point type called "Gold Coins", set its display name and description, and immediately see it available in the system for assignment to users.

**Acceptance Scenarios**:

1. **Given** operations staff has access to the configuration panel, **When** they create a new point type with name "Gold Coins" and description "Basic currency for rewards", **Then** the point type is saved and available for use
2. **Given** multiple point types exist, **When** operations staff views the point types list, **Then** all configured point types are displayed with their names and properties
3. **Given** a point type named "Diamonds" exists, **When** operations staff updates its description, **Then** the changes are saved and reflected immediately
4. **Given** a point type has no user balances, **When** operations staff attempts to delete it, **Then** the deletion succeeds and the point type is removed

---

### User Story 2 - Points Balance Management (Priority: P1)

System can increase or decrease user points balances for any configured point type, with transaction history tracking for auditing purposes.

**Why this priority**: Core functionality that enables all point-based interactions. Users must be able to earn and spend points for the system to be useful.

**Independent Test**: Administrator awards 100 Gold Coins to a user, the user's balance increases from 0 to 100, and the transaction appears in the history with timestamp and reason.

**Acceptance Scenarios**:

1. **Given** a user has 0 Gold Coins, **When** the system awards them 100 Gold Coins with reason "Welcome bonus", **Then** their balance becomes 100 and a transaction record is created
2. **Given** a user has 150 Gold Coins, **When** the system deducts 50 Gold Coins with reason "Reward redemption", **Then** their balance becomes 100 and a deduction record is created
3. **Given** a user has 30 Gold Coins, **When** the system attempts to deduct 50 Gold Coins, **Then** the operation fails with insufficient balance error and the balance remains 30
4. **Given** a user has transactions for Gold Coins, **When** they view their transaction history, **Then** all transactions are displayed in reverse chronological order with amounts, reasons, and timestamps
5. **Given** a user has multiple point types (Gold Coins and Diamonds), **When** operations are performed on Gold Coins, **Then** only Gold Coins balance changes and Diamond balance remains unchanged

---

### User Story 3 - Points Ranking System (Priority: P2)

System can generate rankings of users based on their point balances for any point type, supporting leaderboards and competitive features.

**Why this priority**: Enables engagement and gamification, but the system is functional without it. Users can still earn and spend points.

**Independent Test**: After multiple users have Gold Coins balances, the system generates a ranking showing users ordered by balance (highest first), and the ranking updates when balances change.

**Acceptance Scenarios**:

1. **Given** three users have Gold Coins balances of 500, 300, and 800, **When** a ranking is generated, **Then** users are ordered as: 1st (800), 2nd (500), 3rd (300)
2. **Given** a ranking exists for Gold Coins, **When** a user's balance increases from 500 to 900, **Then** the ranking updates and they move to 1st place
3. **Given** multiple users have the same balance, **When** a ranking is generated, **Then** users with equal balances receive the same rank (e.g., both rank #2) and the next rank skips accordingly (next is #4)
4. **Given** rankings exist for both Gold Coins and Diamonds, **When** a user views rankings, **Then** they can see separate rankings for each point type

---

### User Story 4 - Rank-Based Reward Distribution (Priority: P2)

System can automatically distribute rewards to users based on their ranking positions, with configurable reward rules for different rank ranges.

**Why this priority**: Adds value through automation and enables competitive campaigns, but manual reward distribution can be used as a workaround initially.

**Independent Test**: Configure a reward rule "Top 3 users get 1000 bonus Gold Coins", execute the distribution, and verify the top 3 users receive the bonus while others don't.

**Acceptance Scenarios**:

1. **Given** a reward rule "Rank 1-3 receives 1000 Gold Coins", **When** the distribution is executed, **Then** users ranked 1st, 2nd, and 3rd each receive 1000 Gold Coins
2. **Given** multiple reward tiers (1st place: 1000, 2nd-5th: 500, 6th-10th: 200), **When** distribution executes, **Then** each user receives rewards according to their rank tier
3. **Given** a reward distribution has already been executed for a ranking period, **When** attempting to execute again for the same period, **Then** the system prevents duplicate distribution
4. **Given** reward distribution is scheduled for a future time, **When** the scheduled time arrives, **Then** the system automatically executes the distribution

---

### User Story 5 - Points Redemption System (Priority: P1)

Users can exchange their points for rewards from a redemption catalog, with real-time balance deduction and redemption confirmation.

**Why this priority**: Core user-facing functionality that provides immediate value - users want to use their earned points. This completes the earn-spend cycle.

**Independent Test**: User with 500 Gold Coins can browse redemption catalog, select a 300 Gold Coins reward, complete redemption, and their balance decreases to 200 with a redemption record created.

**Acceptance Scenarios**:

1. **Given** a user has 500 Gold Coins, **When** they redeem a reward costing 300 Gold Coins, **Then** their balance decreases to 200 and a redemption record is created
2. **Given** a user has 100 Gold Coins, **When** they attempt to redeem a reward costing 300 Gold Coins, **Then** the redemption is rejected with insufficient balance error
3. **Given** a reward requires 100 Gold Coins and 50 Diamonds, **When** a user with both balances redeems it, **Then** both point types are deducted accordingly
4. **Given** multiple users attempt to redeem a limited quantity reward simultaneously, **When** the quantity is exhausted, **Then** only users who completed redemption before stock-out succeed
5. **Given** a user completes a redemption, **When** they view their redemption history, **Then** the redemption appears with reward details, cost, and timestamp

---

### Edge Cases

- What happens when a user's balance would become negative after a deduction operation?
  - System rejects the operation with error code INSUFFICIENT_BALANCE and returns structured error with current balance
- How does the system handle concurrent point operations for the same user?
  - Use transaction locking to ensure balance consistency
- What happens when a point type is deleted that users have balances for?
  - System prevents deletion if active balances exist; returns error code CANNOT_DELETE_ACTIVE_POINT_TYPE
- How does ranking handle users with 0 balance?
  - Users with 0 balance are included in ranking but listed after all positive balances
- What happens when reward inventory runs out during redemption?
  - Transaction fails atomically with error code REWARD_OUT_OF_STOCK - points are not deducted if reward cannot be allocated
- How does system handle ranking ties when distributing rank-based rewards?
  - Users with tied ranks all receive the reward for that rank tier
- What happens when a scheduled reward distribution fails?
  - System retries and logs the failure with error code DISTRIBUTION_FAILED; administrator can manually trigger retry
- How are fractional points handled?
  - All point values are integers only; fractional points are not supported

## Requirements *(mandatory)*

### Functional Requirements

**Deployment & Integration**
- **FR-001**: System MUST be usable as a standalone HTTP service with RESTful APIs
- **FR-002**: System MUST be usable as an embeddable Go library with direct function calls
- **FR-003**: Core business logic MUST be shared between HTTP service and library modes
- **FR-004**: Library mode MUST expose clean, documented Go interfaces for all operations
- **FR-005**: Service mode MUST provide comprehensive API documentation (OpenAPI/Swagger)

**Point Type Management**
- **FR-006**: System MUST allow operations staff to create multiple point types with unique identifiers, display names, and descriptions
- **FR-007**: System MUST support enabling/disabling point types without data loss
- **FR-008**: System MUST prevent deletion of point types that have active user balances
- **FR-009**: System MUST validate point type names are unique within the system

**Points Balance Operations**
- **FR-010**: System MUST support increasing user points balance with a transaction reason
- **FR-011**: System MUST support decreasing user points balance with a transaction reason
- **FR-012**: System MUST reject deduction operations that would result in negative balance
- **FR-013**: System MUST record all balance changes with timestamp, amount, reason, and operation type (credit/debit)
- **FR-014**: System MUST maintain separate balances for each point type per user
- **FR-015**: System MUST ensure balance consistency during concurrent operations using appropriate locking mechanisms
- **FR-016**: System MUST store and process all point values as integers (no decimal support)

**Transaction History**
- **FR-017**: System MUST allow users to view their complete transaction history for each point type
- **FR-018**: System MUST display transactions in reverse chronological order (newest first)
- **FR-019**: System MUST include transaction details: timestamp, amount, balance before/after, reason, and operation type
- **FR-020**: System MUST support filtering transaction history by point type, date range, and operation type

**Ranking System**
- **FR-021**: System MUST generate rankings based on point balances for any point type
- **FR-022**: System MUST order rankings by balance in descending order (highest first)
- **FR-023**: System MUST assign the same rank to users with identical balances
- **FR-024**: System MUST support pagination for large ranking lists
- **FR-025**: System MUST update rankings in real-time immediately after point balance changes
- **FR-026**: System MUST maintain separate rankings for each point type

**Reward Distribution**
- **FR-027**: System MUST allow configuration of rank-based reward rules (rank ranges and reward amounts)
- **FR-028**: System MUST distribute rewards according to configured rules for specific ranking snapshots
- **FR-029**: System MUST prevent duplicate reward distributions for the same ranking period
- **FR-030**: System MUST support scheduling automatic reward distributions
- **FR-031**: System MUST log all reward distribution operations with affected users and amounts
- **FR-032**: System MUST handle distribution failures with retry capabilities

**Points Redemption**
- **FR-033**: System MUST allow users to browse available redemption rewards
- **FR-034**: System MUST display point costs for each reward (supporting multiple point types)
- **FR-035**: System MUST validate user has sufficient balance of all required point types before redemption
- **FR-036**: System MUST deduct point costs atomically upon successful redemption
- **FR-037**: System MUST create redemption records with reward details, costs, and timestamp
- **FR-038**: System MUST handle limited quantity rewards with proper inventory management
- **FR-039**: System MUST prevent redemption when reward inventory is exhausted
- **FR-040**: System MUST allow users to view their redemption history

**Data Integrity & Concurrency**
- **FR-041**: System MUST ensure all balance operations are atomic (either fully complete or fully rollback)
- **FR-042**: System MUST handle concurrent operations on the same user balance safely using row-level locking
- **FR-043**: System MUST validate all operations have proper authorization based on user role
- **FR-044**: System MUST maintain audit trails for all critical operations
- **FR-045**: System MUST use Read Committed isolation level with explicit row-level locks (SELECT FOR UPDATE) for balance updates

**Authorization & Security**
- **FR-046**: System MUST support two user roles: operations staff (admin) and end users
- **FR-047**: System MUST restrict point type management operations to operations staff role only
- **FR-048**: System MUST restrict reward rule configuration and distribution operations to operations staff role only
- **FR-049**: System MUST allow end users to view their own balances, transaction history, and redemption history
- **FR-050**: System MUST allow end users to redeem rewards using their own points
- **FR-051**: System MUST prevent users from viewing or modifying other users' point balances

**Error Handling**
- **FR-052**: System MUST return structured error responses with error code, message, and optional details
- **FR-053**: System MUST use consistent error codes across all operations for the same error types
- **FR-054**: System MUST provide clear, actionable error messages for common failure scenarios
- **FR-055**: HTTP service mode MUST map error codes to appropriate HTTP status codes (400, 403, 404, 409, 500)
- **FR-056**: Library mode MUST return typed errors that can be programmatically checked by calling code

### Key Entities

- **Point Type**: Represents a category of points (e.g., Gold Coins, Diamonds, Star Points)
  - Attributes: unique identifier, display name, description, enabled status, creation timestamp
  - Relationships: multiple users can have balances of this point type

- **User Points Balance**: Represents a user's current balance for a specific point type
  - Attributes: user identifier, point type, current balance (integer), last updated timestamp
  - Relationships: belongs to one user, belongs to one point type
  - Constraints: balance cannot be negative, balance must be an integer

- **Points Transaction**: Records a single increase or decrease operation
  - Attributes: transaction ID, user identifier, point type, amount (integer), operation type (credit/debit), reason, balance before (integer), balance after (integer), timestamp
  - Relationships: belongs to one user, belongs to one point type
  - Business rules: immutable once created (audit trail), all amounts must be integers

- **Ranking**: Represents users ordered by their point balance for a specific point type
  - Attributes: point type, snapshot timestamp, ranking entries
  - Relationships: references multiple users and their balances at snapshot time
  - Note: Updated in real-time as balances change; can be efficiently implemented using sorted data structures (e.g., Redis ZSET)

- **Ranking Entry**: Individual entry in a ranking
  - Attributes: rank position, user identifier, point balance, tied rank indicator
  - Relationships: belongs to one ranking, references one user

- **Reward Rule**: Configuration for rank-based reward distribution
  - Attributes: rule ID, point type, rank range (min-max), reward amount, reward point type, active status
  - Relationships: defines rewards for a specific point type's ranking
  - Example: "Rank 1-3 receives 1000 Gold Coins"

- **Reward Distribution**: Records execution of reward distribution
  - Attributes: distribution ID, execution timestamp, ranking snapshot reference, reward rule applied, status (pending/completed/failed)
  - Relationships: references reward rules used, links to transaction records for distributed rewards
  - Constraints: one distribution per ranking period to prevent duplicates

- **Redemption Reward**: Item or benefit available for point redemption
  - Attributes: reward ID, name, description, point costs (can be multiple point types), available quantity, total redeemed count, enabled status
  - Relationships: can require multiple point types for redemption
  - Example: "Premium Gift Box" costs 500 Gold Coins + 100 Diamonds

- **Redemption Record**: Records a user's redemption of a reward
  - Attributes: redemption ID, user identifier, reward details, point costs deducted, redemption timestamp, status (completed/pending/cancelled)
  - Relationships: belongs to one user, references the redeemed reward, links to points transaction records
  - Business rules: immutable once completed

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Functionality & Completeness**
- **SC-001**: Operations staff can create and configure at least 5 different point types without system limitations
- **SC-002**: System successfully processes 1,000 concurrent point transactions without data inconsistency
- **SC-003**: Users can complete a point redemption transaction from browsing to confirmation in under 30 seconds
- **SC-004**: Point balance operations (credit/debit) complete in under 100 milliseconds for 95% of requests
- **SC-005**: Library mode: Go packages can be imported and used directly with same functionality as HTTP service
- **SC-006**: Service mode: RESTful HTTP APIs provide complete access to all point system features

**Data Integrity & Reliability**
- **SC-007**: Zero balance inconsistencies occur during concurrent operations across 10,000 test transactions
- **SC-008**: All point transactions are successfully recorded in transaction history with 100% accuracy
- **SC-009**: System maintains data integrity during reward distribution to 1,000+ users simultaneously

**Ranking Performance**
- **SC-010**: Rankings for up to 100,000 users can be generated in under 5 seconds
- **SC-011**: Ranking updates are reflected in real-time (within 100ms) after point balance changes
- **SC-012**: System correctly handles ranking ties with proper rank assignment for all tied users

**User Experience**
- **SC-013**: 90% of users successfully complete their first redemption without errors or confusion
- **SC-014**: Users can view their complete transaction history (100+ transactions) in under 2 seconds
- **SC-015**: Redemption failures due to insufficient balance provide clear, actionable error messages

**Business Metrics**
- **SC-016**: System supports at least 3 different point types in production use simultaneously
- **SC-017**: Automated reward distributions execute successfully on schedule with 99% reliability
- **SC-018**: Support tickets related to point balance discrepancies are reduced by 80% compared to manual systems

**Scalability**
- **SC-019**: System handles 10,000 daily active users with peak loads of 500 concurrent transactions
- **SC-020**: Database can store transaction history for 1 million transactions without performance degradation
- **SC-021**: Redemption catalog supports at least 100 active rewards without impacting browse performance
- **SC-022**: Library mode: Can be embedded in applications with minimal memory overhead (< 50MB additional RAM)
- **SC-023**: Service mode: Supports horizontal scaling with multiple instances sharing the same database

## Assumptions

1. **Dual Deployment Modes**: System designed to work both as a standalone HTTP service and as an embeddable Go library
   - **Library Mode**: Core business logic exposed as Go packages that can be imported and called directly
   - **Service Mode**: HTTP REST API wrapping the library functionality for remote access
   - Both modes share the same business logic implementation
2. **User Authentication**: Assumes an existing user authentication system provides user identifiers; points system does not handle user login/registration
3. **Currency Precision**: All point values are integers only; no decimal or fractional points are supported
4. **Time Zones**: All timestamps stored in UTC; display formatting handled by client applications
5. **Database Transactions**: Uses Read Committed isolation level with explicit row-level locking for balance operations to ensure consistency
6. **Notification System**: Reward distributions and redemptions may trigger notifications, but notification delivery is handled by a separate system
5. **Database Transactions**: Uses Read Committed isolation level with explicit row-level locking for balance operations to ensure consistency
6. **Notification System**: Reward distributions and redemptions may trigger notifications, but notification delivery is handled by a separate system
7. **Payment Integration**: Redemption rewards are fulfilled by external systems; points system only tracks the redemption transaction
8. **Ranking Periods**: Ranking snapshots for reward distribution are created manually or via scheduled jobs; ranking period definition (daily/weekly/monthly) is configurable
9. **Access Control**: Two user roles are supported - operations staff (admin role) with full system access, and end users with access to their own data only
10. **Data Retention**: Transaction history and redemption records are retained indefinitely for audit purposes unless specific retention policies are defined
11. **Internationalization**: Point type names and reward descriptions support Unicode for multiple languages; specific i18n implementation is outside scope

## Dependencies

- **User Management System**: Required for user identification and authentication
- **Authorization System**: Required for operations staff access control
- **Notification Service** (optional): For sending alerts about rewards, low balance, redemption confirmations
- **Inventory Management** (optional): For tracking physical reward inventory if redemptions involve physical items
- **Scheduling Service** (optional): For automated reward distributions at specified times

## Out of Scope

1. **User Authentication**: Login, registration, password management (system accepts authenticated user IDs)
2. **API Gateway Features**: Rate limiting, API key management, request throttling (can be added externally in service mode)
3. **Campaign Management**: Creating marketing campaigns that award points automatically
4. **Payment Processing**: Purchasing points with real money
5. **Fraud Detection**: Identifying suspicious point earning or redemption patterns
6. **Multi-tenancy**: Supporting multiple organizations with isolated point systems
7. **Point Expiration**: Automatic expiration of points after a certain period (future enhancement)
8. **Point Transfers**: Transferring points between users (future enhancement)
9. **Social Features**: Sharing rankings, challenging friends, social leaderboards
10. **Analytics Dashboard**: Detailed reporting and analytics on point usage patterns (raw data available for external analytics)
11. **Mobile Apps**: Native mobile applications (HTTP APIs enable future mobile integration)
12. **GraphQL API**: Only RESTful HTTP APIs provided in service mode (GraphQL can be added as a separate layer)
