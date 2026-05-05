# Architecture

## System overview

`hyperliquid-go` is a Go SDK for the [Hyperliquid](https://hyperliquid.gitbook.io/hyperliquid-docs/) exchange. It is a feature-complete port of the canonical TypeScript SDK [`nktkas/hyperliquid`](https://github.com/nktkas/hyperliquid), with action signing producing byte-identical bytes verified against captured TS fixtures. It exposes three method groups — `Info` (read-only REST), `Exchange` (signed write actions), and `Subscriptions` (multiplexed WebSocket) — through a single top-level `hyperliquid.Client`.

## Tech stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| Language | Go 1.24+ | Library runtime |
| HTTP | `net/http` | REST transport |
| WebSocket | `github.com/gorilla/websocket` v1.4.2 | WS transport |
| Crypto / EIP-712 | `github.com/ethereum/go-ethereum` v1.17.2 | secp256k1 signing, keccak |
| Crypto (extras) | `golang.org/x/crypto` | Auxiliary primitives |
| Action serialization | Internal `internal/msgpack` | Deterministic msgpack matching TS `@msgpack/msgpack` output |
| Typed-data hashing | Internal `internal/eip712` | EIP-712 domain/struct hashing |

No database. No frontend. No deployable services — this is a library consumed by user code.

## High-level diagram

```
                     hyperliquid.Client
                     /       |        \
                Info     Exchange   Subscriptions
                  |          |            |
            transport.HTTP   |       transport.WS (gorilla)
                  └────┬─────┘            │
                       │                  │
                  HTTP POST              WSS
                       │                  │
                       └──── Hyperliquid API ────┘
                              (mainnet / testnet)

Exchange writes:
  user code → Exchange.<Action> → exchange/sign.go → signer.Signer.SignTypedData
            → msgpack-encoded action + signature → POST /exchange
```

## Directory structure

```
├── client.go, network.go, errors.go    # Top-level package: Client, Config, Network, public errors
├── info/                                # Read-only /info endpoints (33 methods)
│   ├── market.go, user.go, vault.go, staking.go, misc.go
│   └── types.go                         # Response types
├── exchange/                            # Signed /exchange actions (~30 methods)
│   ├── exchange.go                      # Client wiring
│   ├── orders.go, transfers.go, account.go, validator.go, multisig.go
│   ├── sign.go, sign_parity_test.go     # Action hashing + fixture parity tests
│   ├── nonce.go                         # Monotonic nonce generator
│   ├── responses.go                     # Action response types + ActionRejected
│   └── testdata/fixtures/               # TS-generated signing fixtures (source of truth)
├── ws/                                  # WebSocket subscriptions + POST-over-WS primitive
│   ├── conn.go, registry.go             # Connection mgmt, subscription multiplexer, reconnect
│   ├── subs.go                          # 16 typed subscription helpers
│   ├── post.go                          # Low-level WS-POST primitive (Info routing TBD)
│   └── types.go
├── signer/                              # Signing interface + default impl
│   ├── signer.go                        # Signer interface, Domain/Types/Signature
│   └── privkey/                         # secp256k1 implementation via go-ethereum
├── transport/                           # HTTP + WS transport adapters
│   ├── http.go                          # Default HTTP client, TransportAPIError
│   └── ws.go                            # WS dialer abstraction
├── types/                               # Cross-package shared types
├── internal/
│   ├── eip712/                          # EIP-712 typed-data hashing
│   └── msgpack/                         # Deterministic msgpack encoder (TS-parity)
├── examples/                            # Runnable example programs
│   ├── basic-trading/, custom-signer-kms/, info-only/, ws-subscriptions/
└── tools/ts-fixtures/                   # TS script that regenerates signing fixtures
```

## Key architectural decisions

### Byte-identical TS-SDK parity for signing
- **Context**: Hyperliquid action signatures are EIP-712 over a deterministic msgpack-serialized action. The smallest deviation in field ordering, integer encoding, or string normalization changes the digest and rejects the action.
- **Decision**: Treat the canonical TS SDK as the wire-format source of truth. Capture fixtures from TS via `tools/ts-fixtures` and assert byte equality in `exchange/sign_parity_test.go`.
- **Alternatives considered**: Hand-derive from API docs (rejected — docs lag); fuzz against testnet (rejected — not reproducible).
- **Consequences**: Any change to msgpack encoding, action structs, or signing helpers requires fixture regeneration. Encoding bugs surface in CI, not in production.

### Signer is an interface, not a concrete type
- **Context**: Users want to sign with private keys, hardware wallets, KMS, or remote services.
- **Decision**: `signer.Signer` is a two-method interface (`Address`, `SignTypedData`). The default `signer/privkey` is one of N possible backends.
- **Alternatives considered**: Concrete struct holding a private key (rejected — couples SDK to secret material).
- **Consequences**: SDK never sees raw key bytes when used with a remote signer. See `examples/custom-signer-kms/`.

### Single multiplexed WebSocket per Client
- **Context**: Most clients want many subscriptions; opening N connections is wasteful and harder to reconnect.
- **Decision**: One WS connection per `Client`, multiplexed via `ws/registry.go`. On disconnect, exponential backoff reconnect; on reconnect, all active subscriptions are replayed automatically.
- **Consequences**: Subscription handlers must be idempotent across reconnects. State held in user closures is the user's responsibility.

### Public errors are typed; transport errors are unwrapped at the boundary
- `transport.TransportAPIError` is internal; `client.go` wraps it into the public `*hyperliquid.APIError` so users only see the public surface. See `wrappingHTTP` in `client.go`.

## Data flow

**Signed write (e.g., placing an order):**
1. User calls `c.Exchange.Order(ctx, req)`.
2. `exchange/orders.go` builds the action struct.
3. `exchange/sign.go` msgpack-encodes the action (deterministic), hashes via EIP-712 with the network-specific domain (mainnet `signature_chain_id=42161`, testnet `421614`), calls `Signer.SignTypedData`.
4. Action + signature + nonce posted to `/exchange` via `transport.HTTP`.
5. HTTP 200 with `status:"err"` body → `*exchange.ActionRejected`. Non-2xx → `*hyperliquid.APIError`.

**Read (e.g., AllMids):**
1. User calls `c.Info.AllMids(ctx)`.
2. `info/market.go` POSTs `{"type":"allMids"}` to `/info` via `transport.HTTP`.
3. JSON unmarshal into typed response.

**WebSocket subscription:**
1. User calls `c.Subscriptions.Trades(ctx, "BTC", handler)`.
2. `ws/registry.go` ensures the connection is dialed (lazy), registers a typed handler, sends a `subscribe` frame.
3. Inbound frames are demultiplexed by channel + key, decoded into the typed payload, dispatched to the user's handler in a goroutine.
4. On disconnect: exponential backoff reconnect, replay all active subscription frames.
5. `Unsubscribe(ctx)` sends an `unsubscribe` frame and removes the registry entry.

## Environment and deployment

This is a library — no environments to deploy. CI runs `go test ./...` against Go 1.24+. The TS fixture regeneration step requires Node.js and `npx tsx`; see `tools/ts-fixtures/`.

Runtime configuration is via `hyperliquid.Config`:
- `Network` — `Mainnet` (default) or `Testnet`
- `Signer` — required for `Exchange` writes
- `HTTP` — optional `*http.Client` override
- `BaseURL`, `WSURL` — overrides for proxies/tests

## Cross-references

- Detailed design decisions: [docs/design-docs/](docs/design-docs/)
- Security model and signing: [docs/SECURITY.md](docs/SECURITY.md)
- Reliability (reconnect, errors): [docs/RELIABILITY.md](docs/RELIABILITY.md)
- Quality bar: [docs/QUALITY_SCORE.md](docs/QUALITY_SCORE.md)
- Roadmap: [docs/PLANS.md](docs/PLANS.md)
