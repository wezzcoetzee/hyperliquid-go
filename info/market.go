package info

import (
	"context"
	"encoding/json"
	"fmt"
)

// Mids is a map from asset symbol to mid price (as a string-encoded decimal).
type Mids map[string]string

// AllMids returns the current mid price for every perp/spot asset.
func (c *Client) AllMids(ctx context.Context) (Mids, error) {
	var out Mids
	if err := c.post(ctx, map[string]any{"type": "allMids"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AssetInfo describes a single perp asset in the universe.
type AssetInfo struct {
	Name        string `json:"name"`
	SzDecimals  int    `json:"szDecimals"`
	MaxLeverage int    `json:"maxLeverage"`
}

// Meta is the perp-universe metadata response.
type Meta struct {
	Universe []AssetInfo `json:"universe"`
}

// Meta returns the perp-asset universe and per-asset metadata.
func (c *Client) Meta(ctx context.Context) (*Meta, error) {
	var out Meta
	if err := c.post(ctx, map[string]any{"type": "meta"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// L2Level is one price level in the L2 order book.
type L2Level struct {
	// Px is the price as a decimal string.
	Px string `json:"px"`
	// Sz is the aggregate size at this price level as a decimal string.
	Sz string `json:"sz"`
	// N is the number of resting orders at this level.
	N int `json:"n"`
}

// L2Book is a snapshot of a coin's order book; Levels[0] is bids, Levels[1] is asks.
type L2Book struct {
	Coin   string       `json:"coin"`
	Time   int64        `json:"time"`
	Levels [2][]L2Level `json:"levels"`
}

// L2Book returns an L2 order-book snapshot for a single coin.
func (c *Client) L2Book(ctx context.Context, coin string) (*L2Book, error) {
	var out L2Book
	if err := c.post(ctx, map[string]any{"type": "l2Book", "coin": coin}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Candle is a single OHLCV bar.
type Candle struct {
	// T is the candle open timestamp in milliseconds.
	T int64 `json:"t"`
	// C is the close price as a decimal string.
	C string `json:"c"`
	// H is the high price as a decimal string.
	H string `json:"h"`
	// L is the low price as a decimal string.
	L string `json:"l"`
	// O is the open price as a decimal string.
	O string `json:"o"`
	// V is the volume as a decimal string.
	V string `json:"v"`
	// N is the number of trades in the bar.
	N int `json:"n"`
	// I is the candle interval string (e.g. "1m", "1h"); omitted in snapshots.
	I string `json:"i,omitempty"`
	// S is the coin/symbol; omitted in snapshots.
	S string `json:"s,omitempty"`
}

// CandleSnapshotReq parameterizes the candleSnapshot endpoint.
type CandleSnapshotReq struct {
	Coin     string `json:"coin"`
	Interval string `json:"interval"`
	StartMs  int64  `json:"startTime"`
	EndMs    int64  `json:"endTime"`
}

// CandleSnapshot returns historical OHLCV bars.
func (c *Client) CandleSnapshot(ctx context.Context, req CandleSnapshotReq) ([]Candle, error) {
	var out []Candle
	if err := c.post(ctx, map[string]any{"type": "candleSnapshot", "req": req}, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AssetCtx is the perp asset context (mark, oracle, funding, OI, etc.).
type AssetCtx struct {
	DayNtlVlm    string   `json:"dayNtlVlm"`
	Funding      string   `json:"funding"`
	MarkPx       string   `json:"markPx"`
	MidPx        string   `json:"midPx"`
	OpenInterest string   `json:"openInterest"`
	OraclePx     string   `json:"oraclePx"`
	Premium      string   `json:"premium"`
	PrevDayPx    string   `json:"prevDayPx"`
	ImpactPxs    []string `json:"impactPxs,omitempty"`
}

// MetaAndAssetCtxs returns the perp universe alongside per-asset live context.
func (c *Client) MetaAndAssetCtxs(ctx context.Context) (*Meta, []AssetCtx, error) {
	var raw []json.RawMessage
	if err := c.post(ctx, map[string]any{"type": "metaAndAssetCtxs"}, &raw); err != nil {
		return nil, nil, err
	}
	if len(raw) != 2 {
		return nil, nil, fmt.Errorf("info: metaAndAssetCtxs expected 2 elements, got %d", len(raw))
	}
	var meta Meta
	var ctxs []AssetCtx
	if err := json.Unmarshal(raw[0], &meta); err != nil {
		return nil, nil, fmt.Errorf("info: decode meta: %w", err)
	}
	if err := json.Unmarshal(raw[1], &ctxs); err != nil {
		return nil, nil, fmt.Errorf("info: decode asset ctxs: %w", err)
	}
	return &meta, ctxs, nil
}

// SpotToken describes a spot token in the universe.
type SpotToken struct {
	Name        string `json:"name"`
	SzDecimals  int    `json:"szDecimals"`
	WeiDecimals int    `json:"weiDecimals"`
	Index       int    `json:"index"`
	TokenID     string `json:"tokenId"`
	IsCanonical bool   `json:"isCanonical"`
	EvmContract string `json:"evmContract,omitempty"`
}

// SpotPair is one spot trading pair.
type SpotPair struct {
	Tokens [2]int `json:"tokens"`
	Name   string `json:"name"`
	Index  int    `json:"index"`
}

// SpotMeta is the spot universe response.
type SpotMeta struct {
	Tokens   []SpotToken `json:"tokens"`
	Universe []SpotPair  `json:"universe"`
}

// SpotMeta returns the spot universe.
func (c *Client) SpotMeta(ctx context.Context) (*SpotMeta, error) {
	var out SpotMeta
	if err := c.post(ctx, map[string]any{"type": "spotMeta"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SpotAssetCtx is the spot asset context.
type SpotAssetCtx struct {
	Coin              string `json:"coin"`
	CirculatingSupply string `json:"circulatingSupply"`
	DayNtlVlm         string `json:"dayNtlVlm"`
	MarkPx            string `json:"markPx"`
	MidPx             string `json:"midPx"`
	PrevDayPx         string `json:"prevDayPx"`
	TotalSupply       string `json:"totalSupply"`
}

// SpotMetaAndAssetCtxs returns the spot universe alongside per-pair context.
func (c *Client) SpotMetaAndAssetCtxs(ctx context.Context) (*SpotMeta, []SpotAssetCtx, error) {
	var raw []json.RawMessage
	if err := c.post(ctx, map[string]any{"type": "spotMetaAndAssetCtxs"}, &raw); err != nil {
		return nil, nil, err
	}
	if len(raw) != 2 {
		return nil, nil, fmt.Errorf("info: spotMetaAndAssetCtxs expected 2 elements, got %d", len(raw))
	}
	var meta SpotMeta
	var ctxs []SpotAssetCtx
	if err := json.Unmarshal(raw[0], &meta); err != nil {
		return nil, nil, fmt.Errorf("info: decode spot meta: %w", err)
	}
	if err := json.Unmarshal(raw[1], &ctxs); err != nil {
		return nil, nil, fmt.Errorf("info: decode spot asset ctxs: %w", err)
	}
	return &meta, ctxs, nil
}

// Funding is one funding-rate observation.
type Funding struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Premium     string `json:"premium"`
	Time        int64  `json:"time"`
}

// FundingHistory returns funding rates for a coin between [startMs, endMs).
// endMs == 0 returns up to the current time.
func (c *Client) FundingHistory(ctx context.Context, coin string, startMs, endMs int64) ([]Funding, error) {
	body := map[string]any{
		"type":      "fundingHistory",
		"coin":      coin,
		"startTime": startMs,
	}
	if endMs > 0 {
		body["endTime"] = endMs
	}
	var out []Funding
	if err := c.post(ctx, body, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// PredictedFundings returns the predicted-funding payload as raw JSON.
func (c *Client) PredictedFundings(ctx context.Context) ([]json.RawMessage, error) {
	var out []json.RawMessage
	if err := c.post(ctx, map[string]any{"type": "predictedFundings"}, &out); err != nil {
		return nil, err
	}
	return out, nil
}
