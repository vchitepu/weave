# Weave Full Rename Design
**Date:** 2026-04-07  
**Status:** Proposed  
**Type:** Naming refactor (no behavior changes)

## Goal

Rename the project identity from `shine` to `weave` everywhere in the repository and runtime surface, including:

- CLI command and entrypoint path
- Go module path and import references
- User-facing text in help/errors/tests/docs
- Existing superpowers design/plan docs that still reference the old name

This change must preserve all existing rendering features and behavior.

## Scope

### In scope

- Move command entrypoint directory from `cmd/shine` to `cmd/weave`
- Update module path in `go.mod` from `github.com/vinaychitepu/shine` to `github.com/vinaychitepu/weave`
- Update repository-local imports that reference the old module path
- Update command strings and user-visible labels from `shine` to `weave`
- Update project docs/specs/plans that are intended to describe the current tool name
- Build/install binary as `weave`

### Out of scope

- Functional renderer changes
- Theme/layout/wrapping behavior changes
- Backward-compatibility alias command (`shine`) or migration script
- Repository rename on remote hosting (outside local code changes)

## Affected Areas

- `go.mod`
- `cmd/shine/` -> `cmd/weave/`
- `cmd/weave/main.go`
- `cmd/weave/main_test.go`
- `internal/**` files that include module import path or command-name assertions
- `docs/superpowers/specs/*.md` and `docs/superpowers/plans/*.md` where project naming is part of current documentation

## Design Decisions

1. **Hard rename only:** keep one canonical name (`weave`) to avoid long-term dual naming.
2. **No behavior edits:** only naming and paths change; tests validate no regressions.
3. **Minimal churn in wording:** replace name references while preserving existing documentation intent and chronology.

## Implementation Strategy

1. Rename/move CLI entry directory to `cmd/weave`.
2. Update module path in `go.mod`.
3. Update imports and literal command/help strings.
4. Update tests that reference `shine`.
5. Update docs/spec/plan files containing stale project name references.
6. Run full test suite.
7. Build/install `weave` binary and smoke-check CLI help/version.

## Testing Strategy

- `go test ./...` must pass.
- `go build -o ~/Documents/Code/bin/weave ./cmd/weave` must succeed.
- Manual smoke checks:
  - `weave --help`
  - `weave --version`

## Risks and Mitigations

- **Risk:** missed string/path references to `shine`.
  - **Mitigation:** run content search for `shine` and review remaining hits manually.
- **Risk:** broken imports after module rename.
  - **Mitigation:** run full tests and targeted build of CLI entrypoint.

## Success Criteria

- All code compiles and tests pass under `weave` naming.
- CLI invocation is `weave`.
- No unintended behavior changes.
- Remaining `shine` references are only historical (if any) and intentionally retained.
