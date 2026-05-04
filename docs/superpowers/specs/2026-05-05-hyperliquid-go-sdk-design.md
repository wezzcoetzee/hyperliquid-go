# Hyperliquid Go SDK — Design

**Date:** 2026-05-05
**Module:** `github.com/wezzcoetzee/hyperliquid`
**Go version:** 1.22+
**Goal:** Full-parity Go port of [`nktkas/hyperliquid`](https://github.com/nktkas/hyperliquid) (TypeScript SDK), usable as a library by Go applications.

## Scope

Full parity with the TS SDK: Info, Exchange, Subscriptions, multi-sig, vaults, sub-accounts, validator/staking actions, mainnet + testnet.

## Architecture Overview

Single entry point producing a `Client` with three method groups:

```go
client, _ := hyperliquid.New(hyperliquid.Config{
    Network: hyperliquid.Mainnet, // or hyperliquid.Testnet
    Signer:  mySigner,            // optional; required for Exchange writes
    HTTP:    nil,                 // optional custom *http.Client
})

mids, _ := client.Info.AllMids(ctx)
resp, _ := client.Exchange.Order(ctx, exchange.OrderParams{...})
sub, _  := client.Subscriptions.Trades(ctx, "BTC", handler)
```

- `client.Info` — read-only REST (`/info`)
- `client.Exchange` — signed actions (`/exchange`), requires `Signer`
- `client.Subscriptions` — WebSocket subscriptions

Transports (`HTTP`, `WebSocket`) are pluggable interfaces for testing/proxying.

## Package Layout

```
hyperliquid/
├── go.mod
├── client.go                  // hyperliquid.Client, Config, New()
├── network.go                 // Mainnet/Testnet endpoints
├── errors.go                  // typed errors
│
├── signer/
│   ├── signer.go              // Signer interface
│   └── privkey/privkey.go     // default impl (go-ethereum/crypto)
│
├── types/                     // shared domain types (Asset, Side, Tif, Cloid, etc.)
│
├── info/
│   ├── info.go                // Info methods
│   └── types.go
│
├── exchange/
│   ├── exchange.go            // Order, Cancel, Modify, Transfer, Withdraw, …
│   ├── actions.go             // action structs + EIP-712 hashing
│   ├── sign.go                // signing pipeline
│   └── types.go
│
├── ws/
│   ├── ws.go                  // SubscriptionClient
│   ├── conn.go                // reconnect/heartbeat
│   └── subs.go                // typed subscription helpers
│
├── transport/
│   ├── http.go                // HTTP interface + default impl
│   └── ws.go                  // WS interface + default impl (gorilla/websocket)
│
└── internal/
    └── eip712/                // EIP-712 encoding (no external dep)
```

### Dependencies (minimal)
- `github.com/ethereum/go-ethereum` — only in `signer/privkey` (pluggable away)
- `github.com/gorilla/websocket` — only in `transport`
- `golang.org/x/sync` — errgroup for subscriptions
- stdlib otherwise

## Signer Interface

```go
type Signer interface {
    Address() common.Address
    SignTypedData(ctx context.Context, domain TypedDataDomain, types TypedDataTypes, message map[string]any) ([]byte, error)
}
```

Default: `signer/privkey.New(hexKey)` using go-ethereum's `crypto.Sign`. Users implement the interface for KMS, Ledger, browser wallet, Fireblocks, etc.

### Two EIP-712 Domains
- **L1 actions** (orders, cancels): domain `Exchange`, chainId `1337`. Signs an `Agent{source, connectionId}` over the action hash.
- **User-signed actions** (transfers, withdrawals, approveAgent): domain `HyperliquidSignTransaction`, chainId `42161` (mainnet) / `421614` (testnet).

### Signing Pipeline
1. Serialize action to msgpack with deterministic field ordering matching TS.
2. Action hash: `keccak256(msgpack(action) || nonce(8 bytes BE) || vaultAddress(0x00 or 21 bytes))`.
3. Build EIP-712 typed-data (`Agent` for L1, action-specific for user-signed).
4. `signer.SignTypedData(...)` → `r`, `s`, `v`.
5. POST `{action, nonce, signature, vaultAddress?}` to `/exchange`.

Msgpack encoding parity with TS is locked by golden-file fixtures.

## API Coverage

### `client.Info`
- **Market:** `AllMids`, `L2Book`, `CandleSnapshot`, `Meta`, `MetaAndAssetCtxs`, `SpotMeta`, `SpotMetaAndAssetCtxs`, `FundingHistory`, `PredictedFundings`
- **User:** `ClearinghouseState`, `SpotClearinghouseState`, `OpenOrders`, `FrontendOpenOrders`, `UserFills`, `UserFillsByTime`, `UserFunding`, `UserNonFundingLedgerUpdates`, `UserRateLimit`, `OrderStatus`, `HistoricalOrders`, `TwapHistory`, `SubAccounts`, `Referral`, `PortfolioPeriods`
- **Vault:** `VaultDetails`, `VaultSummaries`, `UserVaultEquities`
- **Validator/staking:** `Delegations`, `DelegatorSummary`, `DelegatorHistory`, `DelegatorRewards`, `ValidatorSummaries`
- **Misc:** `MaxBuilderFee`, `LegalCheck`, `ExtraAgents`

### `client.Exchange`
- **Trading:** `Order`, `Cancel`, `CancelByCloid`, `BatchModify`, `Modify`, `ScheduleCancel`, `UpdateLeverage`, `UpdateIsolatedMargin`, `TwapOrder`, `TwapCancel`
- **Transfers:** `UsdSend`, `SpotSend`, `Withdraw3`, `UsdClassTransfer`, `SubAccountTransfer`, `SubAccountSpotTransfer`, `VaultTransfer`, `TokenDelegate`
- **Account:** `ApproveAgent`, `ApproveBuilderFee`, `CreateSubAccount`, `SubAccountModify`, `SetReferrer`, `RegisterReferrer`, `CreateVault`, `VaultModify`, `VaultDistribute`
- **Validator:** `CDeposit`, `CWithdraw`, `CSignerAction`, `CValidatorAction`
- **Multi-sig:** `MultiSig` — wraps an inner action with N signatures.

### `client.Subscriptions`
`AllMids`, `Notification`, `WebData2`, `Candle`, `L2Book`, `Trades`, `OrderUpdates`, `UserEvents`, `UserFills`, `UserFundings`, `UserNonFundingLedgerUpdates`, `ActiveAssetCtx`, `ActiveAssetData`, `Bbo`, `TwapSliceFills`, `TwapHistory`.

POST-over-WS for Info/Exchange opt-in via `Config.WebSocketPosts: true`.

Each subscription returns a `Subscription` handle with `Unsubscribe()`. Delivery is via callback (`func(Event)`), composes naturally with `errgroup` and `ctx` cancellation.

## Errors

```go
type APIError    struct{ Status int; Body string; Type string }   // non-2xx HTTP
type ActionError struct{ Action string; Response string }          // /exchange status:"err"
type SignError   struct{ Err error }
type WSError     struct{ Code int; Reason string }
```

All wrap-compatible with `errors.Is/As`. Network errors pass through unwrapped.

## Context & Concurrency

- Every public method takes `ctx context.Context` as first arg.
- `Client` is safe for concurrent use.
- HTTP: shared `*http.Client`, 30s default timeout (configurable).
- WS: single connection multiplexed across subscriptions; write-mutex; read goroutine fans out to subscribers; automatic reconnect with exponential backoff; subscriptions auto-resubscribe on reconnect.
- Nonces: monotonic millisecond timestamps, per-`Client` mutex bumps on collision (matches TS).

## Testing

- **Unit:** `httptest.Server` for HTTP; in-memory WS server for ws transport.
- **Signing fixtures:** `testdata/signing/*.json` golden files capturing `(action, nonce, vault) → expected_hash, expected_signature` from the TS SDK using a fixed test key.
- **Msgpack fixtures:** golden bytes per action shape, regenerated by `cmd/genfixtures` (Go) plus a TS reference script under `testdata/ts-ref/`.
- **Integration:** `//go:build integration` tag, gated by `HL_TESTNET_KEY` env. Covers order placement, cancel, info reads, one WS subscription against testnet.
- **CI:** unit tests on every push; integration nightly using a dedicated low-balance testnet account.

## Docs

- `examples/`: `info-only`, `basic-trading`, `ws-subscriptions`, `custom-signer-kms`.
- Godoc on all exported symbols.
- README with quickstart and links to the upstream TS docs.
