package hyperliquid

import "testing"

func TestNetwork_Endpoints(t *testing.T) {
	cases := []struct {
		name      string
		net       Network
		wantHTTP  string
		wantWS    string
	}{
		{"mainnet", Mainnet, "https://api.hyperliquid.xyz", "wss://api.hyperliquid.xyz/ws"},
		{"testnet", Testnet, "https://api.hyperliquid-testnet.xyz", "wss://api.hyperliquid-testnet.xyz/ws"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.net.HTTPURL(); got != c.wantHTTP {
				t.Errorf("HTTPURL() = %q, want %q", got, c.wantHTTP)
			}
			if got := c.net.WSURL(); got != c.wantWS {
				t.Errorf("WSURL() = %q, want %q", got, c.wantWS)
			}
		})
	}
}

func TestNetwork_ChainID(t *testing.T) {
	if Mainnet.SignatureChainID() != 42161 {
		t.Errorf("mainnet chainID = %d, want 42161", Mainnet.SignatureChainID())
	}
	if Testnet.SignatureChainID() != 421614 {
		t.Errorf("testnet chainID = %d, want 421614", Testnet.SignatureChainID())
	}
}
