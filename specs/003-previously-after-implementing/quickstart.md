# Quickstart: DI with uber-go/dig

## Goals
- Centralize dependency wiring
- Keep handlers framework-agnostic
- Enable easy test containers

## Architecture

The DI system is organized into clear layer modules:

- **ConfigModule**: Configuration loading
- **InfraModule**: Infrastructure dependencies (DB, Redis)
- **RepoModule**: Repository implementations
- **ServiceModule**: Business services
- **DeliveryModule**: HTTP handlers

## Steps
1. Create providers in `lib/container/modules.go` grouped by layer
2. Add them to `lib/container/builder.go` using dig.Provide
3. In `app/main.go`, build container and register routes via adapter (no manual constructors in routes)

## Module Usage

```go
// Build full container
c, err := container.Build()

// Build test container with overrides
c, err := container.BuildTest(func(c *dig.Container) error {
    return c.Provide(fakeDBProvider)
})
```

## Testing
- Build a test container replacing infra providers with fakes (see `lib/container/testing_test.go`)
- Invoke services/handlers via dig.Invoke in tests
- Test modules in isolation using `TestModuleIsolation`

## Framework-Agnostic Delivery

The system uses a `RouteRegistrar` interface that allows swapping HTTP frameworks:

```go
// Current: gorilla/mux
adapter := lib.NewMuxAdapter(muxRouter)
lib.RegisterRoutes(adapter, "/api/v1", library)

// Alternative: chi router
type ChiAdapter struct { r *chi.Mux }
func (c ChiAdapter) Handle(method, path string, h http.Handler) {
    c.r.Method(method, path, h)
}
adapter := ChiAdapter{chiRouter}
lib.RegisterRoutes(adapter, "/api/v1", library)
```

Handlers remain unchanged when swapping frameworks - they only depend on service interfaces.

## Adding a new dependency
1. Define an interface (consumer side) in the appropriate layer
2. Implement in infra layer
3. Provide the implementation in the correct module function
4. Inject via constructor in the dependent component
