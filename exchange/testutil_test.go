package exchange

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/signer"
	"github.com/wezzcoetzee/hyperliquid/transport"
)

func newExchangeTestClient(t *testing.T, s signer.Signer, source Source, sigChainID uint64, respond func(map[string]any) any) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &req)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(respond(req))
	}))
	c := &Client{
		HTTP:             transport.NewDefaultHTTP(srv.URL, nil),
		Signer:           s,
		Source:           source,
		SignatureChainID: sigChainID,
	}
	return c, srv
}
