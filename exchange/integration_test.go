//go:build integration

package exchange

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/wezzcoetzee/hyperliquid-go/signer/privkey"
	"github.com/wezzcoetzee/hyperliquid-go/transport"
)

// TestIntegration_TestnetOrderCancel places a way-out-of-market limit order
// on the Hyperliquid testnet and immediately cancels it. Run with:
//
//   HL_TESTNET_KEY=0x... go test -tags=integration ./exchange/...
//
// HL_TESTNET_KEY must be a funded testnet account.
func TestIntegration_TestnetOrderCancel(t *testing.T) {
	pk := os.Getenv("HL_TESTNET_KEY")
	if pk == "" {
		t.Skip("HL_TESTNET_KEY not set")
	}
	s, err := privkey.New(pk)
	if err != nil {
		t.Fatal(err)
	}
	c := &Client{
		HTTP:             transport.NewDefaultHTTP("https://api.hyperliquid-testnet.xyz", nil),
		Signer:           s,
		Source:           SourceTestnet,
		SignatureChainID: 421614,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	resp, err := c.Order(ctx, OrderRequest{
		Orders: []OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "1", Sz: "0.001",
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		}},
		Grouping: "na",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Statuses) != 1 || resp.Statuses[0].Resting == nil {
		t.Fatalf("expected resting: %+v", resp.Statuses)
	}
	oid := resp.Statuses[0].Resting.Oid
	t.Logf("placed oid=%d", oid)

	if _, err := c.Cancel(ctx, []CancelParams{{Asset: 0, Oid: oid}}); err != nil {
		t.Fatal(err)
	}
}
