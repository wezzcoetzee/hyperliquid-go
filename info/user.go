package info

import (
	"context"
	"encoding/json"
)

// ClearinghouseState returns the perp clearinghouse snapshot for `user`.
func (c *Client) ClearinghouseState(ctx context.Context, user string) (*ClearinghouseState, error) {
	var out ClearinghouseState
	if err := c.post(ctx, map[string]any{"type": "clearinghouseState", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SpotClearinghouseState returns spot balances for `user`.
func (c *Client) SpotClearinghouseState(ctx context.Context, user string) (*SpotClearinghouseState, error) {
	var out SpotClearinghouseState
	if err := c.post(ctx, map[string]any{"type": "spotClearinghouseState", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// OpenOrders returns the user's resting orders.
func (c *Client) OpenOrders(ctx context.Context, user string) ([]OpenOrder, error) {
	var out []OpenOrder
	if err := c.post(ctx, map[string]any{"type": "openOrders", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// FrontendOpenOrders returns the user's resting orders with order-type metadata.
func (c *Client) FrontendOpenOrders(ctx context.Context, user string) ([]FrontendOpenOrder, error) {
	var out []FrontendOpenOrder
	if err := c.post(ctx, map[string]any{"type": "frontendOpenOrders", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserFills returns the user's recent fills.
func (c *Client) UserFills(ctx context.Context, user string) ([]Fill, error) {
	var out []Fill
	if err := c.post(ctx, map[string]any{"type": "userFills", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserFillsByTime returns the user's fills in [startMs, endMs); endMs == 0
// returns up to the current time.
func (c *Client) UserFillsByTime(ctx context.Context, user string, startMs, endMs int64) ([]Fill, error) {
	body := map[string]any{
		"type":      "userFillsByTime",
		"user":      user,
		"startTime": startMs,
	}
	if endMs > 0 {
		body["endTime"] = endMs
	}
	var out []Fill
	if err := c.post(ctx, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserFunding returns the user's funding-payment ledger entries in [startMs, endMs).
func (c *Client) UserFunding(ctx context.Context, user string, startMs, endMs int64) ([]LedgerEntry, error) {
	body := map[string]any{
		"type":      "userFunding",
		"user":      user,
		"startTime": startMs,
	}
	if endMs > 0 {
		body["endTime"] = endMs
	}
	var out []LedgerEntry
	if err := c.post(ctx, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserNonFundingLedgerUpdates returns deposits/withdrawals/transfers in [startMs, endMs).
func (c *Client) UserNonFundingLedgerUpdates(ctx context.Context, user string, startMs, endMs int64) ([]LedgerEntry, error) {
	body := map[string]any{
		"type":      "userNonFundingLedgerUpdates",
		"user":      user,
		"startTime": startMs,
	}
	if endMs > 0 {
		body["endTime"] = endMs
	}
	var out []LedgerEntry
	if err := c.post(ctx, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// UserRateLimit returns the caller's current rate-limit budget.
func (c *Client) UserRateLimit(ctx context.Context, user string) (*UserRateLimit, error) {
	var out UserRateLimit
	if err := c.post(ctx, map[string]any{"type": "userRateLimit", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// OrderStatusByOID looks up an order by its numeric oid.
func (c *Client) OrderStatusByOID(ctx context.Context, user string, oid uint64) (*OrderStatusResponse, error) {
	var out OrderStatusResponse
	if err := c.post(ctx, map[string]any{"type": "orderStatus", "user": user, "oid": oid}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// OrderStatusByCloid looks up an order by its client-supplied 0x-prefixed cloid.
func (c *Client) OrderStatusByCloid(ctx context.Context, user, cloid string) (*OrderStatusResponse, error) {
	var out OrderStatusResponse
	if err := c.post(ctx, map[string]any{"type": "orderStatus", "user": user, "oid": cloid}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HistoricalOrders returns up to 2000 of the user's historical orders.
func (c *Client) HistoricalOrders(ctx context.Context, user string) ([]HistoricalOrder, error) {
	var out []HistoricalOrder
	if err := c.post(ctx, map[string]any{"type": "historicalOrders", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// TwapHistory returns the user's TWAP orders.
func (c *Client) TwapHistory(ctx context.Context, user string) ([]TwapState, error) {
	var out []TwapState
	if err := c.post(ctx, map[string]any{"type": "twapHistory", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SubAccounts returns sub-accounts owned by the master `user`.
func (c *Client) SubAccounts(ctx context.Context, user string) ([]SubAccount, error) {
	var out []SubAccount
	if err := c.post(ctx, map[string]any{"type": "subAccounts", "user": user}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Referral returns the referral state for `user`.
func (c *Client) Referral(ctx context.Context, user string) (*ReferralState, error) {
	var out ReferralState
	if err := c.post(ctx, map[string]any{"type": "referral", "user": user}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// PortfolioPeriods returns the user's portfolio history bucketed into named periods.
func (c *Client) PortfolioPeriods(ctx context.Context, user string) (map[string]PortfolioPeriod, error) {
	var raw [][]json.RawMessage
	if err := c.post(ctx, map[string]any{"type": "portfolio", "user": user}, &raw); err != nil {
		return nil, err
	}
	out := make(map[string]PortfolioPeriod, len(raw))
	for _, pair := range raw {
		if len(pair) != 2 {
			continue
		}
		var name string
		if err := json.Unmarshal(pair[0], &name); err != nil {
			return nil, err
		}
		var period PortfolioPeriod
		if err := json.Unmarshal(pair[1], &period); err != nil {
			return nil, err
		}
		out[name] = period
	}
	return out, nil
}
