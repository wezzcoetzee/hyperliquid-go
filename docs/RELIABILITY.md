# Reliability

This is a library, not a service — "reliability" here means the contracts the SDK upholds for callers under network failure, server errors, and reconnects.

## Error handling strategy

Every public error is typed and `errors.As` / `errors.Is` friendly. Callers should match on type, not on string content.

| Error | Source | Caller action |
|---|---|---|
| `*hyperliquid.APIError` | Non-2xx HTTP response | Inspect `Status` + `Body`; retry transport-level errors with backoff. **Do not** retry an Exchange write with the same nonce. |
| `*exchange.ActionRejected` | HTTP 200 with `status:"err"` body | Action validation failure (e.g., insufficient margin). Inspect reason; **do not** blindly retry. |
| `exchange.ErrNoSigner` | Write attempted with `Config.Signer == nil` | Programmer error; wire a signer. |
| `*hyperliquid.SignError` | Wrap of signer-layer errors | Bubble up; signer-specific recovery. |
| `*hyperliquid.WSError` | WebSocket protocol error | Caller may unsubscribe + resubscribe; auto-reconnect handles transient cases. |

`transport.TransportAPIError` is internal and is wrapped to `*hyperliquid.APIError` at the package boundary by `wrappingHTTP` in `client.go`. Users should never see the internal type.

## WebSocket reliability

- **Single multiplexed connection.** All subscriptions on a `Client` share one WS connection (`ws/registry.go`).
- **Auto-reconnect with exponential backoff.** Disconnects trigger a reconnect attempt; backoff grows exponentially, bounded by an upper cap.
- **Subscription replay on reconnect.** All active subscriptions are re-sent automatically. Handlers may receive a brief duplicate or gap window across the reconnect; handlers must be idempotent.
- **No buffering of dropped messages.** If a message is missed during reconnect, it is gone. Callers needing strict ordering should reconcile via REST after reconnect.

## Concurrency contract

- `hyperliquid.Client` is safe for concurrent use across goroutines.
- `info.Client`, `exchange.Client`, and `ws.Client` are individually safe.
- `signer.Signer` implementations **must** be safe for concurrent `SignTypedData` calls (the SDK may invoke from multiple goroutines).
- The default `signer/privkey` implementation is concurrent-safe.

## Nonce safety

- `exchange/nonce.go` produces monotonic nonces (millisecond timestamp + collision resolution).
- Sharing a `Client` across goroutines is safe.
- Sharing a key across **separate** `Client` instances or processes risks nonce collisions; coordinate externally if you must.

## Logging

The SDK does not log. Callers control all logging via their `*http.Client` or by wrapping handlers. The SDK deliberately avoids any built-in logger to stay dependency-light and surprise-free.

What callers should **not** log:
- Raw `Signer` instances
- Full signed-action bodies (contain signatures bound to nonces)
- Private keys (obviously)

## Recovery procedures

| Failure mode | Detection | Mitigation |
|---|---|---|
| Network blip during HTTP write | `*APIError` with 5xx or context error | Retry with backoff; **do not** reuse nonce — generate a fresh action. |
| WS disconnect | Auto-reconnect kicks in; handlers pause until reconnect | None required. Optionally reconcile via REST for strict ordering. |
| Signer error (KMS down, etc.) | `*SignError` | Caller-specific; circuit-break or fail upstream. |
| Action rejected (margin, etc.) | `*ActionRejected` | Inspect reason; surface to user, do not auto-retry. |

## Deployment safety

N/A — library. Callers control deployment of their own services.

## Cross-references

- Public error types: [../errors.go](../errors.go), [../exchange/responses.go](../exchange/responses.go)
- WS internals: [../ws/registry.go](../ws/registry.go), [../ws/conn.go](../ws/conn.go)
- Security context for nonces & signing: [SECURITY.md](SECURITY.md)
