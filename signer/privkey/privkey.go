// Package privkey provides a Signer backed by a secp256k1 private key, using
// go-ethereum's crypto package for ECDSA. Use this when your application
// holds the key directly. For HSM/KMS/hardware wallets, implement
// signer.Signer with the appropriate backend.
package privkey

import (
	"context"
	"errors"
	"strings"

	"crypto/ecdsa"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/wezzcoetzee/hyperliquid-go/internal/eip712"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
)

// PrivKey is a signer.Signer backed by a raw secp256k1 private key.
// Use New to construct; the zero value is invalid.
type PrivKey struct {
	key  *ecdsa.PrivateKey
	addr [20]byte
}

// New parses a 64-character hex-encoded secp256k1 private key (with or without
// a leading "0x") and returns a ready-to-use PrivKey. The derived Ethereum
// address is cached on construction.
func New(hexKey string) (*PrivKey, error) {
	hexKey = strings.TrimPrefix(strings.ToLower(strings.TrimSpace(hexKey)), "0x")
	if len(hexKey) != 64 {
		return nil, errors.New("privkey: expected 64-character hex private key")
	}
	k, err := ethcrypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
	}
	pk := &PrivKey{key: k}
	addr := ethcrypto.PubkeyToAddress(k.PublicKey)
	copy(pk.addr[:], addr[:])
	return pk, nil
}

// Address returns the 20-byte Ethereum address derived from the private key.
func (p *PrivKey) Address() [20]byte { return p.addr }

// SignTypedData hashes the EIP-712 typed data and signs it with the private key,
// returning a 65-byte r/s/v signature with v adjusted to 27 or 28.
func (p *PrivKey) SignTypedData(ctx context.Context, d signer.Domain, t signer.Types, primary string, msg map[string]any) (signer.Signature, error) {
	domain := eip712.Domain{
		Name:              d.Name,
		Version:           d.Version,
		ChainID:           d.ChainID,
		VerifyingContract: d.VerifyingContract,
	}
	types := make(eip712.Types, len(t))
	for k, fs := range t {
		ff := make([]eip712.Field, len(fs))
		for i, f := range fs {
			ff[i] = eip712.Field{Name: f.Name, Type: f.Type}
		}
		types[k] = ff
	}
	hash, err := eip712.HashTypedData(domain, types, primary, msg)
	if err != nil {
		return signer.Signature{}, err
	}
	sigBytes, err := ethcrypto.Sign(hash, p.key)
	if err != nil {
		return signer.Signature{}, err
	}
	if len(sigBytes) != 65 {
		return signer.Signature{}, errors.New("privkey: unexpected signature length")
	}
	var sig signer.Signature
	copy(sig.R[:], sigBytes[0:32])
	copy(sig.S[:], sigBytes[32:64])
	sig.V = sigBytes[64] + 27
	return sig, nil
}
