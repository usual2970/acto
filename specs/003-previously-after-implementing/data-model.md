# Data Model: DI Modules

## Entities

### Module
- Fields: name (string), providers (list)
- Relationships: none

### Provider
- Fields: name (string), dependencies (list of names)
- Notes: Resolved by container

### RouteRegistrar
- Fields: methods (register functions)
- Notes: Adapter over HTTP framework
