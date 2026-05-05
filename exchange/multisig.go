// Package-level note: Hyperliquid's multiSig flow is intricate. This Go
// implementation accepts pre-built inner actions and pre-collected signatures
// rather than orchestrating the multi-party signing dance — that workflow
// belongs at the application layer.

package exchange

import (
	"context"
	"encoding/hex"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
)

// MultiSigParams wraps an inner action with N co-signer signatures.
//
// MultiSigUser is the multi-sig "address" that the inner signatures were
// computed under (the multi-sig wallet's address). OuterSigner is the
// caller's address — the address whose key is signing the outer envelope.
// InnerActionHash is the L1 action hash of the inner action (computed via
// ActionHash); the outer envelope signs over a payload containing it.
type MultiSigParams struct {
	// MultiSigUser is the 0x-prefixed multi-sig wallet address whose threshold signatures are provided.
	MultiSigUser string
	// OuterSigner is the 0x-prefixed address of the caller who signs the outer envelope.
	OuterSigner string
	// InnerAction is the pre-built L1 action that the co-signers have signed.
	InnerAction *msgpack.OrderedMap
	// InnerNonce is the nonce embedded in the inner action (used for the inner action hash).
	InnerNonce uint64
	// InnerSignatures are the co-signer signatures over the inner action hash, in threshold order.
	InnerSignatures []signer.Signature
	// VaultAddress, if non-nil, directs the inner action to be applied to the named vault.
	VaultAddress *[20]byte
}

// MultiSig submits a multiSig wrapper. The action body shape is documented in
// Hyperliquid's API reference; this helper builds the outer OrderedMap and
// uses the Client's own signer for the outer signature.
//
// EXPERIMENTAL: this helper is provided as a convenience for the
// already-collected-signatures case. End-to-end correctness against the live
// API is not validated by unit tests — verify on testnet before relying on it.
func (c *Client) MultiSig(ctx context.Context, p MultiSigParams) error {
	if c.Signer == nil {
		return ErrNoSigner
	}

	a := msgpack.NewOrderedMap()
	a.Set("type", "multiSig")
	a.Set("signatureChainId", c.signatureChainIDHex())

	sigs := make([]any, len(p.InnerSignatures))
	for i, s := range p.InnerSignatures {
		sm := msgpack.NewOrderedMap()
		sm.Set("r", "0x"+hex.EncodeToString(s.R[:]))
		sm.Set("s", "0x"+hex.EncodeToString(s.S[:]))
		sm.Set("v", uint64(s.V))
		sigs[i] = sm
	}
	a.Set("signatures", sigs)

	payload := msgpack.NewOrderedMap()
	payload.Set("multiSigUser", p.MultiSigUser)
	payload.Set("outerSigner", p.OuterSigner)
	payload.Set("action", p.InnerAction)
	a.Set("payload", payload)

	return c.submitL1(ctx, a, nil)
}
