package exchange

import (
	"context"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/signer/privkey"
)

func TestBuildOrderAction_MatchesFixtureHash(t *testing.T) {
	f := loadFixture(t, "order_l1")
	req := OrderRequest{
		Orders: []OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "30000", Sz: "0.1", ReduceOnly: false,
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		}},
		Grouping: "na",
	}
	action := buildOrderAction(req)
	nonce, _ := strconv.ParseUint(f.Nonce, 10, 64)

	got, err := ActionHash(action, nonce, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != stripHexPrefix(f.ActionHash) {
		t.Fatalf("hash mismatch:\n got %x\nwant %s", got, f.ActionHash)
	}
}

func TestClient_Order_HappyPath(t *testing.T) {
	pk, err := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}

	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action, _ := req["action"].(map[string]any)
		if action == nil || action["type"] != "order" {
			t.Fatalf("unexpected action: %v", action)
		}
		return map[string]any{
			"status": "ok",
			"response": map[string]any{
				"type": "order",
				"data": map[string]any{
					"statuses": []map[string]any{
						{"resting": map[string]any{"oid": 99}},
					},
				},
			},
		}
	})
	defer srv.Close()

	resp, err := c.Order(context.Background(), OrderRequest{
		Orders: []OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "30000", Sz: "0.1",
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		}},
		Grouping: "na",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Statuses) != 1 || resp.Statuses[0].Resting == nil || resp.Statuses[0].Resting.Oid != 99 {
		t.Fatalf("unexpected: %+v", resp.Statuses)
	}
}

func TestClient_Order_RejectsActionError(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		return map[string]any{
			"status":   "err",
			"response": "bad tick size",
		}
	})
	defer srv.Close()

	_, err := c.Order(context.Background(), OrderRequest{
		Orders: []OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "30000", Sz: "0.1",
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		}},
		Grouping: "na",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	rej, ok := err.(*ActionRejected)
	if !ok {
		t.Fatalf("err type %T (want *ActionRejected)", err)
	}
	if rej.Action != "order" {
		t.Errorf("action: %q", rej.Action)
	}
}

func TestClient_Order_NoSigner(t *testing.T) {
	c, srv := newExchangeTestClient(t, nil, SourceMainnet, 42161, func(req map[string]any) any {
		t.Fatal("server should not be hit")
		return nil
	})
	defer srv.Close()

	_, err := c.Order(context.Background(), OrderRequest{
		Orders: []OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "30000", Sz: "0.1",
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		}},
		Grouping: "na",
	})
	if err != ErrNoSigner {
		t.Fatalf("expected ErrNoSigner, got %v", err)
	}
}
