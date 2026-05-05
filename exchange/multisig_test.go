package exchange

import (
	"context"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid/signer"
)

func TestClient_MultiSig_BuildsEnvelope(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "multiSig" {
			t.Fatalf("type: %v", action["type"])
		}
		sigs := action["signatures"].([]any)
		if len(sigs) != 2 {
			t.Fatalf("signatures len = %d", len(sigs))
		}
		payload := action["payload"].(map[string]any)
		if payload["multiSigUser"] != "0xmulti" || payload["outerSigner"] != "0xouter" {
			t.Fatalf("payload: %v", payload)
		}
		inner := payload["action"].(map[string]any)
		if inner["type"] != "noop" {
			t.Fatalf("inner: %v", inner)
		}
		return okResponse()
	})
	defer srv.Close()

	inner := msgpack.NewOrderedMap()
	inner.Set("type", "noop")

	var s1, s2 signer.Signature
	s1.R[0] = 0x11
	s2.R[0] = 0x22
	s1.V = 27
	s2.V = 28

	err := c.MultiSig(context.Background(), MultiSigParams{
		MultiSigUser:    "0xmulti",
		OuterSigner:     "0xouter",
		InnerAction:     inner,
		InnerNonce:      1,
		InnerSignatures: []signer.Signature{s1, s2},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_MultiSig_NoSigner(t *testing.T) {
	c, srv := newExchangeTestClient(t, nil, SourceMainnet, 42161, func(req map[string]any) any {
		t.Fatal("server should not be hit")
		return nil
	})
	defer srv.Close()

	inner := msgpack.NewOrderedMap()
	inner.Set("type", "noop")

	if err := c.MultiSig(context.Background(), MultiSigParams{InnerAction: inner}); err != ErrNoSigner {
		t.Fatalf("expected ErrNoSigner, got %v", err)
	}
}
