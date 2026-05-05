# hyperliquid-go

Go SDK for [Hyperliquid](https://hyperliquid.gitbook.io/hyperliquid-docs/). Port of the canonical TypeScript SDK [`nktkas/hyperliquid`](https://github.com/nktkas/hyperliquid). Action signing produces byte-identical bytes to the TS SDK, verified against captured fixtures.

## Status

Pre-1.0. The Info, Exchange, and WebSocket clients are feature-complete with full test coverage. Mainnet and testnet supported. Multi-sig is implemented but only validated by unit tests; verify on testnet before production use.

## Install

```bash
go get github.com/wezzcoetzee/hyperliquid-go
```

Requires Go 1.24+.

## Quickstart

### Read-only

```go
c, _ := hyperliquid.New(hyperliquid.Config{Network: hyperliquid.Mainnet})
mids, _ := c.Info.AllMids(context.Background())
fmt.Println(mids["BTC"])
```

### Trading

```go
signer, _ := privkey.New(os.Getenv("HL_KEY"))
c, _ := hyperliquid.New(hyperliquid.Config{Network: hyperliquid.Mainnet, Signer: signer})

resp, err := c.Exchange.Order(ctx, exchange.OrderRequest{
    Orders: []exchange.OrderParams{{
        Asset: 0, IsBuy: true, LimitPx: "30000", Sz: "0.1",
        OrderType: exchange.OrderType{Limit: &exchange.LimitOrder{Tif: exchange.TifGtc}},
    }},
    Grouping: "na",
})
```

### WebSocket

```go
sub, _ := c.Subscriptions.Trades(ctx, "BTC", func(trades []ws.Trade) {
    for _, t := range trades {
        fmt.Printf("%s %s @ %s\n", t.Side, t.Sz, t.Px)
    }
})
defer sub.Unsubscribe(ctx)
<-ctx.Done()
```

See `examples/` for runnable programs.

## Networks

- `hyperliquid.Mainnet` — `https://api.hyperliquid.xyz`, signature chain id `42161` (Arbitrum One)
- `hyperliquid.Testnet` — `https://api.hyperliquid-testnet.xyz`, signature chain id `421614` (Arbitrum Sepolia)

The zero `Config{}` defaults to mainnet.

## Signing

`signer.Signer` is the swap-in interface for any secp256k1 backend. The default `signer/privkey` implementation uses `go-ethereum/crypto`. To plug in a hardware wallet, KMS, or remote signer, implement the two-method interface:

```go
type Signer interface {
    Address() [20]byte
    SignTypedData(ctx context.Context, domain Domain, types Types, primaryType string, message map[string]any) (Signature, error)
}
```

See `examples/custom-signer-kms/` for a sketch.

## Errors

| Type | Meaning |
|---|---|
| `*hyperliquid.APIError` | Non-2xx HTTP response |
| `*exchange.ActionRejected` | HTTP 200 with `status: "err"` body (action validation failure) |
| `exchange.ErrNoSigner` | Write attempted without a Signer wired into Config |
| `*hyperliquid.SignError` | Wrap around signing-layer errors |
| `*hyperliquid.WSError` | WebSocket protocol error |

All work with `errors.As` / `errors.Is`.

## Concurrency

`hyperliquid.Client` is safe for concurrent use. The HTTP transport is shared across method groups; the WebSocket connection is multiplexed across all active subscriptions and auto-reconnects with exponential backoff.

## Contributing

To regenerate the TS-derived signing fixtures:

```bash
cd tools/ts-fixtures
npm install
npx tsx index.ts
```

This rewrites the JSON files under `exchange/testdata/fixtures/`. Run `go test ./exchange/...` afterward to confirm Go-side parity.

## License

MIT — see [LICENSE](LICENSE).
