package hyperliquid

import (
	"net/http"

	"github.com/wezzcoetzee/hyperliquid/exchange"
	"github.com/wezzcoetzee/hyperliquid/info"
	"github.com/wezzcoetzee/hyperliquid/transport"
	"github.com/wezzcoetzee/hyperliquid/ws"
)

// Signer is a placeholder for the signing interface, filled in by Plan 02.
// Once defined, hyperliquid.Config.Signer will hold a signer.Signer rather
// than this empty interface.
type Signer interface{}

// Config controls how a Client is constructed. The zero value has Network==Mainnet
// (see Network's doc), no Signer, and the default *http.Client (30s timeout backstop).
type Config struct {
	// Network selects mainnet or testnet endpoints. Zero value is Mainnet.
	Network Network

	// Signer signs Exchange (write) actions. Required for client.Exchange writes;
	// optional otherwise.
	Signer Signer

	// HTTP overrides the default *http.Client. Use to install proxies, retries, etc.
	HTTP *http.Client

	// BaseURL overrides Network.HTTPURL(). Used for tests and proxies.
	BaseURL string

	// WSURL overrides Network.WSURL(). Used for tests and proxies.
	WSURL string
}

// Client is the top-level Hyperliquid SDK entry point. Its three method-group
// fields are always populated by New; the zero Client is unusable.
//
// Client is safe for concurrent use.
type Client struct {
	Network       Network
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
	httpTr := transport.NewDefaultHTTP(baseURL, cfg.HTTP)
	return &Client{
		Network:       cfg.Network,
		Info:          &info.Client{HTTP: httpTr},
		Exchange:      &exchange.Client{HTTP: httpTr},
		Subscriptions: &ws.Client{},
	}, nil
}
