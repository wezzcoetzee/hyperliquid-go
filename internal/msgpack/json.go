package msgpack

import (
	"bytes"
	"encoding/json"
)

// MarshalOrderedJSON renders an OrderedMap (and nested values) to JSON,
// preserving key insertion order. Nested arrays and OrderedMaps are
// recursed; primitive values fall back to encoding/json.Marshal.
func MarshalOrderedJSON(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := writeJSON(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeJSON(buf *bytes.Buffer, v any) error {
	switch x := v.(type) {
	case *OrderedMap:
		buf.WriteByte('{')
		for i, k := range x.keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			kb, err := json.Marshal(k)
			if err != nil {
				return err
			}
			buf.Write(kb)
			buf.WriteByte(':')
			if err := writeJSON(buf, x.values[k]); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
		return nil
	case []any:
		buf.WriteByte('[')
		for i, e := range x {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeJSON(buf, e); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf.Write(b)
		return nil
	}
}
