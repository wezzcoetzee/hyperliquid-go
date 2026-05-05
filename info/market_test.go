package info

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/wezzcoetzee/hyperliquid-go/transport"
)

// newTestClient creates an info.Client wired to a one-off httptest.Server.
// The handler should validate the request and write the canned JSON response.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &Client{HTTP: transport.NewDefaultHTTP(srv.URL, nil)}
}

// loadResponse loads a canned response fixture from testdata/responses/.
func loadResponse(t *testing.T, name string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "responses", name))
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// readReqType decodes the POST body and returns the value of "type".
func readReqType(t *testing.T, body io.Reader) string {
	t.Helper()
	var req map[string]any
	if err := json.NewDecoder(body).Decode(&req); err != nil {
		t.Fatal(err)
	}
	s, _ := req["type"].(string)
	return s
}

func TestClient_AllMids(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "allMids" {
			t.Fatalf("unexpected type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "all_mids.json"))
	})

	mids, err := c.AllMids(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if mids["BTC"] != "30000.5" {
		t.Fatalf("BTC = %v", mids["BTC"])
	}
	if len(mids) != 3 {
		t.Errorf("expected 3 entries, got %d", len(mids))
	}
}

func TestClient_Meta(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "meta" {
			t.Fatalf("unexpected type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"universe":[{"name":"BTC","szDecimals":5,"maxLeverage":50}]}`))
	})

	meta, err := c.Meta(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Universe) != 1 || meta.Universe[0].Name != "BTC" {
		t.Fatalf("unexpected: %+v", meta)
	}
}

func TestClient_L2Book(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &req)
		if req["type"] != "l2Book" || req["coin"] != "BTC" {
			t.Fatalf("unexpected req: %v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "l2_book.json"))
	})
	book, err := c.L2Book(context.Background(), "BTC")
	if err != nil {
		t.Fatal(err)
	}
	if book.Coin != "BTC" || len(book.Levels[0]) != 1 || book.Levels[0][0].Px != "30000" {
		t.Fatalf("unexpected: %+v", book)
	}
}

func TestClient_CandleSnapshot(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "candleSnapshot" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "candle_snapshot.json"))
	})
	out, err := c.CandleSnapshot(context.Background(), CandleSnapshotReq{
		Coin: "BTC", Interval: "1m", StartMs: 1700000000000, EndMs: 1700003600000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].C != "30000" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestClient_MetaAndAssetCtxs(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "metaAndAssetCtxs" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "meta_and_asset_ctxs.json"))
	})
	meta, ctxs, err := c.MetaAndAssetCtxs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Universe) != 1 || meta.Universe[0].Name != "BTC" {
		t.Fatalf("meta: %+v", meta)
	}
	if len(ctxs) != 1 || ctxs[0].MarkPx != "30000" {
		t.Fatalf("ctxs: %+v", ctxs)
	}
}

func TestClient_SpotMeta(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "spotMeta" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "spot_meta.json"))
	})
	out, err := c.SpotMeta(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Tokens) != 1 || out.Tokens[0].Name != "USDC" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestClient_SpotMetaAndAssetCtxs(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "spotMetaAndAssetCtxs" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "spot_meta_and_asset_ctxs.json"))
	})
	_, ctxs, err := c.SpotMetaAndAssetCtxs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(ctxs) != 1 || ctxs[0].Coin != "PURR/USDC" {
		t.Fatalf("ctxs: %+v", ctxs)
	}
}

func TestClient_FundingHistory(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &req)
		if req["type"] != "fundingHistory" || req["coin"] != "BTC" {
			t.Fatalf("req: %v", req)
		}
		if got := req["startTime"]; got == nil {
			t.Fatalf("missing startTime")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "funding_history.json"))
	})
	out, err := c.FundingHistory(context.Background(), "BTC", 1700000000000, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].FundingRate != "0.0001" {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestClient_PredictedFundings(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "predictedFundings" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "predicted_fundings.json"))
	})
	out, err := c.PredictedFundings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Fatalf("len = %d", len(out))
	}
}
