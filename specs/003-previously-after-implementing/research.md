# Research: DI cleanup with uber-go/dig

## Decisions

- Decision: Use uber-go/dig as the DI container
  - Rationale: Mature, reflection-based container with Provide/Invoke patterns; easy to swap providers in tests
  - Alternatives: wire (codegen, stricter compile-time), fx (opinionated lifecycle), manual wiring (error-prone)

- Decision: Define modules per layer (config, infra, repositories, services, delivery)
  - Rationale: Matches Clean Architecture; simplifies testing and reuse
  - Alternatives: Single mega-module (hard to test), per-feature modules only (duplication)

- Decision: Route registration via framework-agnostic registrar interface
  - Rationale: Keeps delivery decoupled from specific HTTP framework
  - Alternatives: Direct framework usage in app (couples handlers to framework)

## Unknowns resolved

- Tests strategy for DI: Provide test container with fake providers; run module-level tests using dig Invoke
- Error handling: Fail fast during container build; log detailed construction errors

## Notes

- Avoid handlers depending on repositories; only inject services
- Keep container in `lib/container` and expose minimal API
