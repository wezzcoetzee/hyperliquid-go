# Core Beliefs

Foundational technical beliefs that guide every decision in this SDK. These rarely change. Anything contradicting them needs an explicit superseding design doc.

## 1. Byte-for-byte parity with the canonical TS SDK is the contract

**Belief**: Action signatures produced by this Go SDK must be byte-identical to those produced by `nktkas/hyperliquid` for the same input.

**Rationale**: Hyperliquid action signing is EIP-712 over deterministic msgpack. A single byte of drift — field order, integer width, string normalization — yields a different digest and a rejected action. The TS SDK is the upstream source of truth used by the protocol team. Anything else is guesswork.

**Implications**:
- Captured TS fixtures live in `exchange/testdata/fixtures/` and are the test oracle.
- Any change touching `exchange/sign.go`, action structs, `internal/msgpack/`, or `internal/eip712/` requires regenerating fixtures via `tools/ts-fixtures` and re-running `exchange/sign_parity_test.go`.
- "Cleaning up" msgpack output, struct ordering, or hash construction without fixture re-validation is a critical defect.

## 2. The SDK never holds keys; signers are pluggable

**Belief**: Key material is the caller's concern. The SDK consumes a `signer.Signer` interface and never sees raw keys unless the caller explicitly uses `signer/privkey`.

**Rationale**: Production trading systems use HSMs, KMS, hardware wallets, or remote signers. Coupling the SDK to in-process keys would force every user to either accept that coupling or fork.

**Implications**:
- New signing surfaces extend the `Signer` interface only with strong justification — every method is a burden on every backend.
- The default `signer/privkey` is one backend among many; it is not special-cased anywhere in the SDK.
- KMS, hardware, and remote-signer backends are out-of-tree (see `examples/custom-signer-kms/`) to keep the dependency surface minimal.

## 3. One multiplexed WebSocket per Client; reconnects are transparent

**Belief**: Callers should never need to manage WS lifecycle. Subscribe, get messages, unsubscribe.

**Rationale**: Multiple WS connections per client wastes server resources, complicates reconnect semantics, and forces every caller to re-implement subscription replay.

**Implications**:
- `ws/registry.go` owns the single connection and demultiplexes by channel + key.
- Reconnect uses exponential backoff and replays all active subscriptions automatically.
- Handlers must be idempotent across reconnects — duplicate-message windows are inherent to replay.
- Callers who need strict ordering reconcile via REST after a disconnect.

## 4. Small, justified dependency surface

**Belief**: Every dependency is a maintenance burden, a CVE surface, and a transitive constraint on users.

**Rationale**: A trading SDK lives in production trading systems. Each dep we add becomes a dep every consumer must accept.

**Implications**:
- Current deps: `go-ethereum` (signing), `gorilla/websocket` (WS), `x/crypto` (auxiliary). Adding to this list requires CHANGELOG + design-doc justification.
- Prefer stdlib (`net/http`, `encoding/json`) over wrappers.
- No logging library, no metrics library, no retry library — callers compose those themselves via `Config.HTTP`.

## 5. Errors are typed; callers match on type, not strings

**Belief**: Every public error path returns a concrete type that works with `errors.As` / `errors.Is`. String comparison is never the right way to handle SDK errors.

**Rationale**: String error matching is brittle and breaks silently when wording changes. Type matching is checked by the compiler.

**Implications**:
- Public error types live in `errors.go` and `exchange/responses.go`.
- Internal transport errors (`transport.TransportAPIError`) are wrapped at the package boundary into public types — see `wrappingHTTP` in `client.go`.
- Adding a new public error type is a CHANGELOG entry.

## 6. Faithful port over Go-flavored "improvement"

**Belief**: Where the TS SDK shape is awkward in Go, mirror it anyway and document the awkwardness.

**Rationale**: Users translate TS examples into Go. Surprises cost more than ergonomics save. The TS SDK is also where new Hyperliquid features land first, so divergence makes upstream tracking harder.

**Implications**:
- Method names track the TS SDK 1:1.
- Field names in request/response structs match the JSON wire format.
- "Helpful" wrappers (auto-retry, decimal helpers, etc.) belong in user code, not the SDK.
