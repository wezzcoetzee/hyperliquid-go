# Tech Debt Tracker

## Active debt

| ID | Description | Impact | Effort | Priority | Added |
| -- | ----------- | ------ | ------ | -------- | ----- |
| TD-001 | `Config.WebSocketPosts` does not yet route Info HTTP calls through an open WS connection. `ws.Client.Post` exposes the primitive but the HTTP transport isn't wired through it. | Missed latency optimization for Info-poll-heavy clients with an active subscription. | M | P3 | 2026-05-05 |
| TD-002 | Multi-sig flow has only unit-test coverage — not validated end-to-end against testnet. | Multi-sig users may hit wire-format issues only discoverable in production. | M | P2 | 2026-05-05 |
| TD-003 | `New(cfg Config)` returns `(*Client, error)` but never returns a non-nil error today. Either start validating or change the signature pre-1.0. | API churn risk at 1.0. | S | P3 | 2026-05-05 |

## Resolved debt

None yet.
