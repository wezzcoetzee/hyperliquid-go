# Quality Score

## Quality dimensions

### Code quality
- **Idiomatic Go.** `gofmt` clean. No `interface{}`/`any` outside transport boundaries. Errors returned, never panicked, in library paths.
- **Godoc on every exported symbol.** Enforced manually for now; new exported types/funcs without godoc are a review blocker.
- **Concurrent-safety stated and tested** for any shared type. `Client` and method-group clients are safe for concurrent use.
- **No new dependencies** without a CHANGELOG note and design-doc justification. Current deps: `go-ethereum`, `gorilla/websocket`, `x/crypto`.

### Test coverage
- **Signing parity is non-negotiable.** Every signed action must have a fixture under `exchange/testdata/fixtures/` and a parity test in `exchange/sign_parity_test.go`. Mainnet *and* testnet variants where the signature chain id matters.
- **Unit tests for every Info / Exchange / WS public method.** Look at existing `*_test.go` files in each package as the pattern.
- **Integration tests** (`*_integration_test.go`) gated on env var, not run by default. They hit live testnet.
- **WS reconnect & subscription replay** covered by `ws/conn_test.go`, `ws/registry_test.go`, `ws/live_test.go`.

### Wire-format quality
- **Byte-identical to TS SDK.** Any diff in `exchange/sign.go`, `internal/msgpack/`, or action structs requires fixture regeneration via `tools/ts-fixtures` and re-confirming parity.
- **No silent field reordering.** msgpack is deterministic in this codebase by construction.

### Public API quality
- **Stable typed errors.** `errors.As` and `errors.Is` work for every public error type. Adding a new error type is a CHANGELOG entry.
- **No breaking changes without CHANGELOG.** Pre-1.0 allows breaks; document them.
- **Examples runnable.** Every `examples/<name>/` must compile and (where applicable) run against testnet.

## Review checklist

Before merging any change:

- [ ] `go test ./...` passes
- [ ] `go vet ./...` clean
- [ ] New exported symbols have godoc
- [ ] Signing-relevant changes regenerate fixtures and pass `exchange/sign_parity_test.go`
- [ ] CHANGELOG `Unreleased` updated for user-visible changes
- [ ] No new dependency without justification
- [ ] Tests cover the happy path and at least one error path
- [ ] Concurrent-safety preserved on shared types

## Cross-references

- Reliability targets: [RELIABILITY.md](RELIABILITY.md)
- Security checklist: [SECURITY.md](SECURITY.md)
- Core invariants: [design-docs/core-beliefs.md](design-docs/core-beliefs.md)
