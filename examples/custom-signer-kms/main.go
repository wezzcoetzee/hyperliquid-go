// Example skeleton: implement signer.Signer with a remote KMS / HSM.
//
// This file does NOT compile against any specific KMS SDK — it shows the
// shape of the interface and the responsibilities of a custom signer.
// See the comments in SignTypedData for how to wire it to AWS KMS, GCP
// Cloud KMS, HashiCorp Vault, Ledger, or any other backend.
package main

import (
	"context"
	"fmt"

	"github.com/wezzcoetzee/hyperliquid-go/signer"
)

type kmsSigner struct {
	keyID string
	addr  [20]byte
}

func (k *kmsSigner) Address() [20]byte { return k.addr }

func (k *kmsSigner) SignTypedData(ctx context.Context, d signer.Domain, t signer.Types, primary string, msg map[string]any) (signer.Signature, error) {
	// 1. Compute the EIP-712 typed-data digest:
	//
	//      hash, err := eip712.HashTypedData(eipDomain, eipTypes, primary, msg)
	//
	//    (Mirror the type-mapping done in signer/privkey/privkey.go.)
	//
	// 2. Send the 32-byte hash to your KMS for raw secp256k1 signing.
	//    For AWS KMS use KeySpec=ECC_SECG_P256K1, MessageType=DIGEST.
	//
	// 3. Parse the DER-encoded signature into r and s.
	// 4. Recover v by trying both candidates (0 and 1) against the known
	//    address. Return r/s/v with v normalized to 27/28.
	return signer.Signature{}, fmt.Errorf("kmsSigner.SignTypedData: not implemented; see comments")
}

func main() {
	_ = &kmsSigner{}
	fmt.Println("This is a skeleton. Wire to your KMS before using.")
}
