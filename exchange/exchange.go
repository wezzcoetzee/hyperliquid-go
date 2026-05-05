// Package exchange will implement signed Hyperliquid /exchange actions.
//
// This file currently contains only the Client shell so the root package can
// reference it. Trading methods, the Signer wiring, nonce generator, and
// signing pipeline are added in Plans 02 and 04.
package exchange

import (
	"github.com/wezzcoetzee/hyperliquid/signer"
	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Client is the signed /exchange client. Construct via hyperliquid.New.
// Trading methods will be added in Plan 04.
type Client struct {
	HTTP   transport.HTTP
	Signer signer.Signer
}
