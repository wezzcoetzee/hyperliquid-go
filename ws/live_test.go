package ws

import (
	"context"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/wezzcoetzee/hyperliquid-go/signer/privkey"
)

// TestLive_OrderUpdates_Testnet subscribes to orderUpdates for an address
// derived from HL_TESTNET_KEY, waits briefly, and asserts no protocol error.
//
// Run with: HL_TESTNET_KEY=0x... go test -run TestLive_OrderUpdates_Testnet ./ws/...
func TestLive_OrderUpdates_Testnet(t *testing.T) {
	rawKey := os.Getenv("HL_TESTNET_KEY")
	if rawKey == "" {
		t.Skip("HL_TESTNET_KEY not set")
	}

	pk, err := privkey.New(rawKey)
	if err != nil {
		t.Fatal(err)
	}
	addr := pk.Address()
	user := "0x" + hex.EncodeToString(addr[:])

	c := &Client{URL: "wss://api.hyperliquid-testnet.xyz/ws"}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	received := make(chan struct{}, 1)
	sub, err := c.OrderUpdates(ctx, user, func(updates []OrderUpdate) {
		select {
		case received <- struct{}{}:
		default:
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe(context.Background())

	select {
	case <-received:
		t.Log("received order update")
	case <-time.After(5 * time.Second):
		t.Log("no order updates received (account may have no open orders — that is OK)")
	case <-ctx.Done():
		t.Fatal("context expired before subscription was established")
	}
}

