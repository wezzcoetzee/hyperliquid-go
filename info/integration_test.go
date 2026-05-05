//go:build integration

package info

import (
	"context"
	"testing"
	"time"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

// TestIntegration_TestnetMeta is a smoke test that hits the real Hyperliquid
// testnet API. Run with: `go test -tags=integration ./info/...`
//
// It asserts that the Meta endpoint returns at least one perp asset.
func TestIntegration_TestnetMeta(t *testing.T) {
	c := &Client{HTTP: transport.NewDefaultHTTP("https://api.hyperliquid-testnet.xyz", nil)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	meta, err := c.Meta(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Universe) == 0 {
		t.Fatal("empty universe")
	}
	t.Logf("perp universe size: %d", len(meta.Universe))
}
