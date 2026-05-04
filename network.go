package hyperliquid

type Network int

const (
	Mainnet Network = iota
	Testnet
)

func (n Network) HTTPURL() string {
	switch n {
	case Testnet:
		return "https://api.hyperliquid-testnet.xyz"
	default:
		return "https://api.hyperliquid.xyz"
	}
}

func (n Network) WSURL() string {
	switch n {
	case Testnet:
		return "wss://api.hyperliquid-testnet.xyz/ws"
	default:
		return "wss://api.hyperliquid.xyz/ws"
	}
}

func (n Network) SignatureChainID() uint64 {
	switch n {
	case Testnet:
		return 421614
	default:
		return 42161
	}
}

func (n Network) String() string {
	if n == Testnet {
		return "testnet"
	}
	return "mainnet"
}
