package privkey

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/signer"
)

func TestPrivKey_AddressDerivation(t *testing.T) {
	s, err := New("0x0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	addr := s.Address()
	got := hex.EncodeToString(addr[:])
	want := "7e5f4552091a69125d5dfcb7b8c2659029395bdf"
	if got != want {
		t.Fatalf("addr = %s want %s", got, want)
	}
}

func TestPrivKey_New_RejectsBadKey(t *testing.T) {
	cases := []string{
		"",
		"0x",
		"0x123",
		"not-hex",
		"0xZZZZ000000000000000000000000000000000000000000000000000000000000",
	}
	for _, c := range cases {
		if _, err := New(c); err == nil {
			t.Errorf("New(%q) expected error", c)
		}
	}
}

func TestPrivKey_SignsAgent(t *testing.T) {
	s, _ := New("0x0000000000000000000000000000000000000000000000000000000000000001")
	domain := signer.Domain{
		Name:              "Exchange",
		Version:           "1",
		ChainID:           1337,
		VerifyingContract: "0x0000000000000000000000000000000000000000",
	}
	types := signer.Types{
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
	cid := make([]byte, 32)
	cid[0] = 0xab
	sig, err := s.SignTypedData(context.Background(), domain, types, "Agent", map[string]any{
		"source":       "a",
		"connectionId": cid,
	})
	if err != nil {
		t.Fatal(err)
	}
	if sig.V != 27 && sig.V != 28 {
		t.Fatalf("v = %d", sig.V)
	}
	var zero [32]byte
	if sig.R == zero {
		t.Error("R is zero")
	}
	if sig.S == zero {
		t.Error("S is zero")
	}
}

func TestPrivKey_SignsConsistently(t *testing.T) {
	s, _ := New("0x0000000000000000000000000000000000000000000000000000000000000001")
	domain := signer.Domain{Name: "Exchange", Version: "1", ChainID: 1337, VerifyingContract: "0x0000000000000000000000000000000000000000"}
	types := signer.Types{
		"EIP712Domain": {
			{Name: "name", Type: "string"},
			{Name: "version", Type: "string"},
			{Name: "chainId", Type: "uint256"},
			{Name: "verifyingContract", Type: "address"},
		},
		"Agent": {{Name: "source", Type: "string"}, {Name: "connectionId", Type: "bytes32"}},
	}
	msg := map[string]any{"source": "a", "connectionId": make([]byte, 32)}

	a, err := s.SignTypedData(context.Background(), domain, types, "Agent", msg)
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.SignTypedData(context.Background(), domain, types, "Agent", msg)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatal("signature is not deterministic")
	}
}
