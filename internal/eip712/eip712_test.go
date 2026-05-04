package eip712

import (
	"encoding/hex"
	"testing"
)

// Spec vector from EIP-712. Hashing the "Mail" message must produce the well-known
// digest 0xbe609aee343fb3c4b28e1df9e632fca64fcfaede20f02e86244efddf30957bd2.
func TestHashTypedData_SpecMail(t *testing.T) {
	domain := Domain{
		Name:              "Ether Mail",
		Version:           "1",
		ChainID:           1,
		VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
	}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Person": {{Name: "name", Type: "string"}, {Name: "wallet", Type: "address"}},
		"Mail": {
			{Name: "from", Type: "Person"},
			{Name: "to", Type: "Person"},
			{Name: "contents", Type: "string"},
		},
	}
	msg := map[string]any{
		"from":     map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
		"to":       map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
		"contents": "Hello, Bob!",
	}

	h, err := HashTypedData(domain, types, "Mail", msg)
	if err != nil {
		t.Fatal(err)
	}
	want := "be609aee343fb3c4b28e1df9e632fca64fcfaede20f02e86244efddf30957bd2"
	if hex.EncodeToString(h) != want {
		t.Fatalf("got %s want %s", hex.EncodeToString(h), want)
	}
}
