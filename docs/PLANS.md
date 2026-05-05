# Plans & Roadmap

## Current priorities

1. **Pre-1.0 hardening** — Status: in progress. Stabilize the public API surface, lock down error types, ensure CHANGELOG covers every breaking change.
2. **Multi-sig testnet validation** — Status: planned. Multi-sig flow has unit-test coverage but has not been exercised end-to-end against testnet. Required before recommending multi-sig for production use.

## Upcoming

- **POST-over-WebSocket for Info reads.** `ws.Client.Post` exposes the primitive today. `Config.WebSocketPosts` should route Info HTTP calls through an open WS connection when available, saving a TCP+TLS handshake per call. Most useful for clients polling Info endpoints in tight loops alongside an existing subscription. See README "Roadmap" section.
- **Endpoint additions** as Hyperliquid ships new actions/info methods upstream. Track via TS SDK diffs.

## Deferred

- **Trading framework abstractions** — strategies, OMS, PnL. Out of scope; this is a transport+signing SDK, not a trading framework. See [PRODUCT_SENSE.md](PRODUCT_SENSE.md).
- **Built-in retry / rate-limit middleware** — users plug their own `*http.Client`. Adding opinionated retries risks duplicate-nonce hazards on Exchange writes.
- **Expanded `signer/` backends in-tree** (KMS, Ledger, etc.) — kept out-of-tree to avoid pulling in cloud SDKs / USB libs as transitive deps. Sketch lives in `examples/custom-signer-kms/`.

## Cross-references

- Active execution plans: [exec-plans/active/](exec-plans/active/)
- Tech debt: [exec-plans/tech-debt-tracker.md](exec-plans/tech-debt-tracker.md)
- Product context: [PRODUCT_SENSE.md](PRODUCT_SENSE.md)

Last updated: 2026-05-05
