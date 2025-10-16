# Implementation Plan: DI cleanup with uber-go/dig

**Branch**: `003-previously-after-implementing` | **Date**: 2025-10-16 | **Spec**: `/specs/003-previously-after-implementing/spec.md`
**Input**: Feature specification from `/specs/003-previously-after-implementing/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Unify dependency injection during route initialization by introducing a single container builder, layered DI modules, and framework-agnostic route registration using uber-go/dig.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.20  
**Primary Dependencies**: dig (uber-go/dig), gorilla/mux (existing), go-sql-driver/mysql (existing), go-redis (existing)  
**Storage**: MySQL, Redis (as-is)  
**Testing**: go test, testify; container tests with replaced providers  
**Target Platform**: Linux server; embeddable library  
**Project Type**: Single backend service/library  
**Performance Goals**: Container build < 200ms locally; zero runtime DI allocations on hot paths  
**Constraints**: Handlers depend on service interfaces only; no repositories in delivery  
**Scale/Scope**: Whole app DI wiring and route registration

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- Interfaces declared at consumer side – PASS
- Clean Architecture layering – PASS
- Test-first and coverage for new modules – NEEDS CLARIFICATION (add tests in tasks)
- Framework independence of business logic – PASS

## Project Structure

### Documentation (this feature)

```
specs/003-previously-after-implementing/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```
lib/
├── container/
│   ├── builder.go          # dig.Container builder and modules registration
│   ├── modules.go          # Provide functions per layer (config/infra/repos/services/delivery)
│   └── testing.go          # Test helpers to replace providers
├── router.go               # Framework-agnostic route registrar interfaces/adapters
└── config.go               # Library-level config contract
```

**Structure Decision**: Introduce `lib/container` for DI builder and modules; keep domain/use cases unchanged; adapt `app/` to use container and route registrar.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
