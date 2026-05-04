// Package eip712 implements the typed-data hashing scheme of EIP-712.
//
// Scope: it covers the value types Hyperliquid signing needs — string, bytes,
// fixed-size byte arrays bytes1..bytes32 (length-validated), address (length-
// validated to 20 bytes), bool, unsigned integer types (uint*/int* — negative
// values are rejected), nested structs, and homogeneous arrays. Salts are
// supported via Domain.Salt routed through the bytes32 field.
//
// HashTypedData(...) of the canonical "Mail" example matches the spec digest.
package eip712

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"golang.org/x/crypto/sha3"
)

type Domain struct {
	Name              string
	Version           string
	ChainID           uint64
	VerifyingContract string // hex address; "" omits this field if not in EIP712Domain types
	Salt              []byte // optional 32-byte salt
}

type Field struct{ Name, Type string }
type Types map[string][]Field

func keccak(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}

// HashTypedData returns keccak256(0x19 0x01 || domainSeparator || hashStruct(message)).
func HashTypedData(domain Domain, types Types, primaryType string, message map[string]any) ([]byte, error) {
	domainSeparator, err := hashStruct("EIP712Domain", types, domainAsMap(domain, types))
	if err != nil {
		return nil, fmt.Errorf("domain hash: %w", err)
	}
	msgHash, err := hashStruct(primaryType, types, message)
	if err != nil {
		return nil, fmt.Errorf("message hash: %w", err)
	}
	return keccak([]byte{0x19, 0x01}, domainSeparator, msgHash), nil
}

func domainAsMap(d Domain, types Types) map[string]any {
	m := map[string]any{}
	for _, f := range types["EIP712Domain"] {
		switch f.Name {
		case "name":
			m["name"] = d.Name
		case "version":
			m["version"] = d.Version
		case "chainId":
			m["chainId"] = new(big.Int).SetUint64(d.ChainID)
		case "verifyingContract":
			m["verifyingContract"] = d.VerifyingContract
		case "salt":
			m["salt"] = d.Salt
		}
	}
	return m
}

func hashStruct(primary string, types Types, data map[string]any) ([]byte, error) {
	enc, err := encodeType(primary, types)
	if err != nil {
		return nil, err
	}
	typeHash := keccak([]byte(enc))
	encoded, err := encodeData(primary, types, data)
	if err != nil {
		return nil, err
	}
	return keccak(typeHash, encoded), nil
}

func encodeType(primary string, types Types) (string, error) {
	if _, ok := types[primary]; !ok {
		return "", fmt.Errorf("unknown primary type %q", primary)
	}
	deps := map[string]bool{}
	collectDeps(primary, types, deps)
	delete(deps, primary)
	sorted := make([]string, 0, len(deps))
	for d := range deps {
		sorted = append(sorted, d)
	}
	sort.Strings(sorted)
	all := append([]string{primary}, sorted...)

	var b strings.Builder
	for _, name := range all {
		fields := types[name]
		b.WriteString(name)
		b.WriteByte('(')
		for i, f := range fields {
			if i > 0 {
				b.WriteByte(',')
			}
			b.WriteString(f.Type)
			b.WriteByte(' ')
			b.WriteString(f.Name)
		}
		b.WriteByte(')')
	}
	return b.String(), nil
}

func collectDeps(name string, types Types, out map[string]bool) {
	if out[name] {
		return
	}
	if _, ok := types[name]; !ok {
		return
	}
	out[name] = true
	for _, f := range types[name] {
		base := strings.TrimSuffix(f.Type, "[]")
		if _, isStruct := types[base]; isStruct {
			collectDeps(base, types, out)
		}
	}
}

func encodeData(primary string, types Types, data map[string]any) ([]byte, error) {
	fields := types[primary]
	out := make([]byte, 0, 32*len(fields))
	for _, f := range fields {
		v, ok := data[f.Name]
		if !ok {
			return nil, fmt.Errorf("missing field %q in %s", f.Name, primary)
		}
		enc, err := encodeValue(f.Type, v, types)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", f.Name, err)
		}
		out = append(out, enc...)
	}
	return out, nil
}

func encodeValue(t string, v any, types Types) ([]byte, error) {
	switch {
	case t == "string":
		s, ok := v.(string)
		if !ok {
			return nil, errors.New("expected string")
		}
		return keccak([]byte(s)), nil
	case t == "bytes":
		b, ok := v.([]byte)
		if !ok {
			return nil, errors.New("expected []byte for bytes")
		}
		return keccak(b), nil
	case strings.HasPrefix(t, "bytes"):
		size, err := parseBytesN(t)
		if err != nil {
			return nil, err
		}
		b, ok := v.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected []byte for %s", t)
		}
		if len(b) != size {
			return nil, fmt.Errorf("eip712: %s expects %d bytes, got %d", t, size, len(b))
		}
		return leftPad32(b), nil
	case t == "address":
		s, ok := v.(string)
		if !ok {
			return nil, errors.New("expected hex string for address")
		}
		s = strings.TrimPrefix(strings.ToLower(s), "0x")
		b, err := hexDecode(s)
		if err != nil {
			return nil, err
		}
		if len(b) != 20 {
			return nil, fmt.Errorf("eip712: address must be 20 bytes, got %d", len(b))
		}
		return leftPad32(b), nil
	case t == "bool":
		b, ok := v.(bool)
		if !ok {
			return nil, errors.New("expected bool")
		}
		out := make([]byte, 32)
		if b {
			out[31] = 1
		}
		return out, nil
	case strings.HasPrefix(t, "uint") || strings.HasPrefix(t, "int"):
		var n *big.Int
		switch x := v.(type) {
		case *big.Int:
			n = x
		case uint64:
			n = new(big.Int).SetUint64(x)
		case int64:
			n = big.NewInt(x)
		case int:
			n = big.NewInt(int64(x))
		default:
			return nil, fmt.Errorf("unsupported int type %T", v)
		}
		if n.Sign() < 0 {
			return nil, fmt.Errorf("eip712: negative %s not supported (Hyperliquid uses unsigned values)", t)
		}
		out := make([]byte, 32)
		nb := n.Bytes()
		copy(out[32-len(nb):], nb)
		return out, nil
	case strings.HasSuffix(t, "[]"):
		base := strings.TrimSuffix(t, "[]")
		arr, ok := v.([]any)
		if !ok {
			return nil, errors.New("expected []any")
		}
		var concat []byte
		for _, e := range arr {
			enc, err := encodeValue(base, e, types)
			if err != nil {
				return nil, err
			}
			concat = append(concat, enc...)
		}
		return keccak(concat), nil
	default:
		if _, ok := types[t]; ok {
			m, ok := v.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("expected map for struct %s", t)
			}
			return hashStruct(t, types, m)
		}
		return nil, fmt.Errorf("unknown type %q", t)
	}
}

// parseBytesN parses a "bytesN" type string (1 ≤ N ≤ 32) into N. It rejects
// inputs like "bytesfoo" or "bytes33" so they fall through to the unknown-type
// error path instead of silently truncating.
func parseBytesN(t string) (int, error) {
	if len(t) <= len("bytes") {
		return 0, fmt.Errorf("eip712: invalid type %q", t)
	}
	rest := t[len("bytes"):]
	n := 0
	for _, c := range rest {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("eip712: invalid type %q", t)
		}
		n = n*10 + int(c-'0')
	}
	if n < 1 || n > 32 {
		return 0, fmt.Errorf("eip712: invalid type %q (N must be 1..32)", t)
	}
	return n, nil
}

func leftPad32(b []byte) []byte {
	if len(b) >= 32 {
		return b[:32]
	}
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}

func hexDecode(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		s = "0" + s
	}
	out := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		hi, lo := fromHex(s[i]), fromHex(s[i+1])
		if hi < 0 || lo < 0 {
			return nil, fmt.Errorf("invalid hex")
		}
		out[i/2] = byte(hi<<4 | lo)
	}
	return out, nil
}

func fromHex(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c - 'a' + 10)
	case c >= 'A' && c <= 'F':
		return int(c - 'A' + 10)
	}
	return -1
}
