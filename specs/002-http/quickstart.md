# Quick Start Guide: Enhanced Library Integration

**Date**: 2024-12-19  
**Feature**: Enhanced Library Integration  
**Branch**: 002-http

## Overview

This guide shows how to integrate the Acto points system into your Go application using the new simplified library initialization. The library provides a simple API for common use cases while maintaining flexibility for advanced scenarios.

## Prerequisites

- Go 1.21 or later
- MySQL database
- Redis server
- Basic understanding of Go and HTTP routing

## Installation

```bash
go get github.com/acto/acto
```

## Basic Usage

### 1. Simple Integration (3 lines of code)

```go
package main

import (
    "database/sql"
    "log"
    
    "acto/lib"
    _ "github.com/go-sql-driver/mysql"
    "github.com/redis/go-redis/v9"
)

func main() {
    // 1. Create database connection
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/acto")
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Create Redis client
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // 3. Initialize library
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer library.Close()
    
    // Use the services
    services := library.GetServices()
    
    // Create a point type
    id, err := services.PointType.Create(context.Background(), "gold-points", "金币积分", "用于购买商品的积分")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Created point type: %s", id)
}
```

### 2. HTTP Route Integration (2 lines of code)

```go
package main

import (
    "net/http"
    
    "acto/lib"
    "github.com/gorilla/mux"
)

func main() {
    // Initialize library (same as above)
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Create router
    router := mux.NewRouter()
    
    // Add Acto routes
    err = lib.RegisterGorillaMuxRoutes(router, library.GetServices())
    if err != nil {
        log.Fatal(err)
    }
    
    // Start server
    http.ListenAndServe(":8080", router)
}
```

## Advanced Usage

### 1. Custom Repository Implementation

```go
package main

import (
    "acto/lib"
    "acto/points"
)

// Custom point type repository
type CustomPointTypeRepository struct {
    // Your custom implementation
}

func (r *CustomPointTypeRepository) CreatePointType(ctx context.Context, pt points.PointType) (string, error) {
    // Custom logic here
    return "custom-id", nil
}

// ... implement other methods

func main() {
    // Initialize with custom repository
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
        Repositories: &lib.RepositoryConfig{
            PointType: &CustomPointTypeRepository{},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Use library with custom repository
    services := library.GetServices()
    // ... use services
}
```

### 2. Custom Route Configuration

```go
package main

import (
    "acto/lib"
    "github.com/gorilla/mux"
)

func main() {
    // Initialize library
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
        Routes: &lib.RouteConfig{
            PathPrefix: "/api/v1/points",
            Middleware: []lib.MiddlewareFunc{
                authMiddleware,
                loggingMiddleware,
            },
            CORS: &lib.CORSConfig{
                AllowedOrigins:   []string{"https://example.com"},
                AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
                AllowedHeaders:   []string{"Content-Type", "Authorization"},
                AllowCredentials: true,
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Create router with custom configuration
    router := mux.NewRouter()
    err = lib.RegisterGorillaMuxRoutes(router, library.GetServices())
    if err != nil {
        log.Fatal(err)
    }
}
```

### 3. Multiple Router Support

```go
package main

import (
    "acto/lib"
    "github.com/gin-gonic/gin"
    "github.com/labstack/echo/v4"
)

func main() {
    // Initialize library
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    services := library.GetServices()
    
    // Gin integration
    ginRouter := gin.Default()
    err = lib.RegisterGinRoutes(ginRouter, services)
    if err != nil {
        log.Fatal(err)
    }
    
    // Echo integration
    echoRouter := echo.New()
    err = lib.RegisterEchoRoutes(echoRouter, services)
    if err != nil {
        log.Fatal(err)
    }
}
```

## Error Handling

### 1. Typed Error Handling

```go
package main

import (
    "acto/lib"
    "errors"
)

func main() {
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        var libErr *lib.LibraryError
        if errors.As(err, &libErr) {
            switch libErr.Type {
            case lib.ErrTypeConnection:
                log.Printf("Connection error: %s", libErr.Message)
            case lib.ErrTypeValidation:
                log.Printf("Validation error: %s", libErr.Message)
            case lib.ErrTypeConfiguration:
                log.Printf("Configuration error: %s", libErr.Message)
            default:
                log.Printf("Unknown error: %s", libErr.Message)
            }
        } else {
            log.Printf("Non-library error: %v", err)
        }
        return
    }
}
```

### 2. Health Checking

```go
package main

import (
    "acto/lib"
    "net/http"
)

func main() {
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Check health
    err = library.Health()
    if err != nil {
        log.Printf("Library health check failed: %v", err)
    }
    
    // Health endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        err := library.Health()
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            w.Write([]byte("Unhealthy"))
            return
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Healthy"))
    })
}
```

## Configuration Options

### 1. Library Configuration

```go
type LibraryConfig struct {
    // Required fields
    DB    *sql.DB                    // Database connection
    Redis redis.UniversalClient      // Redis client
    
    // Optional fields
    Repositories *RepositoryConfig    // Custom repository overrides
    Routes      *RouteConfig         // HTTP route configuration
}
```

### 2. Repository Configuration

```go
type RepositoryConfig struct {
    PointType points.PointTypeRepository  // Custom point type repository
    Balance   points.BalanceRepository     // Custom balance repository
    Ranking   points.RankingRepository     // Custom ranking repository
}
```

### 3. Route Configuration

```go
type RouteConfig struct {
    PathPrefix string           // URL path prefix (default: "/api/v1")
    Middleware []MiddlewareFunc // Custom middleware functions
    CORS      *CORSConfig      // CORS configuration
}
```

## Best Practices

### 1. Connection Management

```go
package main

import (
    "context"
    "time"
)

func main() {
    // Configure connection timeouts
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Initialize with context
    library, err := lib.NewLibraryWithContext(ctx, lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer library.Close()
}
```

### 2. Graceful Shutdown

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer library.Close()
    
    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    <-c
    log.Println("Shutting down...")
    
    // Library will be closed by defer
}
```

### 3. Testing

```go
package main

import (
    "testing"
    "acto/lib"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/alicebob/miniredis"
)

func TestLibraryIntegration(t *testing.T) {
    // Mock database
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()
    
    // Mock Redis
    s := miniredis.RunT(t)
    defer s.Close()
    
    redisClient := redis.NewClient(&redis.Options{
        Addr: s.Addr(),
    })
    
    // Initialize library
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        t.Fatal(err)
    }
    defer library.Close()
    
    // Test library functionality
    services := library.GetServices()
    // ... test services
}
```

## Troubleshooting

### 1. Common Issues

**Database Connection Failed**
```go
// Check database connection
err := db.Ping()
if err != nil {
    log.Printf("Database connection failed: %v", err)
}
```

**Redis Connection Failed**
```go
// Check Redis connection
err := redisClient.Ping(context.Background()).Err()
if err != nil {
    log.Printf("Redis connection failed: %v", err)
}
```

**Route Conflicts**
```go
// Check for route conflicts
router := mux.NewRouter()
router.PathPrefix("/api/v1").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Your existing routes
}))

// Acto routes will conflict with existing routes
// Use different path prefix
err := lib.RegisterGorillaMuxRoutes(router, services, lib.RouteConfig{
    PathPrefix: "/acto/api/v1",
})
```

### 2. Debug Mode

```go
package main

import (
    "acto/lib"
    "log"
)

func main() {
    // Enable debug logging
    lib.SetDebug(true)
    
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Debug information will be logged
}
```

## Migration from Direct Usage

### 1. Before (Direct Usage)

```go
// Old way - direct repository usage
ptRepo := repoMysql.NewPointTypeRepository(db)
balanceRepo := repoMysql.NewBalanceTxRepository(db)
rankingRepo := repoRedis.NewRankingRepository(redisClient)

ptService := points.NewPointTypeService(ptRepo)
balanceService := points.NewBalanceService(balanceRepo, rankingRepo, ptRepo)
```

### 2. After (Library Usage)

```go
// New way - library usage
library, err := lib.NewLibrary(lib.LibraryConfig{
    DB:    db,
    Redis: redisClient,
})
if err != nil {
    log.Fatal(err)
}

services := library.GetServices()
// Use services.PointType, services.Balance, etc.
```

## Performance Considerations

### 1. Connection Pooling

```go
package main

import (
    "database/sql"
    "time"
)

func main() {
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

### 2. Memory Usage

```go
package main

import (
    "runtime"
    "time"
)

func main() {
    library, err := lib.NewLibrary(lib.LibraryConfig{
        DB:    db,
        Redis: redisClient,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Monitor memory usage
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            var m runtime.MemStats
            runtime.ReadMemStats(&m)
            log.Printf("Memory usage: %d KB", m.Alloc/1024)
        }
    }()
}
```

## Next Steps

1. **Explore Examples**: Check the `examples/` directory for more usage examples
2. **Custom Repositories**: Implement custom repository interfaces for advanced use cases
3. **Middleware**: Add custom middleware for authentication, logging, etc.
4. **Monitoring**: Integrate with monitoring systems for production use
5. **Testing**: Write comprehensive tests for your integration

## Support

- **Documentation**: [Full API Documentation](./contracts/openapi.yaml)
- **Examples**: [Usage Examples](../examples/)
- **Issues**: [GitHub Issues](https://github.com/acto/acto/issues)
- **Discussions**: [GitHub Discussions](https://github.com/acto/acto/discussions)
