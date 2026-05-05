// Package exchange implements signed Hyperliquid /exchange actions.
package exchange

import (
	"github.com/wezzcoetzee/hyperliquid/signer"
	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Client is the signed /exchange client. Construct via hyperliquid.New.
//
// Source selects mainnet ("a") vs testnet ("b") for L1 signing.
// SignatureChainID is the EIP-712 domain chainId for user-signed actions
// (typically 42161 for Arbitrum mainnet, 421614 for Arbitrum Sepolia).
type Client struct {
	HTTP             transport.HTTP
	Signer           signer.Signer
	Source           Source
	SignatureChainID uint64

	// VaultAddress, if set, signs every action on behalf of the vault.
	VaultAddress *[20]byte

	nonces nonceGen
}
