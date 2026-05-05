package ws

import (
	"context"
	"encoding/json"
)

// AllMids subscribes to the all-mids channel.
func (c *Client) AllMids(ctx context.Context, h func(AllMidsEvent)) (*Subscription, error) {
	return c.subscribe(ctx, "allMids", map[string]any{"type": "allMids"}, func(data []byte) {
		var e AllMidsEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// Notification subscribes to per-user system notifications.
func (c *Client) Notification(ctx context.Context, user string, h func(Notification)) (*Subscription, error) {
	key := "notification:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "notification", "user": user}, func(data []byte) {
		var e Notification
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// WebData2 subscribes to the consolidated web-frontend feed.
func (c *Client) WebData2(ctx context.Context, user string, h func(json.RawMessage)) (*Subscription, error) {
	key := "webData2:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "webData2", "user": user}, func(data []byte) {
		h(json.RawMessage(data))
	})
}

// Trades subscribes to per-coin trade events.
func (c *Client) Trades(ctx context.Context, coin string, h func([]Trade)) (*Subscription, error) {
	key := "trades:" + coin
	return c.subscribe(ctx, key, map[string]any{"type": "trades", "coin": coin}, func(data []byte) {
		var t []Trade
		if err := json.Unmarshal(data, &t); err == nil {
			h(t)
		}
	})
}

// L2Book subscribes to per-coin L2 book snapshots.
func (c *Client) L2Book(ctx context.Context, coin string, h func(L2BookEvent)) (*Subscription, error) {
	key := "l2Book:" + coin
	return c.subscribe(ctx, key, map[string]any{"type": "l2Book", "coin": coin}, func(data []byte) {
		var e L2BookEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// Candle subscribes to per-coin per-interval candle updates.
func (c *Client) Candle(ctx context.Context, coin, interval string, h func(CandleEvent)) (*Subscription, error) {
	key := "candle:" + coin + ":" + interval
	return c.subscribe(ctx, key, map[string]any{"type": "candle", "coin": coin, "interval": interval}, func(data []byte) {
		var e CandleEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// Bbo subscribes to per-coin best-bid/best-offer updates.
func (c *Client) Bbo(ctx context.Context, coin string, h func(BboEvent)) (*Subscription, error) {
	key := "bbo:" + coin
	return c.subscribe(ctx, key, map[string]any{"type": "bbo", "coin": coin}, func(data []byte) {
		var e BboEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// OrderUpdates subscribes to the user's order-status updates.
func (c *Client) OrderUpdates(ctx context.Context, user string, h func([]OrderUpdate)) (*Subscription, error) {
	key := "orderUpdates:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "orderUpdates", "user": user}, func(data []byte) {
		var e []OrderUpdate
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// UserEvents subscribes to assorted user events (liquidations, funding, etc.).
func (c *Client) UserEvents(ctx context.Context, user string, h func(UserEvent)) (*Subscription, error) {
	key := "userEvents:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "userEvents", "user": user}, func(data []byte) {
		var e UserEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// UserFills subscribes to the user's fills.
func (c *Client) UserFills(ctx context.Context, user string, h func(UserFillsEvent)) (*Subscription, error) {
	key := "userFills:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "userFills", "user": user}, func(data []byte) {
		var e UserFillsEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// UserFundings subscribes to the user's funding ledger.
func (c *Client) UserFundings(ctx context.Context, user string, h func(UserFundingsEvent)) (*Subscription, error) {
	key := "userFundings:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "userFundings", "user": user}, func(data []byte) {
		var e UserFundingsEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// UserNonFundingLedgerUpdates subscribes to deposits/withdrawals/transfers.
func (c *Client) UserNonFundingLedgerUpdates(ctx context.Context, user string, h func(UserLedgerEvent)) (*Subscription, error) {
	key := "userNonFundingLedgerUpdates:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "userNonFundingLedgerUpdates", "user": user}, func(data []byte) {
		var e UserLedgerEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// ActiveAssetCtx subscribes to per-coin active-asset context.
func (c *Client) ActiveAssetCtx(ctx context.Context, coin string, h func(ActiveAssetCtxEvent)) (*Subscription, error) {
	key := "activeAssetCtx:" + coin
	return c.subscribe(ctx, key, map[string]any{"type": "activeAssetCtx", "coin": coin}, func(data []byte) {
		var e ActiveAssetCtxEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// ActiveAssetData subscribes to per-user per-asset trade-state data.
func (c *Client) ActiveAssetData(ctx context.Context, user, coin string, h func(ActiveAssetDataEvent)) (*Subscription, error) {
	key := "activeAssetData:" + user + ":" + coin
	return c.subscribe(ctx, key, map[string]any{"type": "activeAssetData", "user": user, "coin": coin}, func(data []byte) {
		var e ActiveAssetDataEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// TwapSliceFills subscribes to the user's TWAP slice fills.
func (c *Client) TwapSliceFills(ctx context.Context, user string, h func(TwapSliceFillsEvent)) (*Subscription, error) {
	key := "twapSliceFills:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "twapSliceFills", "user": user}, func(data []byte) {
		var e TwapSliceFillsEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}

// TwapHistory subscribes to the user's TWAP history.
func (c *Client) TwapHistory(ctx context.Context, user string, h func(TwapHistoryEvent)) (*Subscription, error) {
	key := "twapHistory:" + user
	return c.subscribe(ctx, key, map[string]any{"type": "twapHistory", "user": user}, func(data []byte) {
		var e TwapHistoryEvent
		if err := json.Unmarshal(data, &e); err == nil {
			h(e)
		}
	})
}
