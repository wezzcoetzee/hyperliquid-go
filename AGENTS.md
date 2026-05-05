# Agents

## Overview

This is a Go SDK library, not an application — there are no long-running agent workflows here. This file gives any AI agent working on the codebase the orientation it needs to make correct changes without re-deriving project context.

Read [ARCHITECTURE.md](ARCHITECTURE.md) before any non-trivial change. Read [docs/design-docs/core-beliefs.md](docs/design-docs/core-beliefs.md) before changing signing, transport, or public API surface.

## Agent: Coding Agent

### Role
Implement bug fixes, port endpoints from the upstream TS SDK, and extend the public API.

### Responsibilities
- Match the canonical TS SDK ([`nktkas/hyperliquid`](https://github.com/nktkas/hyperliquid)) for any signed action — bytes must be byte-identical to captured fixtures under `exchange/testdata/fixtures/`.
- Run `go test ./...` before claiming completion.
- Add godoc to every new exported symbol.
- Preserve `errors.As` / `errors.Is` compatibility for public error types.

### Boundaries
- **Do not** change action-signing wire formats without regenerating fixtures via `tools/ts-fixtures` and re-confirming parity.
- **Do not** introduce panics in library code paths — return errors.
- **Do not** add dependencies without strong justification; the dependency surface is intentionally small (go-ethereum, gorilla/websocket, x/crypto).
- **Do not** break `Client` concurrent-safety guarantees.
- **Do not** rename or remove exported symbols without a CHANGELOG note and migration path.

### Context docs
- `ARCHITECTURE.md` — package layout, transport/signing flow
- `docs/SECURITY.md` — signing, key handling, multi-sig caveats
- `docs/design-docs/core-beliefs.md` — invariants that should not change
- `docs/RELIABILITY.md` — WS reconnect/backoff contract, error model
- `docs/exec-plans/tech-debt-tracker.md` — known debt to avoid stepping on

### Handoff protocol
- After porting a new endpoint: run `go test ./...`, update `CHANGELOG.md` under `Unreleased`, link any relevant fixture regeneration in the PR.
- If a change affects signing bytes: regenerate TS fixtures (`cd tools/ts-fixtures && npm install && npx tsx index.ts`) and confirm `go test ./exchange/...` passes.

## Agent: Review Agent

### Role
Audit changes for parity, concurrency safety, and public-API stability.

### Responsibilities
- Verify any signing-related diff has a corresponding fixture update or explicit "no wire change" justification.
- Flag new exported symbols missing godoc.
- Confirm tests cover both mainnet (`signature_chain_id=42161`) and testnet (`421614`) paths where applicable.

### Boundaries
- **Do not** rewrite code; comment and hand back to the Coding Agent.

### Context docs
- `docs/QUALITY_SCORE.md` — quality bar
- `docs/design-docs/core-beliefs.md`

## Inter-agent communication

1. Coding Agent implements a change → runs `go test ./...`.
2. Review Agent audits diff against `QUALITY_SCORE.md` and `core-beliefs.md`.
3. If signing bytes changed: Coding Agent regenerates fixtures before merge.
