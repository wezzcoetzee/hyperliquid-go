// Package msgpack implements a minimal deterministic msgpack encoder for
// Hyperliquid action shapes. It is NOT a general-purpose msgpack library.
//
// Supported value types: nil, bool, string, []byte, int, int64, uint64,
// []any, *OrderedMap. Any other type (including float64, int32, uint32,
// uint8, map[string]any) returns an error.
//
// All output uses the smallest representation per the msgpack spec, with
// big-endian multi-byte integers. Maps are encoded in OrderedMap insertion
// order, which is what the TS SDK produces for object literals — using a
// regular Go map would give nondeterministic byte output and break signing
// parity, so map[string]any is deliberately not supported.
package msgpack

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// OrderedMap preserves insertion order so encoded output matches TS object
// iteration order.
type OrderedMap struct {
	keys   []string
	values map[string]any
}

func NewOrderedMap() *OrderedMap { return &OrderedMap{values: map[string]any{}} }

func (m *OrderedMap) Set(k string, v any) {
	if _, ok := m.values[k]; !ok {
		m.keys = append(m.keys, k)
	}
	m.values[k] = v
}

func (m *OrderedMap) Len() int       { return len(m.keys) }
func (m *OrderedMap) Keys() []string { return m.keys }

func (m *OrderedMap) Get(k string) (any, bool) {
	v, ok := m.values[k]
	return v, ok
}

func Encode(v any) ([]byte, error) {
	return encode(nil, v)
}

func encode(buf []byte, v any) ([]byte, error) {
	switch x := v.(type) {
	case nil:
		return append(buf, 0xc0), nil
	case bool:
		if x {
			return append(buf, 0xc3), nil
		}
		return append(buf, 0xc2), nil
	case string:
		return encodeString(buf, x), nil
	case []byte:
		return encodeBin(buf, x), nil
	case int:
		return encodeInt(buf, int64(x)), nil
	case int64:
		return encodeInt(buf, x), nil
	case uint64:
		return encodeUint(buf, x), nil
	case float64:
		return nil, errors.New("msgpack: float not supported (Hyperliquid uses string-encoded decimals)")
	case []any:
		buf = encodeArrayHeader(buf, len(x))
		var err error
		for _, e := range x {
			buf, err = encode(buf, e)
			if err != nil {
				return nil, err
			}
		}
		return buf, nil
	case *OrderedMap:
		buf = encodeMapHeader(buf, x.Len())
		var err error
		for _, k := range x.keys {
			buf = encodeString(buf, k)
			buf, err = encode(buf, x.values[k])
			if err != nil {
				return nil, err
			}
		}
		return buf, nil
	default:
		return nil, fmt.Errorf("msgpack: unsupported type %T", v)
	}
}

func encodeUint(buf []byte, n uint64) []byte {
	switch {
	case n <= 0x7f:
		return append(buf, byte(n))
	case n <= 0xff:
		return append(buf, 0xcc, byte(n))
	case n <= 0xffff:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(n))
		return append(append(buf, 0xcd), b[:]...)
	case n <= 0xffffffff:
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(n))
		return append(append(buf, 0xce), b[:]...)
	default:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], n)
		return append(append(buf, 0xcf), b[:]...)
	}
}

func encodeInt(buf []byte, n int64) []byte {
	if n >= 0 {
		return encodeUint(buf, uint64(n))
	}
	switch {
	case n >= -32:
		return append(buf, byte(n))
	case n >= -128:
		return append(buf, 0xd0, byte(int8(n)))
	case n >= -32768:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(int16(n)))
		return append(append(buf, 0xd1), b[:]...)
	case n >= -(1 << 31):
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(int32(n)))
		return append(append(buf, 0xd2), b[:]...)
	default:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(n))
		return append(append(buf, 0xd3), b[:]...)
	}
}

func encodeString(buf []byte, s string) []byte {
	n := len(s)
	switch {
	case n < 32:
		buf = append(buf, 0xa0|byte(n))
	case n < 256:
		buf = append(buf, 0xd9, byte(n))
	case n < 65536:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(n))
		buf = append(append(buf, 0xda), b[:]...)
	default:
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(n))
		buf = append(append(buf, 0xdb), b[:]...)
	}
	return append(buf, s...)
}

func encodeBin(buf []byte, b []byte) []byte {
	n := len(b)
	switch {
	case n < 256:
		buf = append(buf, 0xc4, byte(n))
	case n < 65536:
		var bb [2]byte
		binary.BigEndian.PutUint16(bb[:], uint16(n))
		buf = append(append(buf, 0xc5), bb[:]...)
	default:
		var bb [4]byte
		binary.BigEndian.PutUint32(bb[:], uint32(n))
		buf = append(append(buf, 0xc6), bb[:]...)
	}
	return append(buf, b...)
}

func encodeArrayHeader(buf []byte, n int) []byte {
	switch {
	case n < 16:
		return append(buf, 0x90|byte(n))
	case n < 65536:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(n))
		return append(append(buf, 0xdc), b[:]...)
	default:
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(n))
		return append(append(buf, 0xdd), b[:]...)
	}
}

func encodeMapHeader(buf []byte, n int) []byte {
	switch {
	case n < 16:
		return append(buf, 0x80|byte(n))
	case n < 65536:
		var b [2]byte
		binary.BigEndian.PutUint16(b[:], uint16(n))
		return append(append(buf, 0xde), b[:]...)
	default:
		var b [4]byte
		binary.BigEndian.PutUint32(b[:], uint32(n))
		return append(append(buf, 0xdf), b[:]...)
	}
}
