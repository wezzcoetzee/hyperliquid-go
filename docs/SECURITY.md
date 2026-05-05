# Security

This is a client SDK; "security" here means safe handling of signing material and correctness of cryptographic operations. There is no server-side auth model.

## Signing model

Hyperliquid authenticates write actions via EIP-712 typed-data signatures over a deterministic msgpack-serialized action body, using secp256k1 (Ethereum-compatible). Two domains exist:

- **L1 actions** (orders, cancels, transfers within Hyperliquid): domain uses `signature_chain_id` matching the network — `42161` (Arbitrum One) for mainnet, `421614` (Arbitrum Sepolia) for testnet.
- **User-signed actions** (account-level: deposits, withdraws, agent registration, etc.): domain explicitly carries the chain id of the signing chain.

`exchange/sign.go` is the single place this is computed. Any change there must be validated against captured TS-SDK fixtures (`exchange/testdata/fixtures/`).

## Key handling

- The SDK **never** holds key material directly. `Config.Signer` is a `signer.Signer` interface; how that signer obtains keys is its concern.
- The default `signer/privkey` implementation accepts a `0x...` hex string and uses `go-ethereum/crypto`. The private key lives in the `*PrivKeySigner` for the lifetime of the process.
- For production, prefer a remote/HSM/KMS-backed `Signer`. See `examples/custom-signer-kms/` for the integration shape.
- Do **not** log `Signer` instances or signed-action payloads in user code — signatures are tied to nonces and replays may be possible if a nonce window is reused.

## Nonces

`exchange/nonce.go` provides a monotonic nonce source (millisecond timestamp with collision resolution). Two clients sharing the same key but different nonce sources can collide; share the `Client` (or coordinate nonces externally) when running multiple instances against the same key.

## Multi-sig

`exchange/multisig.go` implements the multi-sig signing flow. Status: **unit-test coverage only — not yet validated against testnet end-to-end** (see [PLANS.md](PLANS.md)). Verify on testnet before relying on it for production funds.

## Transport

- All HTTP traffic uses `https://` (mainnet/testnet endpoints in `network.go`); WS uses `wss://`.
- `Config.HTTP` lets users install proxies, custom TLS configs, or rate-limited transports.
- No request signing beyond the EIP-712 action signature — there is no API key model.

## Action validation

- HTTP 200 with `status:"err"` body → `*exchange.ActionRejected`. Treat these the same as authentication failures: do **not** auto-retry without inspecting the reason (insufficient margin, etc.).
- HTTP non-2xx → `*hyperliquid.APIError`. Safe to retry transport-level errors; **never** blindly retry an action with the same nonce.

## Threat model (caller-side)

Risks the SDK does not protect against (caller's responsibility):
- Compromised host with access to the signer process — out of scope; mitigated by remote signers.
- Replay of captured signed actions — Hyperliquid enforces nonce windows server-side; SDK relies on this.
- Hostile `*http.Client` injected via `Config.HTTP` — caller is trusted.

## Security checklist for new features

When adding a new signed action:
- [ ] Action struct field tags match TS SDK exactly (msgpack determinism)
- [ ] `signature_chain_id` correct for L1 vs user-signed
- [ ] Fixture captured via `tools/ts-fixtures` and parity-tested in `exchange/sign_parity_test.go`
- [ ] Both mainnet and testnet fixture variants where the action differs by network
- [ ] No key material in struct `String()` / log output
- [ ] Nonce source threaded through correctly (`exchange/nonce.go`)
- [ ] Godoc on every exported symbol

## Cross-references

- Architecture & data flow: [../ARCHITECTURE.md](../ARCHITECTURE.md)
- Reliability / error model: [RELIABILITY.md](RELIABILITY.md)
- Core beliefs (parity invariant): [design-docs/core-beliefs.md](design-docs/core-beliefs.md)
