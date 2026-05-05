// Package signer defines the EIP-712 signing interface used by Hyperliquid
// Exchange writes. Implementations sign typed data using the private key
// (or HSM/KMS/wallet) of choice. A default secp256k1 implementation lives
// at signer/privkey.
package signer

import "context"

// Domain mirrors the EIP-712 EIP712Domain struct fields actually used by
// Hyperliquid. Salt is unused.
type Domain struct {
	Name              string
	Version           string
	ChainID           uint64
	VerifyingContract string // 0x-prefixed hex; case is normalized by the EIP-712 layer
}

// Field describes one entry in an EIP-712 type definition.
type Field struct {
	Name string
	Type string
}

// Types maps an EIP-712 primary type name to its ordered field list.
// Implementations MUST include the "EIP712Domain" entry the signer expects.
type Types map[string][]Field

// Signature is the canonical 65-byte ECDSA signature split into r/s/v.
// v is 27 or 28 to match Ethereum/EIP-155 normalization used by viem and
// go-ethereum's crypto.Sign output (after the +27 adjustment).
type Signature struct {
	R [32]byte
	S [32]byte
	V byte
}

// Signer abstracts EIP-712 signing. Address returns the 20-byte address
// associated with the signing key. SignTypedData hashes the typed data per
// EIP-712 and returns r/s/v.
type Signer interface {
	Address() [20]byte
	SignTypedData(ctx context.Context, domain Domain, types Types, primaryType string, message map[string]any) (Signature, error)
}
