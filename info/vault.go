package info

import (
	"context"
	"encoding/json"
)

// VaultDetails describes a single vault's metadata, performance, and follower set.
type VaultDetails struct {
	Name                  string          `json:"name"`
	VaultAddress          string          `json:"vaultAddress"`
	Leader                string          `json:"leader"`
	Description           string          `json:"description"`
	PortfolioPeriods      json.RawMessage `json:"portfolio"`
	Apr                   float64         `json:"apr"`
	FollowerState         json.RawMessage `json:"followerState"`
	LeaderFraction        float64         `json:"leaderFraction"`
	LeaderCommission      float64         `json:"leaderCommission"`
	Followers             []VaultFollower `json:"followers"`
	MaxDistributable      string          `json:"maxDistributable"`
	MaxWithdrawable       string          `json:"maxWithdrawable"`
	IsClosed              bool            `json:"isClosed"`
	AllowDeposits         bool            `json:"allowDeposits"`
	AlwaysCloseOnWithdraw bool            `json:"alwaysCloseOnWithdraw"`
}

// VaultFollower describes one user following a vault.
type VaultFollower struct {
	User           string `json:"user"`
	VaultEquity    string `json:"vaultEquity"`
	Pnl            string `json:"pnl"`
	AllTimePnl     string `json:"allTimePnl"`
	DaysFollowing  int    `json:"daysFollowing"`
	VaultEntryTime int64  `json:"vaultEntryTime"`
	LockupUntil    int64  `json:"lockupUntil"`
}

// VaultDetails returns details for a vault. user (caller address) is optional;
// pass "" to omit.
func (c *Client) VaultDetails(ctx context.Context, vaultAddress, user string) (*VaultDetails, error) {
	body := map[string]any{"type": "vaultDetails", "vaultAddress": vaultAddress}
	if user != "" {
		body["user"] = user
	}
	var out VaultDetails
	if err := c.post(ctx, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// VaultSummary is the brief vault listing entry.
type VaultSummary struct {
	Name             string          `json:"name"`
	VaultAddress     string          `json:"vaultAddress"`
	Leader           string          `json:"leader"`
	TVL              string          `json:"tvl"`
	IsClosed         bool            `json:"isClosed"`
	Relationship     json.RawMessage `json:"relationship"`
	CreateTimeMillis int64           `json:"createTimeMillis"`
}

// VaultSummaries returns a brief listing of all vaults.
func (c *Client) VaultSummaries(ctx context.Context) ([]VaultSummary, error) {
	var out []VaultSummary
	if err := c.post(ctx, map[string]any{"type": "vaultSummaries"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// VaultEquity is the user's equity in one vault.
type VaultEquity struct {
	VaultAddress         string `json:"vaultAddress"`
	Equity               string `json:"equity"`
	LockedUntilTimestamp int64  `json:"lockedUntilTimestamp"`
}

// UserVaultEquities returns the user's equity in each vault they follow.
func (c *Client) UserVaultEquities(ctx context.Context, user string) ([]VaultEquity, error) {
	var out []VaultEquity
	if err := c.post(ctx, map[string]any{"type": "userVaultEquities", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
