# Data Model: Enhanced Library Integration

**Date**: 2024-12-19  
**Feature**: Enhanced Library Integration  
**Branch**: 002-http

## Overview

This document defines the data structures and entities for the library initialization system, including configuration objects, service containers, and error types.

## Core Entities

### LibraryConfig

**Purpose**: Main configuration object for library initialization

**Fields**:
- `DB *sql.DB` - Database connection (required)
- `Redis redis.UniversalClient` - Redis client (required)
- `Repositories *RepositoryConfig` - Custom repository overrides (optional)
- `Routes *RouteConfig` - HTTP route configuration (optional)

**Validation Rules**:
- DB must not be nil
- Redis must not be nil
- If Repositories provided, all fields must implement correct interfaces
- If Routes provided, PathPrefix must be valid URL path

**JSON Tags**: All fields use camelCase
```go
type LibraryConfig struct {
    DB          *sql.DB                    `json:"db,omitempty"`
    Redis       redis.UniversalClient      `json:"redis,omitempty"`
    Repositories *RepositoryConfig         `json:"repositories,omitempty"`
    Routes      *RouteConfig               `json:"routes,omitempty"`
}
```

### RepositoryConfig

**Purpose**: Configuration for custom repository implementations

**Fields**:
- `PointType points.PointTypeRepository` - Custom point type repository (optional)
- `Balance points.BalanceRepository` - Custom balance repository (optional)
- `Ranking points.RankingRepository` - Custom ranking repository (optional)

**Validation Rules**:
- All fields are optional
- If provided, must implement the corresponding interface
- Cannot be nil if provided (use interface{} for optional override)

**JSON Tags**: All fields use camelCase
```go
type RepositoryConfig struct {
    PointType points.PointTypeRepository `json:"pointType,omitempty"`
    Balance   points.BalanceRepository   `json:"balance,omitempty"`
    Ranking   points.RankingRepository   `json:"ranking,omitempty"`
}
```

### RouteConfig

**Purpose**: Configuration for HTTP route integration

**Fields**:
- `PathPrefix string` - URL path prefix for Acto routes (default: "/api/v1")
- `Middleware []MiddlewareFunc` - Custom middleware functions (optional)
- `CORS *CORSConfig` - CORS configuration (optional)

**Validation Rules**:
- PathPrefix must be valid URL path (start with "/")
- Middleware functions must be compatible with target router
- CORS config must be valid if provided

**JSON Tags**: All fields use camelCase
```go
type RouteConfig struct {
    PathPrefix string           `json:"pathPrefix"`
    Middleware []MiddlewareFunc `json:"middleware,omitempty"`
    CORS      *CORSConfig       `json:"cors,omitempty"`
}
```

### CORSConfig

**Purpose**: CORS configuration for HTTP routes

**Fields**:
- `AllowedOrigins []string` - Allowed origins (default: ["*"])
- `AllowedMethods []string` - Allowed HTTP methods (default: ["GET", "POST", "PUT", "DELETE"])
- `AllowedHeaders []string` - Allowed headers (default: ["Content-Type", "Authorization"])
- `AllowCredentials bool` - Allow credentials (default: false)

**Validation Rules**:
- AllowedOrigins cannot be empty
- AllowedMethods must contain valid HTTP methods
- AllowedHeaders must contain valid header names

**JSON Tags**: All fields use camelCase
```go
type CORSConfig struct {
    AllowedOrigins   []string `json:"allowedOrigins"`
    AllowedMethods   []string `json:"allowedMethods"`
    AllowedHeaders   []string `json:"allowedHeaders"`
    AllowCredentials bool     `json:"allowCredentials"`
}
```

### Library

**Purpose**: Main library instance containing all services

**Fields**:
- `Services *Services` - Container for all business services
- `Config LibraryConfig` - Original configuration used for initialization
- `Repositories *RepositoryContainer` - Repository implementations used

**Methods**:
- `Close() error` - Cleanup resources
- `Health() error` - Check service health
- `GetServices() *Services` - Get service container

**JSON Tags**: All fields use camelCase
```go
type Library struct {
    Services    *Services           `json:"services,omitempty"`
    Config      LibraryConfig      `json:"config,omitempty"`
    Repositories *RepositoryContainer `json:"repositories,omitempty"`
}
```

### Services

**Purpose**: Container for all business services

**Fields**:
- `PointType points.PointTypeService` - Point type service
- `Balance points.BalanceService` - Balance service
- `Distribution points.DistributionService` - Distribution service
- `Redemption points.RedemptionService` - Redemption service

**JSON Tags**: All fields use camelCase
```go
type Services struct {
    PointType    points.PointTypeService    `json:"pointType,omitempty"`
    Balance      points.BalanceService      `json:"balance,omitempty"`
    Distribution points.DistributionService  `json:"distribution,omitempty"`
    Redemption  points.RedemptionService   `json:"redemption,omitempty"`
}
```

### RepositoryContainer

**Purpose**: Container for repository implementations

**Fields**:
- `PointType points.PointTypeRepository` - Point type repository
- `Balance points.BalanceRepository` - Balance repository
- `Ranking points.RankingRepository` - Ranking repository

**JSON Tags**: All fields use camelCase
```go
type RepositoryContainer struct {
    PointType points.PointTypeRepository `json:"pointType,omitempty"`
    Balance   points.BalanceRepository   `json:"balance,omitempty"`
    Ranking   points.RankingRepository   `json:"ranking,omitempty"`
}
```

## Error Types

### LibraryError

**Purpose**: Typed error for library operations

**Fields**:
- `Type ErrorType` - Error type classification
- `Message string` - Human-readable error message
- `Cause error` - Underlying error (optional)
- `Context map[string]interface{}` - Additional context (optional)

**JSON Tags**: All fields use camelCase
```go
type LibraryError struct {
    Type    ErrorType                `json:"type"`
    Message string                  `json:"message"`
    Cause   error                   `json:"cause,omitempty"`
    Context map[string]interface{}  `json:"context,omitempty"`
}
```

### ErrorType

**Purpose**: Classification of error types

**Values**:
- `ErrTypeConnection` - Database/Redis connection errors
- `ErrTypeValidation` - Input validation errors
- `ErrTypeConfiguration` - Configuration errors
- `ErrTypeRepository` - Repository operation errors
- `ErrTypeService` - Service operation errors
- `ErrTypeRoute` - Route integration errors

**JSON Tags**: Uses string representation
```go
type ErrorType int

const (
    ErrTypeConnection ErrorType = iota
    ErrTypeValidation
    ErrTypeConfiguration
    ErrTypeRepository
    ErrTypeService
    ErrTypeRoute
)
```

## State Transitions

### Library Initialization

1. **Start**: Empty state
2. **Validate Config**: Check configuration validity
3. **Create Repositories**: Initialize repository implementations
4. **Create Services**: Initialize business services
5. **Ready**: Library ready for use

### Error States

1. **Validation Error**: Configuration invalid
2. **Connection Error**: Database/Redis connection failed
3. **Repository Error**: Repository initialization failed
4. **Service Error**: Service initialization failed

## Relationships

### Library → Services
- One-to-one relationship
- Library contains Services
- Services are created during Library initialization

### Library → RepositoryContainer
- One-to-one relationship
- Library contains RepositoryContainer
- RepositoryContainer is created during Library initialization

### LibraryConfig → RepositoryConfig
- One-to-one optional relationship
- LibraryConfig may contain RepositoryConfig
- RepositoryConfig overrides default repositories

### LibraryConfig → RouteConfig
- One-to-one optional relationship
- LibraryConfig may contain RouteConfig
- RouteConfig configures HTTP route integration

## Validation Rules

### LibraryConfig Validation
- DB must not be nil
- Redis must not be nil
- If Repositories provided, must be valid
- If Routes provided, must be valid

### RepositoryConfig Validation
- All fields are optional
- If provided, must implement correct interface
- Cannot be nil if provided

### RouteConfig Validation
- PathPrefix must be valid URL path
- Middleware functions must be compatible
- CORS config must be valid if provided

## JSON Serialization

### Field Naming
- All fields use camelCase in JSON
- Database fields remain snake_case
- API fields use camelCase

### Omitempty Tags
- Optional fields use `omitempty` tag
- Required fields do not use `omitempty`
- Pointer fields use `omitempty` for nil values

### Example JSON
```json
{
  "db": null,
  "redis": null,
  "repositories": {
    "pointType": null,
    "balance": null,
    "ranking": null
  },
  "routes": {
    "pathPrefix": "/api/v1",
    "middleware": [],
    "cors": {
      "allowedOrigins": ["*"],
      "allowedMethods": ["GET", "POST", "PUT", "DELETE"],
      "allowedHeaders": ["Content-Type", "Authorization"],
      "allowCredentials": false
    }
  }
}
```

## Database Schema

### No New Tables Required
- Library configuration is runtime-only
- No persistent storage needed
- Configuration passed at initialization time

### Existing Tables
- Uses existing database schema
- No modifications to existing tables
- Maintains backward compatibility

## Security Considerations

### Input Validation
- All configuration parameters validated
- SQL injection prevention
- Redis command injection prevention
- Path traversal prevention

### Error Information
- No sensitive information in error messages
- Proper error logging
- Secure connection handling
- Input sanitization

## Performance Considerations

### Memory Usage
- Minimal memory allocation during initialization
- Efficient struct layout
- No unnecessary data duplication

### Initialization Speed
- Fast configuration validation
- Efficient service creation
- Minimal overhead

### Runtime Performance
- No performance impact on service calls
- Efficient dependency injection
- Minimal wrapper overhead
