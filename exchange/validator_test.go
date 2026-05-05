package exchange

import (
	"context"
	"testing"
)

func TestClient_CDeposit(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "cDeposit" || action["wei"].(float64) != 1000 {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.CDeposit(context.Background(), 1000); err != nil {
		t.Fatal(err)
	}
}

func TestClient_CWithdraw(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "cWithdraw" || action["wei"].(float64) != 500 {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.CWithdraw(context.Background(), 500); err != nil {
		t.Fatal(err)
	}
}
