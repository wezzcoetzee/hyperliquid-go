package hyperliquid

// Network identifies which Hyperliquid environment to talk to.
//
// The zero value is Mainnet, so an uninitialized hyperliquid.Config{} defaults
// to mainnet. Any value other than Testnet is treated as Mainnet by HTTPURL,
// WSURL, SignatureChainID, and String.
type Network int

const (
	// Mainnet targets the production Hyperliquid network (Arbitrum One, chainId 42161).
	Mainnet Network = iota
	// Testnet targets the Hyperliquid test network (Arbitrum Sepolia, chainId 421614).
	Testnet
)

// HTTPURL returns the base REST API URL for the network.
func (n Network) HTTPURL() string {
	switch n {
	case Testnet:
		return "https://api.hyperliquid-testnet.xyz"
	default:
		return "https://api.hyperliquid.xyz"
	}
}

// WSURL returns the WebSocket endpoint URL for the network.
func (n Network) WSURL() string {
	switch n {
	case Testnet:
		return "wss://api.hyperliquid-testnet.xyz/ws"
	default:
		return "wss://api.hyperliquid.xyz/ws"
	}
}

// SignatureChainID returns the EIP-712 domain chainId for the network.
func (n Network) SignatureChainID() uint64 {
	switch n {
	case Testnet:
		return 421614
	default:
		return 42161
	}
}

// String returns "mainnet" or "testnet".
func (n Network) String() string {
	if n == Testnet {
		return "testnet"
	}
	return "mainnet"
}
