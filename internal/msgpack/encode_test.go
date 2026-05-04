package msgpack

import (
	"bytes"
	"testing"
)

func TestEncode_Nil(t *testing.T) {
	got, err := Encode(nil)
	if err != nil || !bytes.Equal(got, []byte{0xc0}) {
		t.Fatalf("nil => %x err=%v", got, err)
	}
}

func TestEncode_Bool(t *testing.T) {
	tt, _ := Encode(true)
	ff, _ := Encode(false)
	if !bytes.Equal(tt, []byte{0xc3}) || !bytes.Equal(ff, []byte{0xc2}) {
		t.Fatalf("bools: %x %x", tt, ff)
	}
}

func TestEncode_UintBoundaries(t *testing.T) {
	cases := []struct {
		in   uint64
		want []byte
	}{
		{0, []byte{0x00}},
		{127, []byte{0x7f}},
		{128, []byte{0xcc, 0x80}},   // smallest uint8
		{255, []byte{0xcc, 0xff}},
		{256, []byte{0xcd, 0x01, 0x00}}, // smallest uint16
		{65535, []byte{0xcd, 0xff, 0xff}},
		{65536, []byte{0xce, 0x00, 0x01, 0x00, 0x00}}, // smallest uint32
		{4294967295, []byte{0xce, 0xff, 0xff, 0xff, 0xff}},
		{4294967296, []byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}}, // smallest uint64
	}
	for _, c := range cases {
		got, err := Encode(c.in)
		if err != nil {
			t.Fatalf("Encode(%d): %v", c.in, err)
		}
		if !bytes.Equal(got, c.want) {
			t.Errorf("Encode(%d) = %x, want %x", c.in, got, c.want)
		}
	}
}

func TestEncode_IntNegativeBoundaries(t *testing.T) {
	cases := []struct {
		in   int64
		want []byte
	}{
		{-1, []byte{0xff}},
		{-32, []byte{0xe0}},                                      // smallest fixint
		{-33, []byte{0xd0, 0xdf}},                                // int8
		{-128, []byte{0xd0, 0x80}},
		{-129, []byte{0xd1, 0xff, 0x7f}},                         // int16
		{-32768, []byte{0xd1, 0x80, 0x00}},
		{-32769, []byte{0xd2, 0xff, 0xff, 0x7f, 0xff}},           // int32
		{-(1 << 31), []byte{0xd2, 0x80, 0x00, 0x00, 0x00}},
		{-(1 << 31) - 1, []byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff}}, // int64
	}
	for _, c := range cases {
		got, err := Encode(c.in)
		if err != nil {
			t.Fatalf("Encode(%d): %v", c.in, err)
		}
		if !bytes.Equal(got, c.want) {
			t.Errorf("Encode(%d) = %x, want %x", c.in, got, c.want)
		}
	}
}

func TestEncode_IntPositiveDelegatesToUint(t *testing.T) {
	got, _ := Encode(int64(1000))
	want, _ := Encode(uint64(1000))
	if !bytes.Equal(got, want) {
		t.Errorf("int64(1000) %x != uint64(1000) %x", got, want)
	}
}

func TestEncode_StringBoundaries(t *testing.T) {
	mk := func(n int) string {
		s := make([]byte, n)
		for i := range s {
			s[i] = 'a'
		}
		return string(s)
	}

	cases := []struct {
		name      string
		in        string
		wantStart []byte // header bytes
		wantLen   int
	}{
		{"fixstr 0", "", []byte{0xa0}, 1},
		{"fixstr 31", mk(31), []byte{0xbf}, 32},
		{"str8 32", mk(32), []byte{0xd9, 32}, 34},
		{"str8 255", mk(255), []byte{0xd9, 0xff}, 257},
		{"str16 256", mk(256), []byte{0xda, 0x01, 0x00}, 259},
		{"str16 65535", mk(65535), []byte{0xda, 0xff, 0xff}, 65538},
		{"str32 65536", mk(65536), []byte{0xdb, 0x00, 0x01, 0x00, 0x00}, 65541},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Encode(c.in)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != c.wantLen {
				t.Errorf("len = %d, want %d", len(got), c.wantLen)
			}
			if !bytes.HasPrefix(got, c.wantStart) {
				t.Errorf("prefix = %x, want %x", got[:len(c.wantStart)], c.wantStart)
			}
		})
	}
}

func TestEncode_BinBoundaries(t *testing.T) {
	mk := func(n int) []byte { return bytes.Repeat([]byte{0xab}, n) }

	cases := []struct {
		name      string
		in        []byte
		wantStart []byte
		wantLen   int
	}{
		{"bin8 0", mk(0), []byte{0xc4, 0x00}, 2},
		{"bin8 255", mk(255), []byte{0xc4, 0xff}, 257},
		{"bin16 256", mk(256), []byte{0xc5, 0x01, 0x00}, 259},
		{"bin16 65535", mk(65535), []byte{0xc5, 0xff, 0xff}, 65538},
		{"bin32 65536", mk(65536), []byte{0xc6, 0x00, 0x01, 0x00, 0x00}, 65541},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Encode(c.in)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != c.wantLen {
				t.Errorf("len = %d, want %d", len(got), c.wantLen)
			}
			if !bytes.HasPrefix(got, c.wantStart) {
				t.Errorf("prefix = %x, want %x", got[:len(c.wantStart)], c.wantStart)
			}
		})
	}
}

func TestEncode_ArrayBoundaries(t *testing.T) {
	mk := func(n int) []any {
		out := make([]any, n)
		for i := range out {
			out[i] = uint64(0)
		}
		return out
	}

	got, _ := Encode(mk(15))
	if got[0] != 0x9f {
		t.Errorf("fixarray 15: header %02x", got[0])
	}
	got, _ = Encode(mk(16))
	if !bytes.HasPrefix(got, []byte{0xdc, 0x00, 0x10}) {
		t.Errorf("array16 16: prefix %x", got[:3])
	}
}

func TestEncode_MapHeaderBoundary(t *testing.T) {
	m := NewOrderedMap()
	for i := 0; i < 16; i++ {
		k := string(rune('a' + i))
		m.Set(k, uint64(0))
	}
	got, err := Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(got, []byte{0xde, 0x00, 0x10}) {
		t.Errorf("map16 boundary: prefix %x", got[:3])
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
	wantPrefix := []byte{0x83, 0xa4, 't', 'y', 'p', 'e', 0xa5, 'o', 'r', 'd', 'e', 'r'}
	if !bytes.HasPrefix(got, wantPrefix) {
		t.Fatalf("got prefix %x", got[:len(wantPrefix)])
	}
}

func TestOrderedMap_OverwriteKeepsPosition(t *testing.T) {
	m := NewOrderedMap()
	m.Set("a", uint64(1))
	m.Set("b", uint64(2))
	m.Set("a", uint64(99)) // overwrite
	keys := m.Keys()
	if len(keys) != 2 || keys[0] != "a" || keys[1] != "b" {
		t.Fatalf("keys = %v", keys)
	}
	v, _ := m.Get("a")
	if v.(uint64) != 99 {
		t.Errorf("a = %v", v)
	}
}

func TestEncode_RejectsFloat(t *testing.T) {
	if _, err := Encode(float64(1.5)); err == nil {
		t.Fatal("expected error for float64")
	}
}

func TestEncode_RejectsMap(t *testing.T) {
	if _, err := Encode(map[string]any{"x": uint64(1)}); err == nil {
		t.Fatal("expected error for map[string]any (use OrderedMap)")
	}
}

func TestEncode_RejectsUnsupportedInts(t *testing.T) {
	if _, err := Encode(int32(1)); err == nil {
		t.Fatal("expected error for int32")
	}
	if _, err := Encode(uint32(1)); err == nil {
		t.Fatal("expected error for uint32")
	}
}
