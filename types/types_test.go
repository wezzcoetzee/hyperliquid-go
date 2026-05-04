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

func TestTif_WireValues(t *testing.T) {
	cases := map[Tif]string{
		TifGtc: "Gtc",
		TifIoc: "Ioc",
		TifAlo: "Alo",
	}
	for tif, want := range cases {
		if string(tif) != want {
			t.Errorf("Tif %v = %q, want %q", tif, string(tif), want)
		}
	}
}

func TestParseCloid_Valid(t *testing.T) {
	c, err := ParseCloid("0x0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	if string(c) != "0x0123456789abcdef0123456789abcdef" {
		t.Errorf("Cloid = %q", c)
	}
	if got := c.Bytes(); len(got) != 16 {
		t.Errorf("Bytes() length = %d, want 16", len(got))
	}
}

func TestParseCloid_Invalid(t *testing.T) {
	cases := []string{
		"",
		"0x",
		"0x0123",                                 // too short
		"0123456789abcdef0123456789abcdef",       // missing 0x
		"0xZZ23456789abcdef0123456789abcdef",     // non-hex
		"0x0123456789abcdef0123456789abcdef00",   // too long
	}
	for _, s := range cases {
		if _, err := ParseCloid(s); err == nil {
			t.Errorf("ParseCloid(%q) expected error", s)
		}
	}
}

func TestCloid_BytesOnInvalid(t *testing.T) {
	if got := Cloid("garbage").Bytes(); got != nil {
		t.Errorf("expected nil for invalid Cloid, got %x", got)
	}
}
