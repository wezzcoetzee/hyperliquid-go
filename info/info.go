// Package info implements read-only Hyperliquid /info endpoints.
//
// Methods on Client require no signer and no authentication. Later plans
// expand the surface; this file scaffolds the package and one method (Meta)
// so the root Client can be wired end-to-end.
package info

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Client is the read-only /info client. Construct via hyperliquid.New.
type Client struct {
	HTTP transport.HTTP
}

// Meta is the perp-universe metadata response.
type Meta struct {
	Universe []AssetInfo `json:"universe"`
}

// AssetInfo describes a single perp asset in the universe.
type AssetInfo struct {
	Name        string `json:"name"`
	SzDecimals  int    `json:"szDecimals"`
	MaxLeverage int    `json:"maxLeverage"`
}

// Meta returns the perp-asset universe and per-asset metadata.
func (c *Client) Meta(ctx context.Context) (*Meta, error) {
	var out Meta
	if err := c.HTTP.PostJSON(ctx, "/info", map[string]any{"type": "meta"}, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
