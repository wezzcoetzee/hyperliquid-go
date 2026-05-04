package hyperliquid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew_Defaults(t *testing.T) {
	c, err := New(Config{Network: Testnet})
	if err != nil {
		t.Fatal(err)
	}
	if c.Network != Testnet {
		t.Errorf("network = %v", c.Network)
	}
	if c.Info == nil || c.Exchange == nil || c.Subscriptions == nil {
		t.Fatal("method groups must be wired")
	}
}

func TestClient_InfoMeta_RoundTrip(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"universe":[{"name":"BTC","szDecimals":5,"maxLeverage":50}]}`))
	}))
	defer srv.Close()

	c, err := New(Config{Network: Testnet, BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}
	meta, err := c.Info.Meta(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Universe) != 1 || meta.Universe[0].Name != "BTC" {
		t.Fatalf("unexpected: %+v", meta)
	}
}
