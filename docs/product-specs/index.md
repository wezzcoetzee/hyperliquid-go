# Product Specs

This SDK has no user-facing features in the product-spec sense — every "feature" is a port of an upstream Hyperliquid endpoint. The feature inventory is the public API surface itself; treat the godoc on each method group as the spec.

| Surface | Defined in | Notes |
|---|---|---|
| `Info` (33 read-only endpoints) | [`info/`](../../info/) | Market, user, vault, staking, misc |
| `Exchange` (~30 signed actions) | [`exchange/`](../../exchange/) | Trading, transfers, account, validator, multi-sig |
| `Subscriptions` (16 typed WS subs) | [`ws/subs.go`](../../ws/subs.go) | Auto-reconnect with subscription replay |
| `Signer` interface | [`signer/`](../../signer/) | Default: `signer/privkey`. Custom: see `examples/custom-signer-kms/` |

If a future change introduces SDK-level UX that *isn't* a transparent port (e.g., a new helper, an opinionated workflow), add a spec file here describing the user problem and the chosen API shape before implementing.
