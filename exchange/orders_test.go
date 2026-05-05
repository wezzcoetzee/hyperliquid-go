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

func TestClient_Cancel(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "cancel" {
			t.Fatalf("type: %v", action["type"])
		}
		cancels := action["cancels"].([]any)
		if len(cancels) != 1 {
			t.Fatalf("cancels: %v", cancels)
		}
		first := cancels[0].(map[string]any)
		if first["a"].(float64) != 0 || first["o"].(float64) != 12345 {
			t.Fatalf("first: %v", first)
		}
		return map[string]any{
			"status": "ok",
			"response": map[string]any{
				"type": "cancel",
				"data": map[string]any{"statuses": []any{"success"}},
			},
		}
	})
	defer srv.Close()

	resp, err := c.Cancel(context.Background(), []CancelParams{{Asset: 0, Oid: 12345}})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Statuses) != 1 || resp.Statuses[0] != "success" {
		t.Fatalf("statuses: %v", resp.Statuses)
	}
}

func TestClient_CancelByCloid(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "cancelByCloid" {
			t.Fatalf("type: %v", action["type"])
		}
		cancels := action["cancels"].([]any)
		first := cancels[0].(map[string]any)
		if first["asset"].(float64) != 0 || first["cloid"] != "0x0123456789abcdef0123456789abcdef" {
			t.Fatalf("first: %v", first)
		}
		return map[string]any{
			"status": "ok",
			"response": map[string]any{
				"type": "cancel",
				"data": map[string]any{"statuses": []any{"success"}},
			},
		}
	})
	defer srv.Close()

	_, err := c.CancelByCloid(context.Background(), []CancelByCloidParams{{Asset: 0, Cloid: "0x0123456789abcdef0123456789abcdef"}})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Modify(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "modify" {
			t.Fatalf("type: %v", action["type"])
		}
		if action["oid"].(float64) != 555 {
			t.Fatalf("oid: %v", action["oid"])
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.Modify(context.Background(), ModifyParams{
		Oid: 555,
		NewOrder: OrderParams{
			Asset: 0, IsBuy: true, LimitPx: "30100", Sz: "0.1",
			OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}},
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_BatchModify(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "batchModify" {
			t.Fatalf("type: %v", action["type"])
		}
		mods := action["modifies"].([]any)
		if len(mods) != 2 {
			t.Fatalf("len(mods)=%d", len(mods))
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.BatchModify(context.Background(), []ModifyParams{
		{Oid: 1, NewOrder: OrderParams{Asset: 0, IsBuy: true, LimitPx: "1", Sz: "0.1", OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}}}},
		{Oid: 2, NewOrder: OrderParams{Asset: 0, IsBuy: false, LimitPx: "2", Sz: "0.1", OrderType: OrderType{Limit: &LimitOrder{Tif: TifGtc}}}},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ScheduleCancel(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "scheduleCancel" {
			t.Fatalf("type: %v", action["type"])
		}
		if action["time"].(float64) != 1700000000000 {
			t.Fatalf("time: %v", action["time"])
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.ScheduleCancel(context.Background(), 1700000000000); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ScheduleCancel_Disarm(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if _, present := action["time"]; present {
			t.Fatalf("time should be omitted for disarm")
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.ScheduleCancel(context.Background(), 0); err != nil {
		t.Fatal(err)
	}
}

func TestClient_UpdateLeverage(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "updateLeverage" || action["asset"].(float64) != 0 || action["isCross"] != true || action["leverage"].(float64) != 10 {
			t.Fatalf("action: %v", action)
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.UpdateLeverage(context.Background(), 0, true, 10); err != nil {
		t.Fatal(err)
	}
}

func TestClient_UpdateIsolatedMargin(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "updateIsolatedMargin" || action["ntli"].(float64) != -100 {
			t.Fatalf("action: %v", action)
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.UpdateIsolatedMargin(context.Background(), 0, true, -100); err != nil {
		t.Fatal(err)
	}
}

func TestClient_TwapOrder(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "twapOrder" {
			t.Fatalf("type: %v", action["type"])
		}
		twap := action["twap"].(map[string]any)
		if twap["m"].(float64) != 30 {
			t.Fatalf("twap: %v", twap)
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.TwapOrder(context.Background(), TwapParams{Asset: 0, IsBuy: true, Sz: "0.1", Minutes: 30, Randomize: true}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_TwapCancel(t *testing.T) {
	pk, _ := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	c, srv := newExchangeTestClient(t, pk, SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "twapCancel" || action["t"].(float64) != 7 {
			t.Fatalf("action: %v", action)
		}
		return map[string]any{"status": "ok", "response": map[string]any{}}
	})
	defer srv.Close()

	if err := c.TwapCancel(context.Background(), 0, 7); err != nil {
		t.Fatal(err)
	}
}
