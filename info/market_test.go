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

	"github.com/wezzcoetzee/hyperliquid/transport"
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
