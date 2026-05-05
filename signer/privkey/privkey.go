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

	"github.com/wezzcoetzee/hyperliquid/internal/eip712"
	"github.com/wezzcoetzee/hyperliquid/signer"
)

type PrivKey struct {
	key  *ecdsa.PrivateKey
	addr [20]byte
}

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

func (p *PrivKey) Address() [20]byte { return p.addr }

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
