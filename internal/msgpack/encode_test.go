package msgpack

import (
	"bytes"
	"testing"
)

func TestEncode_PositiveFixint(t *testing.T) {
	got, err := Encode(uint64(0))
	if err != nil || !bytes.Equal(got, []byte{0x00}) {
		t.Fatalf("0 => %x err=%v", got, err)
	}
	got, _ = Encode(uint64(127))
	if !bytes.Equal(got, []byte{0x7f}) {
		t.Fatalf("127 => %x", got)
	}
}

func TestEncode_Int64(t *testing.T) {
	got, _ := Encode(int64(-1))
	if !bytes.Equal(got, []byte{0xff}) {
		t.Fatalf("-1 => %x", got)
	}
	got, _ = Encode(int64(-32))
	if !bytes.Equal(got, []byte{0xe0}) {
		t.Fatalf("-32 => %x", got)
	}
}

func TestEncode_String(t *testing.T) {
	got, _ := Encode("hi")
	if !bytes.Equal(got, []byte{0xa2, 'h', 'i'}) {
		t.Fatalf("'hi' => %x", got)
	}
}

func TestEncode_Bool(t *testing.T) {
	tt, _ := Encode(true)
	ff, _ := Encode(false)
	if tt[0] != 0xc3 || ff[0] != 0xc2 {
		t.Fatalf("bools: %x %x", tt, ff)
	}
}

func TestEncode_MapPreservesOrder(t *testing.T) {
	m := NewOrderedMap()
	m.Set("type", "order")
	m.Set("orders", []any{})
	m.Set("grouping", "na")

	got, err := Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] != 0x83 {
		t.Fatalf("first byte = %02x", got[0])
	}
	wantPrefix := []byte{0x83, 0xa4, 't', 'y', 'p', 'e', 0xa5, 'o', 'r', 'd', 'e', 'r'}
	if !bytes.HasPrefix(got, wantPrefix) {
		t.Fatalf("got prefix %x", got[:len(wantPrefix)])
	}
}
