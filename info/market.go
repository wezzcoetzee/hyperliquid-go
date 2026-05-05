package info

import "context"

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
