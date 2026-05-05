package exchange

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
)

type fixture struct {
	PrivateKey string          `json:"privateKey"`
	Nonce      string          `json:"nonce"`
	Action     json.RawMessage `json:"action"`
	ActionHash string          `json:"actionHash"`
	Signature  struct {
		R string `json:"r"`
		S string `json:"s"`
		V int    `json:"v"`
	} `json:"signature"`
	UserSigned json.RawMessage `json:"userSigned"`
}

func loadFixture(t *testing.T, name string) fixture {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "fixtures", name+".json"))
	if err != nil {
		t.Fatal(err)
	}
	var f fixture
	if err := json.Unmarshal(b, &f); err != nil {
		t.Fatal(err)
	}
	return f
}

func orderedFromJSON(raw json.RawMessage) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	v, err := decodeValue(dec)
	if err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, errors.New("trailing json")
	}
	return v, nil
}

func decodeValue(dec *json.Decoder) (any, error) {
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	switch t := tok.(type) {
	case json.Delim:
		if t == '{' {
			m := msgpack.NewOrderedMap()
			for dec.More() {
				keyTok, err := dec.Token()
				if err != nil {
					return nil, err
				}
				key, ok := keyTok.(string)
				if !ok {
					return nil, errors.New("expected string key")
				}
				val, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				m.Set(key, val)
			}
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			return m, nil
		}
		if t == '[' {
			var arr []any
			for dec.More() {
				v, err := decodeValue(dec)
				if err != nil {
					return nil, err
				}
				arr = append(arr, v)
			}
			if _, err := dec.Token(); err != nil {
				return nil, err
			}
			if arr == nil {
				return []any{}, nil
			}
			return arr, nil
		}
		return nil, errors.New("unexpected delim")
	case string:
		return t, nil
	case bool:
		return t, nil
	case nil:
		return nil, nil
	case json.Number:
		i, err := t.Int64()
		if err != nil {
			return nil, err
		}
		return i, nil
	}
	return nil, errors.New("unsupported json token")
}

func stripHexPrefix(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func TestActionHash_OrderL1(t *testing.T) {
	f := loadFixture(t, "order_l1")
	v, err := orderedFromJSON(f.Action)
	if err != nil {
		t.Fatal(err)
	}
	action, ok := v.(*msgpack.OrderedMap)
	if !ok {
		t.Fatalf("action is %T, want *msgpack.OrderedMap", v)
	}

	nonce, err := strconv.ParseUint(f.Nonce, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ActionHash(action, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != stripHexPrefix(f.ActionHash) {
		t.Fatalf("hash mismatch:\n got %x\nwant %s", got, f.ActionHash)
	}
}

func TestActionHash_CancelL1(t *testing.T) {
	f := loadFixture(t, "cancel_l1")
	v, err := orderedFromJSON(f.Action)
	if err != nil {
		t.Fatal(err)
	}
	action := v.(*msgpack.OrderedMap)

	nonce, _ := strconv.ParseUint(f.Nonce, 10, 64)

	got, err := ActionHash(action, nonce, nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != stripHexPrefix(f.ActionHash) {
		t.Fatalf("hash mismatch:\n got %x\nwant %s", got, f.ActionHash)
	}
}

func TestActionHash_RejectsBadVault(t *testing.T) {
	m := msgpack.NewOrderedMap()
	m.Set("type", "noop")
	if _, err := ActionHash(m, 1, []byte{0x01, 0x02}); err == nil {
		t.Fatal("expected error for short vault")
	}
}
