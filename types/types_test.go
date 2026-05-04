package types

import "testing"

func TestAddress_Normalize(t *testing.T) {
	a, err := ParseAddress("0xABCDEF0123456789abcdef0123456789abcdef01")
	if err != nil {
		t.Fatal(err)
	}
	if a.Hex() != "0xabcdef0123456789abcdef0123456789abcdef01" {
		t.Errorf("hex = %q", a.Hex())
	}
}

func TestAddress_Invalid(t *testing.T) {
	if _, err := ParseAddress("not-an-address"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSide_String(t *testing.T) {
	if Buy.String() != "B" || Sell.String() != "A" {
		t.Errorf("buy=%s sell=%s", Buy, Sell)
	}
}
