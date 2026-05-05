package exchange

import (
	"context"
	"encoding/binary"

	"golang.org/x/crypto/sha3"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid/signer"
)

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

func BuildL1Signature(ctx context.Context, s signer.Signer, action *msgpack.OrderedMap, nonce uint64, vault *[20]byte, expiresAfter *uint64, mainnet bool) (signer.Signature, error) {
	hash, err := ActionHash(action, nonce, vault, expiresAfter)
	if err != nil {
		return signer.Signature{}, err
	}
	source := "b"
	if mainnet {
		source = "a"
	}
	domain := signer.Domain{
		Name:              l1DomainName,
		Version:           l1DomainVersion,
		ChainID:           l1DomainChain,
		VerifyingContract: zeroAddress,
	}
	return s.SignTypedData(ctx, domain, l1AgentTypes, "Agent", map[string]any{
		"source":       source,
		"connectionId": hash,
	})
}

func BuildUserSignature(ctx context.Context, s signer.Signer, primaryType string, fields []signer.Field, message map[string]any, signatureChainID uint64) (signer.Signature, error) {
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
