# hyperliquid-go

Go SDK for Hyperliquid. Port of [nktkas/hyperliquid](https://github.com/nktkas/hyperliquid).

## Install

    go get github.com/wezzcoetzee/hyperliquid

## Quickstart

```go
client, _ := hyperliquid.New(hyperliquid.Config{Network: hyperliquid.Mainnet})
meta, _ := client.Info.Meta(context.Background())
fmt.Println(meta.Universe)
```

Status: under construction. See `docs/superpowers/plans/` for roadmap.
