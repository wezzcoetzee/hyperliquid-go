package exchange

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
)

// CDeposit deposits HYPE from the caller's spot wallet into the staking pool.
// wei is the raw token amount.
//
// CSignerAction and CValidatorAction are explicitly out of scope for this SDK.
// They are validator-operator-only actions targeting infrastructure operators,
// not traders. See docs/superpowers/plans/2026-05-05-04-exchange.md for the
// rationale.
func (c *Client) CDeposit(ctx context.Context, wei uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "cDeposit")
	a.Set("wei", wei)
	return c.submitL1(ctx, a, nil)
}

// CWithdraw withdraws unlocked HYPE from staking back to the spot wallet.
func (c *Client) CWithdraw(ctx context.Context, wei uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "cWithdraw")
	a.Set("wei", wei)
	return c.submitL1(ctx, a, nil)
}
