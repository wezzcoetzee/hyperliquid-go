package exchange

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
)

// ApproveAgent authorizes an agent wallet to sign L1 actions on the user's behalf.
func (c *Client) ApproveAgent(ctx context.Context, agentAddress, agentName string) error {
	n := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"agentAddress":     agentAddress,
		"agentName":        agentName,
		"nonce":            n,
	}
	action := c.userSignWrapper("approveAgent", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:ApproveAgent", userSignActionFields.ApproveAgent, msg, nil)
}

// ApproveBuilderFee authorizes a builder address to charge up to maxFeeRate.
func (c *Client) ApproveBuilderFee(ctx context.Context, builder, maxFeeRate string) error {
	n := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"maxFeeRate":       maxFeeRate,
		"builder":          builder,
		"nonce":            n,
	}
	action := c.userSignWrapper("approveBuilderFee", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:ApproveBuilderFee", userSignActionFields.ApproveBuilderFee, msg, nil)
}

// CreateSubAccount creates a sub-account with the given display name.
func (c *Client) CreateSubAccount(ctx context.Context, name string) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "createSubAccount")
	a.Set("name", name)
	return c.submitL1(ctx, a, nil)
}

// SubAccountModify renames an existing sub-account.
func (c *Client) SubAccountModify(ctx context.Context, subAccountUser, name string) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "subAccountModify")
	a.Set("subAccountUser", subAccountUser)
	a.Set("name", name)
	return c.submitL1(ctx, a, nil)
}

// SetReferrer attaches a referral code to the caller's account (one-time).
func (c *Client) SetReferrer(ctx context.Context, code string) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "setReferrer")
	a.Set("code", code)
	return c.submitL1(ctx, a, nil)
}

// RegisterReferrer creates a referral code owned by the caller.
func (c *Client) RegisterReferrer(ctx context.Context, code string) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "registerReferrer")
	a.Set("code", code)
	return c.submitL1(ctx, a, nil)
}

// CreateVaultRequest describes a new vault.
type CreateVaultRequest struct {
	Name        string
	Description string
	InitialUsd  uint64
}

// CreateVault creates a vault with the caller as leader.
func (c *Client) CreateVault(ctx context.Context, req CreateVaultRequest) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "createVault")
	a.Set("name", req.Name)
	a.Set("description", req.Description)
	a.Set("initialUsd", req.InitialUsd)
	return c.submitL1(ctx, a, nil)
}

// VaultModifyRequest updates vault settings. nil pointers leave fields unchanged.
type VaultModifyRequest struct {
	VaultAddress          string
	AllowDeposits         *bool
	AlwaysCloseOnWithdraw *bool
}

// VaultModify updates vault settings.
func (c *Client) VaultModify(ctx context.Context, req VaultModifyRequest) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "vaultModify")
	a.Set("vaultAddress", req.VaultAddress)
	if req.AllowDeposits != nil {
		a.Set("allowDeposits", *req.AllowDeposits)
	}
	if req.AlwaysCloseOnWithdraw != nil {
		a.Set("alwaysCloseOnWithdraw", *req.AlwaysCloseOnWithdraw)
	}
	return c.submitL1(ctx, a, nil)
}

// VaultDistribute distributes vault profits to followers.
func (c *Client) VaultDistribute(ctx context.Context, vaultAddress string, usd uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "vaultDistribute")
	a.Set("vaultAddress", vaultAddress)
	a.Set("usd", usd)
	return c.submitL1(ctx, a, nil)
}
