package exchange

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
)

// CDeposit deposits HYPE from the caller's spot wallet into the staking pool.
// wei is the raw token amount.
//
// Note: validator-only actions (CSignerAction, CValidatorAction) are
// intentionally not implemented. Their schemas are validator-internal and
// likely to change; most users will never need them.
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
