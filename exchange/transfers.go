package exchange

import (
	"context"
	"strconv"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
)

// hyperliquidChain returns the string the user-signed actions embed in their
// "hyperliquidChain" field. "Mainnet" or "Testnet".
func (c *Client) hyperliquidChain() string {
	if c.SignatureChainID == 421614 {
		return "Testnet"
	}
	return "Mainnet"
}

// signatureChainIDHex returns the chainId formatted as 0x-hex (e.g. "0xa4b1"
// for 42161). The TS SDK puts this in the action body.
func (c *Client) signatureChainIDHex() string {
	return "0x" + strconv.FormatUint(c.SignatureChainID, 16)
}

// userSignActionFields are the typed-field lists per user-signed action.
var userSignActionFields = struct {
	UsdSend           []signer.Field
	SpotSend          []signer.Field
	Withdraw3         []signer.Field
	UsdClassTransfer  []signer.Field
	TokenDelegate     []signer.Field
	ApproveAgent      []signer.Field
	ApproveBuilderFee []signer.Field
}{
	UsdSend: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	},
	SpotSend: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "token", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	},
	Withdraw3: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	},
	UsdClassTransfer: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "toPerp", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	},
	TokenDelegate: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "validator", Type: "address"},
		{Name: "wei", Type: "uint64"},
		{Name: "isUndelegate", Type: "bool"},
		{Name: "nonce", Type: "uint64"},
	},
	ApproveAgent: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "agentAddress", Type: "address"},
		{Name: "agentName", Type: "string"},
		{Name: "nonce", Type: "uint64"},
	},
	ApproveBuilderFee: []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "maxFeeRate", Type: "string"},
		{Name: "builder", Type: "address"},
		{Name: "nonce", Type: "uint64"},
	},
}

func (c *Client) userSignWrapper(actionType string, message map[string]any) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", actionType)
	a.Set("signatureChainId", c.signatureChainIDHex())
	a.Set("hyperliquidChain", message["hyperliquidChain"])
	for _, k := range orderedUserSignKeys(actionType) {
		if v, ok := message[k]; ok {
			a.Set(k, v)
		}
	}
	return a
}

func orderedUserSignKeys(actionType string) []string {
	switch actionType {
	case "usdSend":
		return []string{"destination", "amount", "time"}
	case "spotSend":
		return []string{"destination", "token", "amount", "time"}
	case "withdraw3":
		return []string{"destination", "amount", "time"}
	case "usdClassTransfer":
		return []string{"amount", "toPerp", "nonce"}
	case "tokenDelegate":
		return []string{"validator", "wei", "isUndelegate", "nonce"}
	case "approveAgent":
		return []string{"agentAddress", "agentName", "nonce"}
	case "approveBuilderFee":
		return []string{"maxFeeRate", "builder", "nonce"}
	}
	return nil
}

// UsdSend sends USDC to another Hyperliquid user.
func (c *Client) UsdSend(ctx context.Context, destination, amount string) error {
	t := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"destination":      destination,
		"amount":           amount,
		"time":             t,
	}
	action := c.userSignWrapper("usdSend", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:UsdSend", userSignActionFields.UsdSend, msg, nil)
}

// SpotSend transfers a spot token to another user.
func (c *Client) SpotSend(ctx context.Context, destination, token, amount string) error {
	t := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"destination":      destination,
		"token":            token,
		"amount":           amount,
		"time":             t,
	}
	action := c.userSignWrapper("spotSend", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:SpotSend", userSignActionFields.SpotSend, msg, nil)
}

// Withdraw3 withdraws USDC to an Arbitrum address.
func (c *Client) Withdraw3(ctx context.Context, destination, amount string) error {
	t := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"destination":      destination,
		"amount":           amount,
		"time":             t,
	}
	action := c.userSignWrapper("withdraw3", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:Withdraw", userSignActionFields.Withdraw3, msg, nil)
}

// UsdClassTransfer moves USDC between perp and spot wallets.
func (c *Client) UsdClassTransfer(ctx context.Context, amount string, toPerp bool) error {
	n := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"amount":           amount,
		"toPerp":           toPerp,
		"nonce":            n,
	}
	action := c.userSignWrapper("usdClassTransfer", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:UsdClassTransfer", userSignActionFields.UsdClassTransfer, msg, nil)
}

// TokenDelegate delegates HYPE to a validator (or undelegates if isUndelegate).
func (c *Client) TokenDelegate(ctx context.Context, validator string, wei uint64, isUndelegate bool) error {
	n := c.nonces.next()
	msg := map[string]any{
		"hyperliquidChain": c.hyperliquidChain(),
		"validator":        validator,
		"wei":              wei,
		"isUndelegate":     isUndelegate,
		"nonce":            n,
	}
	action := c.userSignWrapper("tokenDelegate", msg)
	return c.submitUser(ctx, action, "HyperliquidTransaction:TokenDelegate", userSignActionFields.TokenDelegate, msg, nil)
}

// SubAccountTransfer transfers USDC between master and sub-account.
func (c *Client) SubAccountTransfer(ctx context.Context, subAccountUser string, isDeposit bool, usd uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "subAccountTransfer")
	a.Set("subAccountUser", subAccountUser)
	a.Set("isDeposit", isDeposit)
	a.Set("usd", usd)
	return c.submitL1(ctx, a, nil)
}

// SubAccountSpotTransfer transfers a spot token between master and sub-account.
func (c *Client) SubAccountSpotTransfer(ctx context.Context, subAccountUser string, isDeposit bool, token string, amount string) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "subAccountSpotTransfer")
	a.Set("subAccountUser", subAccountUser)
	a.Set("isDeposit", isDeposit)
	a.Set("token", token)
	a.Set("amount", amount)
	return c.submitL1(ctx, a, nil)
}

// VaultTransfer deposits to or withdraws from a vault.
func (c *Client) VaultTransfer(ctx context.Context, vaultAddress string, isDeposit bool, usd uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "vaultTransfer")
	a.Set("vaultAddress", vaultAddress)
	a.Set("isDeposit", isDeposit)
	a.Set("usd", usd)
	return c.submitL1(ctx, a, nil)
}
