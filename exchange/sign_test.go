package exchange

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid/signer"
	"github.com/wezzcoetzee/hyperliquid/signer/privkey"
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

	got, err := ActionHash(action, nonce, nil, nil)
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

	got, err := ActionHash(action, nonce, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != stripHexPrefix(f.ActionHash) {
		t.Fatalf("hash mismatch:\n got %x\nwant %s", got, f.ActionHash)
	}
}

func TestActionHash_VaultChangesHash(t *testing.T) {
	m := msgpack.NewOrderedMap()
	m.Set("type", "noop")

	noVault, err := ActionHash(m, 1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var vault [20]byte
	for i := range vault {
		vault[i] = byte(i + 1)
	}
	withVault, err := ActionHash(m, 1, &vault, nil)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(noVault, withVault) {
		t.Fatal("vault must change the hash")
	}
}

func TestActionHash_ExpiresAfterChangesHash(t *testing.T) {
	m := msgpack.NewOrderedMap()
	m.Set("type", "noop")

	base, err := ActionHash(m, 1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	exp := uint64(1700000000000)
	withExp, err := ActionHash(m, 1, nil, &exp)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(base, withExp) {
		t.Fatal("expiresAfter must change the hash")
	}
}

func TestBuildL1Signature_OrderFixture(t *testing.T) {
	f := loadFixture(t, "order_l1")
	pk, err := privkey.New(f.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	v, err := orderedFromJSON(f.Action)
	if err != nil {
		t.Fatal(err)
	}
	action := v.(*msgpack.OrderedMap)
	nonce, _ := strconv.ParseUint(f.Nonce, 10, 64)

	sig, err := BuildL1Signature(context.Background(), pk, action, nonce, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(sig.R[:]) != stripHexPrefix(f.Signature.R) {
		t.Errorf("r mismatch: got %x want %s", sig.R, f.Signature.R)
	}
	if hex.EncodeToString(sig.S[:]) != stripHexPrefix(f.Signature.S) {
		t.Errorf("s mismatch: got %x want %s", sig.S, f.Signature.S)
	}
	if int(sig.V) != f.Signature.V {
		t.Errorf("v mismatch: got %d want %d", sig.V, f.Signature.V)
	}
}

func TestBuildL1Signature_CancelFixture(t *testing.T) {
	f := loadFixture(t, "cancel_l1")
	pk, _ := privkey.New(f.PrivateKey)
	v, _ := orderedFromJSON(f.Action)
	action := v.(*msgpack.OrderedMap)
	nonce, _ := strconv.ParseUint(f.Nonce, 10, 64)

	sig, err := BuildL1Signature(context.Background(), pk, action, nonce, nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(sig.R[:]) != stripHexPrefix(f.Signature.R) {
		t.Errorf("r mismatch")
	}
	if hex.EncodeToString(sig.S[:]) != stripHexPrefix(f.Signature.S) {
		t.Errorf("s mismatch")
	}
	if int(sig.V) != f.Signature.V {
		t.Errorf("v mismatch: got %d want %d", sig.V, f.Signature.V)
	}
}

func TestBuildUserSignature_UsdSendFixture(t *testing.T) {
	f := loadFixture(t, "usd_send")
	pk, _ := privkey.New(f.PrivateKey)

	var actionFields map[string]json.RawMessage
	if err := json.Unmarshal(f.Action, &actionFields); err != nil {
		t.Fatal(err)
	}

	var hlChain, dest, amount string
	if err := json.Unmarshal(actionFields["hyperliquidChain"], &hlChain); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(actionFields["destination"], &dest); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(actionFields["amount"], &amount); err != nil {
		t.Fatal(err)
	}
	var n json.Number
	if err := json.Unmarshal(actionFields["time"], &n); err != nil {
		t.Fatal(err)
	}
	ti, _ := n.Int64()
	time := uint64(ti)

	fields := []signer.Field{
		{Name: "hyperliquidChain", Type: "string"},
		{Name: "destination", Type: "string"},
		{Name: "amount", Type: "string"},
		{Name: "time", Type: "uint64"},
	}
	msg := map[string]any{
		"hyperliquidChain": hlChain,
		"destination":      dest,
		"amount":           amount,
		"time":             time,
	}

	var us struct {
		Type    string `json:"type"`
		ChainID uint64 `json:"chainId"`
	}
	if err := json.Unmarshal(f.UserSigned, &us); err != nil {
		t.Fatal(err)
	}

	sig, err := BuildUserSignature(context.Background(), pk, us.Type, fields, msg, us.ChainID)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(sig.R[:]) != stripHexPrefix(f.Signature.R) {
		t.Errorf("r mismatch: got %x want %s", sig.R, f.Signature.R)
	}
	if hex.EncodeToString(sig.S[:]) != stripHexPrefix(f.Signature.S) {
		t.Errorf("s mismatch: got %x want %s", sig.S, f.Signature.S)
	}
	if int(sig.V) != f.Signature.V {
		t.Errorf("v mismatch: got %d want %d", sig.V, f.Signature.V)
	}
}
