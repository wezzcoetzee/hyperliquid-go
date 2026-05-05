package info

import "context"

// Delegation is one user→validator delegation.
type Delegation struct {
	Validator            string `json:"validator"`
	Amount               string `json:"amount"`
	LockedUntilTimestamp int64  `json:"lockedUntilTimestamp"`
}

// Delegations returns the user's active delegations.
func (c *Client) Delegations(ctx context.Context, user string) ([]Delegation, error) {
	var out []Delegation
	if err := c.post(ctx, map[string]any{"type": "delegations", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DelegatorSummary aggregates a delegator's stake across validators.
type DelegatorSummary struct {
	Delegated              string `json:"delegated"`
	Undelegated            string `json:"undelegated"`
	TotalPendingWithdrawal string `json:"totalPendingWithdrawal"`
	NPendingWithdrawals    int    `json:"nPendingWithdrawals"`
}

// DelegatorSummary returns aggregate stake info for the user.
func (c *Client) DelegatorSummary(ctx context.Context, user string) (*DelegatorSummary, error) {
	var out DelegatorSummary
	if err := c.post(ctx, map[string]any{"type": "delegatorSummary", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DelegatorHistoryEntry is one delegation event.
type DelegatorHistoryEntry struct {
	Time  int64          `json:"time"`
	Hash  string         `json:"hash"`
	Delta map[string]any `json:"delta"`
}

// DelegatorHistory returns the user's stake-change history.
func (c *Client) DelegatorHistory(ctx context.Context, user string) ([]DelegatorHistoryEntry, error) {
	var out []DelegatorHistoryEntry
	if err := c.post(ctx, map[string]any{"type": "delegatorHistory", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DelegatorReward is one staking reward payout.
type DelegatorReward struct {
	Time        int64  `json:"time"`
	Source      string `json:"source"`
	TotalAmount string `json:"totalAmount"`
}

// DelegatorRewards returns the user's reward history.
func (c *Client) DelegatorRewards(ctx context.Context, user string) ([]DelegatorReward, error) {
	var out []DelegatorReward
	if err := c.post(ctx, map[string]any{"type": "delegatorRewards", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ValidatorSummary describes one validator.
type ValidatorSummary struct {
	Validator       string `json:"validator"`
	Signer          string `json:"signer"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	NRecentBlocks   int    `json:"nRecentBlocks"`
	Stake           string `json:"stake"`
	IsJailed        bool   `json:"isJailed"`
	UnjailableAfter int64  `json:"unjailableAfter"`
	IsActive        bool   `json:"isActive"`
	Commission      string `json:"commission"`
}

// ValidatorSummaries returns the active validator set with stake/commission info.
func (c *Client) ValidatorSummaries(ctx context.Context) ([]ValidatorSummary, error) {
	var out []ValidatorSummary
	if err := c.post(ctx, map[string]any{"type": "validatorSummaries"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
