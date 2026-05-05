package hyperliquid

import (
	"context"
	"errors"
	"net/http"

	"github.com/wezzcoetzee/hyperliquid-go/exchange"
	"github.com/wezzcoetzee/hyperliquid-go/info"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
	"github.com/wezzcoetzee/hyperliquid-go/transport"
	"github.com/wezzcoetzee/hyperliquid-go/ws"
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
	WSURL string

	// WebSocketPosts routes Exchange POST requests through the WebSocket
	// connection instead of HTTP. Requires the WebSocket connection to be
	// established first (via a subscription or explicit dial).
	//
	// Note: WS-based posting is exposed as ws.Client.Post for now.
	// Full HTTP-transport routing via WS is a planned follow-up.
	WebSocketPosts bool
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
	wsURL := cfg.WSURL
	if wsURL == "" {
		wsURL = cfg.Network.WSURL()
	}
	httpTr := &wrappingHTTP{inner: transport.NewDefaultHTTP(baseURL, cfg.HTTP)}
	exchClient := &exchange.Client{HTTP: httpTr, Signer: cfg.Signer}
	switch cfg.Network {
	case Testnet:
		exchClient.Source = exchange.SourceTestnet
		exchClient.SignatureChainID = 421614
	default:
		exchClient.Source = exchange.SourceMainnet
		exchClient.SignatureChainID = 42161
	}
	return &Client{
		Network:       cfg.Network,
		Signer:        cfg.Signer,
		Info:          &info.Client{HTTP: httpTr},
		Exchange:      exchClient,
		Subscriptions: &ws.Client{URL: wsURL},
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
