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

func TestEncode_BytesN_Padding(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "data", Type: "bytes32"}},
	}
	data := make([]byte, 32)
	data[0] = 0xab
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"data": data}); err != nil {
		t.Fatalf("32-byte: %v", err)
	}
}

func TestEncode_BytesN_RejectsWrongLength(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "data", Type: "bytes32"}},
	}
	short := make([]byte, 5)
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"data": short}); err == nil {
		t.Fatal("expected error for 5-byte bytes32")
	}
}

func TestEncode_BytesN_RejectsInvalidType(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "data", Type: "bytes33"}},
	}
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"data": make([]byte, 33)}); err == nil {
		t.Fatal("expected error for bytes33")
	}
}

func TestEncode_Address_RejectsMalformed(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "addr", Type: "address"}},
	}
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"addr": "0x1234"}); err == nil {
		t.Fatal("expected error for short address")
	}
}

func TestEncode_Address_MixedCase(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "addr", Type: "address"}},
	}
	hLower, _ := HashTypedData(domain, types, "Wrapper", map[string]any{"addr": "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"})
	hUpper, _ := HashTypedData(domain, types, "Wrapper", map[string]any{"addr": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"})
	if string(hLower) != string(hUpper) {
		t.Fatal("mixed-case and lowercase addresses must hash the same")
	}
}

func TestEncode_RejectsNegativeInt(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "n", Type: "int256"}},
	}
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"n": int64(-1)}); err == nil {
		t.Fatal("expected error for negative int")
	}
}

func TestEncode_MissingField(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "x", Type: "string"}, {Name: "y", Type: "string"}},
	}
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"x": "ok"}); err == nil {
		t.Fatal("expected error for missing y")
	}
}

func TestEncode_UnknownType(t *testing.T) {
	domain := Domain{Name: "x", Version: "1", ChainID: 1, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Wrapper": {{Name: "v", Type: "weirdType"}},
	}
	if _, err := HashTypedData(domain, types, "Wrapper", map[string]any{"v": "x"}); err == nil {
		t.Fatal("expected error for unknown type")
	}
}
