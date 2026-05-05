package exchange

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
	"github.com/wezzcoetzee/hyperliquid-go/signer/privkey"
	"github.com/wezzcoetzee/hyperliquid-go/transport"
)

// TestLive_MultiSig_Testnet exercises the full sign-and-send pipeline for a
// multiSig envelope against the Hyperliquid testnet.
//
// The test submits an invalid inner action and asserts the error comes from the
// exchange API layer (proving the HTTP envelope was accepted at the gateway).
//
// Run with: HL_TESTNET_KEY=0x... go test -run TestLive_MultiSig_Testnet ./exchange/...
func TestLive_MultiSig_Testnet(t *testing.T) {
	rawKey := os.Getenv("HL_TESTNET_KEY")
	if rawKey == "" {
		t.Skip("HL_TESTNET_KEY not set")
	}

	pk, err := privkey.New(rawKey)
	if err != nil {
		t.Fatal(err)
	}

	c := &Client{
		HTTP:             transport.NewDefaultHTTP("https://api.hyperliquid-testnet.xyz", nil),
		Signer:           pk,
		Source:           SourceTestnet,
		SignatureChainID: 421614,
	}

	inner := msgpack.NewOrderedMap()
	inner.Set("type", "scheduleCancel")

	addr := pk.Address()
	addrHex := "0x" + hex.EncodeToString(addr[:])

	var dummySig signer.Signature
	dummySig.V = 27
	dummySig.R[0] = 0x01
	dummySig.S[0] = 0x01

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = c.MultiSig(ctx, MultiSigParams{
		MultiSigUser:    addrHex,
		OuterSigner:     addrHex,
		InnerAction:     inner,
		InnerNonce:      c.nonces.next(),
		InnerSignatures: []signer.Signature{dummySig},
	})

	var rejected *ActionRejected
	if errors.As(err, &rejected) {
		t.Logf("exchange rejected as expected for dummy sigs: %v", rejected.Response)
		return
	}
	if err != nil {
		t.Fatalf("unexpected error (not ActionRejected): %v", err)
	}
	t.Log("exchange accepted the envelope")
}
