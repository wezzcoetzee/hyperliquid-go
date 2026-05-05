// Package info implements read-only Hyperliquid /info endpoints.
//
// Methods on Client require no signer and no authentication. Construct via
// hyperliquid.New; the HTTP transport is wired in by the root package.
package info

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Client is the read-only /info client. Construct via hyperliquid.New.
type Client struct {
	HTTP transport.HTTP
}

// post is the shared POST-to-/info helper used by every method on Client.
func (c *Client) post(ctx context.Context, body any, out any) error {
	return c.HTTP.PostJSON(ctx, "/info", body, out)
}
