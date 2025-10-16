# Integration Test Scaffold

This document outlines how to write an integration-like test for the library initialization and basic routes using sqlmock (for DB) and miniredis (for Redis).

## Example (pseudo-Go)

```go
package lib_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "acto/lib"
    sqlmock "github.com/DATA-DOG/go-sqlmock"
    miniredis "github.com/alicebob/miniredis/v2"
    goRedis "github.com/redis/go-redis/v9"
)

func TestLibrary_Health_Integration(t *testing.T) {
    // sqlmock
    db, _, err := sqlmock.New()
    if err != nil { t.Fatal(err) }
    defer db.Close()

    // miniredis
    s, _ := miniredis.Run()
    defer s.Close()
    rc := goRedis.NewClient(&goRedis.Options{Addr: s.Addr()})

    // init library
    library, err := lib.NewLibrary(lib.LibraryConfig{DB: db, Redis: rc})
    if err != nil { t.Fatal(err) }

    // health route via fake mux
    fr := &fakeRouter{}
    _ = lib.RegisterGorillaMuxRoutes(fr, library)
    rr := httptest.NewRecorder()
    req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
    fr.sub.handlers["/health"].ServeHTTP(rr, req)
    if rr.Code != http.StatusOK { t.Fatalf("want 200, got %d", rr.Code) }
}
```

Notes:
- The fake mux types in `lib/router_test.go` can be copied for route testing without importing Gorilla.
- Use `sqlmock` only for connection creation; repository calls are not required for the health path.
- `miniredis` provides an in-memory Redis server; `Ping` will succeed.
