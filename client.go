package hyperliquid

import (
	"context"
	"errors"
	"net/http"

	"github.com/wezzcoetzee/hyperliquid/exchange"
	"github.com/wezzcoetzee/hyperliquid/info"
	"github.com/wezzcoetzee/hyperliquid/signer"
	"github.com/wezzcoetzee/hyperliquid/transport"
	"github.com/wezzcoetzee/hyperliquid/ws"
)

// Config controls how a Client is constructed. The zero value has Network==Mainnet
// (see Network's doc), no Signer, and the default *http.Client (30s timeout backstop).
type Config struct {
	// Network selects mainnet or testnet endpoints. Zero value is Mainnet.
	Network Network

	// Signer signs Exchange (write) actions. Required for client.Exchange writes;
	// optional otherwise.
	Signer signer.Signer

	// HTTP overrides the default *http.Client. Use to install proxies, retries, etc.
	HTTP *http.Client

	// BaseURL overrides Network.HTTPURL(). Used for tests and proxies.
	BaseURL string

	// WSURL overrides Network.WSURL(). Used for tests and proxies.
	// Currently accepted but unused; wired into the WebSocket client in Plan 05.
	WSURL string
}

// Client is the top-level Hyperliquid SDK entry point. Its three method-group
// fields are always populated by New; the zero Client is unusable.
//
// Client is safe for concurrent use.
type Client struct {
	Network       Network
	Signer        signer.Signer
	Info          *info.Client
	Exchange      *exchange.Client
	Subscriptions *ws.Client
}

// New returns a configured Client. It never returns a non-nil error today, but
// the signature reserves room for future validation (e.g., key derivation).
func New(cfg Config) (*Client, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = cfg.Network.HTTPURL()
	}
	httpTr := &wrappingHTTP{inner: transport.NewDefaultHTTP(baseURL, cfg.HTTP)}
	return &Client{
		Network:       cfg.Network,
		Signer:        cfg.Signer,
		Info:          &info.Client{HTTP: httpTr},
		Exchange:      &exchange.Client{HTTP: httpTr, Signer: cfg.Signer},
		Subscriptions: &ws.Client{},
	}, nil
}

// wrappingHTTP adapts a transport.HTTP to convert *transport.TransportAPIError
// into the public *APIError type at the package boundary.
type wrappingHTTP struct {
	inner transport.HTTP
}

func (w *wrappingHTTP) PostJSON(ctx context.Context, path string, body any, out any) error {
	err := w.inner.PostJSON(ctx, path, body, out)
	if err == nil {
		return nil
	}
	var tErr *transport.TransportAPIError
	if errors.As(err, &tErr) {
		return &APIError{Status: tErr.Status, Body: tErr.Body}
	}
	return err
}
