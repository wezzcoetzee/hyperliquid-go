package exchange

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
	"github.com/wezzcoetzee/hyperliquid-go/signer"
	"github.com/wezzcoetzee/hyperliquid-go/signer/privkey"
)

func assertL1Fixture(t *testing.T, name string, action *msgpack.OrderedMap) {
	t.Helper()
	f := loadFixture(t, name)
	pk, err := privkey.New(f.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	nonce, _ := strconv.ParseUint(f.Nonce, 10, 64)

	got, err := ActionHash(action, nonce, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if hex.EncodeToString(got) != stripHexPrefix(f.ActionHash) {
		t.Fatalf("hash mismatch:\n got %x\nwant %s", got, f.ActionHash)
	}

	sig, err := BuildL1Signature(context.Background(), pk, action, nonce, nil, nil, SourceMainnet)
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

func assertUserFixture(t *testing.T, name string, fields []signer.Field, msg map[string]any) {
	t.Helper()
	f := loadFixture(t, name)
	pk, err := privkey.New(f.PrivateKey)
	if err != nil {
		t.Fatal(err)
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

func TestBuildL1Signature_CancelByCloidFixture(t *testing.T) {
	cancel := msgpack.NewOrderedMap()
	cancel.Set("asset", uint64(0))
	cancel.Set("cloid", "0xdeadbeef00000000000000000000000000000000000000000000000000000001")
	a := msgpack.NewOrderedMap()
	a.Set("type", "cancelByCloid")
	a.Set("cancels", []any{cancel})
	assertL1Fixture(t, "cancel_by_cloid_l1", a)
}

func TestBuildL1Signature_ModifyFixture(t *testing.T) {
	order := msgpack.NewOrderedMap()
	order.Set("a", uint64(0))
	order.Set("b", true)
	order.Set("p", "31000")
	order.Set("s", "0.2")
	order.Set("r", false)
	lim := msgpack.NewOrderedMap()
	lim.Set("tif", "Gtc")
	ot := msgpack.NewOrderedMap()
	ot.Set("limit", lim)
	order.Set("t", ot)
	a := msgpack.NewOrderedMap()
	a.Set("type", "modify")
	a.Set("oid", uint64(999))
	a.Set("order", order)
	assertL1Fixture(t, "modify_l1", a)
}

func TestBuildL1Signature_BatchModifyFixture(t *testing.T) {
	buildOrder := func(asset uint64, isBuy bool, px, sz string, reduceOnly bool, tif string) *msgpack.OrderedMap {
		o := msgpack.NewOrderedMap()
		o.Set("a", asset)
		o.Set("b", isBuy)
		o.Set("p", px)
		o.Set("s", sz)
		o.Set("r", reduceOnly)
		lim := msgpack.NewOrderedMap()
		lim.Set("tif", tif)
		ot := msgpack.NewOrderedMap()
		ot.Set("limit", lim)
		o.Set("t", ot)
		return o
	}

	m1 := msgpack.NewOrderedMap()
	m1.Set("oid", uint64(111))
	m1.Set("order", buildOrder(0, true, "31000", "0.2", false, "Gtc"))

	m2 := msgpack.NewOrderedMap()
	m2.Set("oid", uint64(222))
	m2.Set("order", buildOrder(1, false, "2000", "1.5", true, "Ioc"))

	a := msgpack.NewOrderedMap()
	a.Set("type", "batchModify")
	a.Set("modifies", []any{m1, m2})
	assertL1Fixture(t, "batch_modify_l1", a)
}

func TestBuildL1Signature_ScheduleCancelFixture(t *testing.T) {
	a := msgpack.NewOrderedMap()
	a.Set("type", "scheduleCancel")
	a.Set("time", uint64(1700000060000))
	assertL1Fixture(t, "schedule_cancel_l1", a)
}

func TestBuildL1Signature_UpdateLeverageFixture(t *testing.T) {
	a := msgpack.NewOrderedMap()
	a.Set("type", "updateLeverage")
	a.Set("asset", uint64(0))
	a.Set("isCross", true)
	a.Set("leverage", uint64(10))
	assertL1Fixture(t, "update_leverage_l1", a)
}

func TestBuildL1Signature_UpdateIsolatedMarginFixture(t *testing.T) {
	a := msgpack.NewOrderedMap()
	a.Set("type", "updateIsolatedMargin")
	a.Set("asset", uint64(0))
	a.Set("isBuy", true)
	a.Set("ntli", int64(1000000))
	assertL1Fixture(t, "update_isolated_margin_l1", a)
}

func TestBuildL1Signature_TwapOrderFixture(t *testing.T) {
	twap := msgpack.NewOrderedMap()
	twap.Set("a", uint64(0))
	twap.Set("b", true)
	twap.Set("s", "1.0")
	twap.Set("r", false)
	twap.Set("m", uint64(5))
	twap.Set("t", false)
	a := msgpack.NewOrderedMap()
	a.Set("type", "twapOrder")
	a.Set("twap", twap)
	assertL1Fixture(t, "twap_order_l1", a)
}

func TestBuildL1Signature_TwapCancelFixture(t *testing.T) {
	a := msgpack.NewOrderedMap()
	a.Set("type", "twapCancel")
	a.Set("a", uint64(0))
	a.Set("t", uint64(42))
	assertL1Fixture(t, "twap_cancel_l1", a)
}

func TestBuildUserSignature_Withdraw3Fixture(t *testing.T) {
	assertUserFixture(t, "withdraw3", userSignActionFields.Withdraw3, map[string]any{
		"hyperliquidChain": "Mainnet",
		"destination":      "0x0000000000000000000000000000000000000002",
		"amount":           "5",
		"time":             uint64(1700000000000),
	})
}

func TestBuildUserSignature_SpotSendFixture(t *testing.T) {
	assertUserFixture(t, "spot_send", userSignActionFields.SpotSend, map[string]any{
		"hyperliquidChain": "Mainnet",
		"destination":      "0x0000000000000000000000000000000000000002",
		"token":            "USDC:0xeb62eee3685fc4c43992febcd9e75443",
		"amount":           "1",
		"time":             uint64(1700000000000),
	})
}

func TestBuildUserSignature_UsdClassTransferFixture(t *testing.T) {
	assertUserFixture(t, "usd_class_transfer", userSignActionFields.UsdClassTransfer, map[string]any{
		"hyperliquidChain": "Mainnet",
		"amount":           "100",
		"toPerp":           true,
		"nonce":            uint64(1700000000000),
	})
}

func TestBuildUserSignature_ApproveAgentFixture(t *testing.T) {
	assertUserFixture(t, "approve_agent", userSignActionFields.ApproveAgent, map[string]any{
		"hyperliquidChain": "Mainnet",
		"agentAddress":     "0x0000000000000000000000000000000000000003",
		"agentName":        "TestAgent",
		"nonce":            uint64(1700000000000),
	})
}

func TestBuildUserSignature_ApproveBuilderFeeFixture(t *testing.T) {
	assertUserFixture(t, "approve_builder_fee", userSignActionFields.ApproveBuilderFee, map[string]any{
		"hyperliquidChain": "Mainnet",
		"maxFeeRate":       "0.1%",
		"builder":          "0x0000000000000000000000000000000000000004",
		"nonce":            uint64(1700000000000),
	})
}

func TestBuildUserSignature_TokenDelegateFixture(t *testing.T) {
	assertUserFixture(t, "token_delegate", userSignActionFields.TokenDelegate, map[string]any{
		"hyperliquidChain": "Mainnet",
		"validator":        "0x0000000000000000000000000000000000000005",
		"wei":              uint64(1000000000000000000),
		"isUndelegate":     false,
		"nonce":            uint64(1700000000000),
	})
}
