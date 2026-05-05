# Product Sense

## What this product is

`hyperliquid-go` is a Go SDK for the Hyperliquid perpetual exchange. The audience is Go developers building trading bots, market-making infrastructure, monitoring tools, and on-chain integrations who want a typed, idiomatic Go interface to Hyperliquid without translating the JS/TS SDK themselves. It exists because Hyperliquid's canonical SDK is TypeScript, and Go shops were either re-implementing signing (and getting it subtly wrong) or shelling out to Node.

## User mental models

Users think of this as a **thin, faithful port** of the canonical TS SDK — not a reimagining. They expect:
- Method names and shapes that map 1:1 to the TS SDK so they can read TS examples and translate.
- Action signing that "just works" — they should never have to think about EIP-712 domains, msgpack ordering, or signature chain IDs.
- Errors that are inspectable with `errors.As` / `errors.Is`, not opaque strings.
- A `Client` they can share across goroutines.

They do **not** think of this as a trading framework. There is no order management, no risk engine, no strategy abstraction — it's a transport + signing layer.

## Core user flows

1. **Read market data** — `c.Info.AllMids(ctx)` / `c.Info.L2Book(ctx, "BTC")`. Success: typed response, no signer needed. Common failure: network errors, rate limits surfacing as `*APIError`.
2. **Place / cancel an order** — Build `OrderRequest`, call `c.Exchange.Order(ctx, req)`. Success: typed response with `oid`. Common failure: `*ActionRejected` for validation (insufficient margin, price out of bounds), `*APIError` for transport.
3. **Subscribe to a stream** — `c.Subscriptions.Trades(ctx, "BTC", handler)`. Success: handler invoked per message; auto-reconnect transparent. Common failure: handler panic crashes the goroutine — handler responsibility.
4. **Plug a custom signer** — Implement `signer.Signer` for KMS / hardware wallet, pass into `Config.Signer`. The SDK never sees raw keys.

## Product principles

1. **Parity over ergonomics.** When TS SDK shape is awkward in Go, mirror it anyway. Users translating examples should not be surprised. Document the awkwardness; don't paper over it.
2. **Fail loud, fail typed.** Every error path returns a typed error users can match on. No string-comparison error handling.
3. **Zero surprises in concurrency.** `Client` is goroutine-safe. The WS connection is shared and multiplexed — users should never need to dial directly.
4. **Small dependency surface.** New deps require strong justification. Today: go-ethereum (signing), gorilla/websocket (WS), x/crypto.
5. **Bytes are the contract.** Signing byte-equality with the TS SDK is the most important property of this library. Test it relentlessly.

## What this product is NOT

- **Not a trading bot framework.** No strategies, no order books, no PnL tracking, no backtester.
- **Not a wallet.** Key custody is the user's problem; the SDK only consumes a `Signer`.
- **Not a Hyperliquid Node / on-chain client.** This talks to the Hyperliquid REST/WS API only.
- **Not opinionated about retries or rate limiting.** Users plug their own `*http.Client` for that.
- **Not a re-export of every TS helper.** Helpers that are trivial in Go (e.g., decimal formatting one-liners) are deliberately omitted.

## Cross-references

- Product specs: [product-specs/](product-specs/)
- Roadmap: [PLANS.md](PLANS.md)
- Design beliefs: [design-docs/core-beliefs.md](design-docs/core-beliefs.md)
