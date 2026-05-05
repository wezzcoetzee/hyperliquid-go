package exchange

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"

	"golang.org/x/crypto/sha3"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
)

// ErrNoSigner is returned when a signing call is made without a Signer wired
// into the Client. Wrapped here rather than in the root package to avoid an
// import cycle.
var ErrNoSigner = errors.New("exchange: Signer required for write operations")

const (
	l1DomainName    = "Exchange"
	l1DomainVersion = "1"
	l1DomainChain   = 1337
	zeroAddress     = "0x0000000000000000000000000000000000000000"

	userSignDomainName    = "HyperliquidSignTransaction"
	userSignDomainVersion = "1"
)

var l1AgentTypes = signer.Types{
	"EIP712Domain": {
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"Agent": {
		{Name: "source", Type: "string"},
		{Name: "connectionId", Type: "bytes32"},
	},
}

// Source identifies the EIP-712 Agent.source byte used for L1 signing.
// Hyperliquid uses "a" on mainnet and "b" on testnet.
type Source string

const (
	SourceMainnet Source = "a"
	SourceTestnet Source = "b"
)

func BuildL1Signature(ctx context.Context, s signer.Signer, action *msgpack.OrderedMap, nonce uint64, vault *[20]byte, expiresAfter *uint64, source Source) (signer.Signature, error) {
	if s == nil {
		return signer.Signature{}, ErrNoSigner
	}
	hash, err := ActionHash(action, nonce, vault, expiresAfter)
	if err != nil {
		return signer.Signature{}, err
	}
	domain := signer.Domain{
		Name:              l1DomainName,
		Version:           l1DomainVersion,
		ChainID:           l1DomainChain,
		VerifyingContract: zeroAddress,
	}
	return s.SignTypedData(ctx, domain, l1AgentTypes, "Agent", map[string]any{
		"source":       string(source),
		"connectionId": hash,
	})
}

func BuildUserSignature(ctx context.Context, s signer.Signer, primaryType string, fields []signer.Field, message map[string]any, signatureChainID uint64) (signer.Signature, error) {
	if s == nil {
		return signer.Signature{}, ErrNoSigner
	}
	types := signer.Types{
		"EIP712Domain": l1AgentTypes["EIP712Domain"],
		primaryType:    fields,
	}
	domain := signer.Domain{
		Name:              userSignDomainName,
		Version:           userSignDomainVersion,
		ChainID:           signatureChainID,
		VerifyingContract: zeroAddress,
	}
	return s.SignTypedData(ctx, domain, types, primaryType, message)
}

// ActionHash returns the L1 action hash used by Hyperliquid's signing scheme.
//
//	keccak256( msgpack(action) || nonce_be8
//	           || (vault ? 0x01 || addr20 : 0x00)
//	           || (expiresAfter ? 0x00 || expiresAfter_be8 : ε) )
//
// The expiresAfter trailer is OMITTED entirely when nil — not zero-padded —
// matching the TS SDK's createL1ActionHash helper.
func ActionHash(action *msgpack.OrderedMap, nonce uint64, vault *[20]byte, expiresAfter *uint64) ([]byte, error) {
	encoded, err := msgpack.Encode(action)
	if err != nil {
		return nil, err
	}

	var nonceBytes [8]byte
	binary.BigEndian.PutUint64(nonceBytes[:], nonce)

	h := sha3.NewLegacyKeccak256()
	h.Write(encoded)
	h.Write(nonceBytes[:])
	if vault == nil {
		h.Write([]byte{0x00})
	} else {
		h.Write([]byte{0x01})
		h.Write(vault[:])
	}
	if expiresAfter != nil {
		var ea [8]byte
		binary.BigEndian.PutUint64(ea[:], *expiresAfter)
		h.Write([]byte{0x00})
		h.Write(ea[:])
	}
	return h.Sum(nil), nil
}

// submitL1 signs and submits an L1 action to /exchange, decoding the inner
// response body into out (skipped if out is nil).
func (c *Client) submitL1(ctx context.Context, action *msgpack.OrderedMap, out any) error {
	if c.Signer == nil {
		return ErrNoSigner
	}
	nonce := c.nonces.next()
	sig, err := BuildL1Signature(ctx, c.Signer, action, nonce, c.VaultAddress, nil, c.Source)
	if err != nil {
		return err
	}
	return c.send(ctx, action, nonce, sig, out)
}

// submitUser signs and submits a user-signed action. message must include
// "time" set to the nonce (Hyperliquid embeds the nonce inside the action).
func (c *Client) submitUser(ctx context.Context, action *msgpack.OrderedMap, primaryType string, fields []signer.Field, message map[string]any, out any) error {
	if c.Signer == nil {
		return ErrNoSigner
	}
	t, ok := message["time"].(uint64)
	if !ok || t == 0 {
		t = c.nonces.next()
		message["time"] = t
	}
	sig, err := BuildUserSignature(ctx, c.Signer, primaryType, fields, message, c.SignatureChainID)
	if err != nil {
		return err
	}
	return c.send(ctx, action, t, sig, out)
}

// send POSTs the signed envelope to /exchange and decodes the inner response.
// HTTP-level errors propagate from the transport. Successful HTTP with a
// status:"err" body returns *ActionRejected.
func (c *Client) send(ctx context.Context, action *msgpack.OrderedMap, nonce uint64, sig signer.Signature, out any) error {
	body := map[string]any{
		"action": msgpackActionAsJSON(action),
		"nonce":  nonce,
		"signature": map[string]any{
			"r": "0x" + hex.EncodeToString(sig.R[:]),
			"s": "0x" + hex.EncodeToString(sig.S[:]),
			"v": int(sig.V),
		},
	}
	if c.VaultAddress != nil {
		body["vaultAddress"] = "0x" + hex.EncodeToString(c.VaultAddress[:])
	}
	var raw StatusResponse
	if err := c.HTTP.PostJSON(ctx, "/exchange", body, &raw); err != nil {
		return err
	}
	if raw.Status != "ok" {
		return &ActionRejected{Action: actionType(action), Response: string(raw.Response)}
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(raw.Response, out)
}

func actionType(m *msgpack.OrderedMap) string {
	v, _ := m.Get("type")
	s, _ := v.(string)
	return s
}

// msgpackActionAsJSON renders an OrderedMap as a json.RawMessage preserving
// insertion order. The HTTP layer json.Marshal would otherwise re-key-sort
// our action.
func msgpackActionAsJSON(m *msgpack.OrderedMap) json.RawMessage {
	b, _ := msgpack.MarshalOrderedJSON(m)
	return b
}
