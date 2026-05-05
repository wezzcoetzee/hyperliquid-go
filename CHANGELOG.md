# Changelog

## Unreleased

- Initial release: full-parity Go port of `nktkas/hyperliquid`.
- `info`: 33 read-only endpoints (market, user, vault, staking, misc).
- `exchange`: ~30 signed actions covering trading, transfers, account, validator, multi-sig.
- `signer.Signer` interface plus `signer/privkey` default secp256k1 implementation.
- `ws.Client`: 16 typed subscription helpers, auto-reconnect with exponential backoff, subscription replay on reconnect.
- Mainnet and testnet support for both HTTP and WebSocket.
- Action signing verified byte-identical to the canonical TS SDK across mainnet, testnet, and user-signed flows.
