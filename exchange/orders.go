package exchange

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid-go/internal/msgpack"
)

// Tif is a time-in-force qualifier for limit orders.
type Tif string

const (
	// TifGtc is Good-Till-Cancelled: the order rests until filled or explicitly cancelled.
	TifGtc Tif = "Gtc"
	// TifIoc is Immediate-Or-Cancel: any unfilled portion is cancelled immediately.
	TifIoc Tif = "Ioc"
	// TifAlo is Add-Liquidity-Only (post-only): the order is cancelled if it would cross.
	TifAlo Tif = "Alo"
)

// LimitOrder is the limit-order order-type leaf.
type LimitOrder struct {
	// Tif is the time-in-force qualifier (Gtc, Ioc, or Alo).
	Tif Tif
}

// TriggerOrder is the trigger-order order-type leaf (TP/SL).
type TriggerOrder struct {
	// IsMarket indicates the trigger fills as a market order (true) or limit order (false).
	IsMarket bool
	// TriggerPx is the trigger price as a decimal string.
	TriggerPx string
	// Tpsl is either "tp" (take-profit) or "sl" (stop-loss).
	Tpsl string
}

// OrderType picks one of Limit or Trigger. Exactly one must be non-nil.
type OrderType struct {
	Limit   *LimitOrder
	Trigger *TriggerOrder
}

// OrderParams describes a single order in an Order request.
type OrderParams struct {
	// Asset is the numeric asset index as returned by info.Meta.
	Asset uint32
	// IsBuy is true for a buy (long) order, false for a sell (short).
	IsBuy bool
	// LimitPx is the limit price as a decimal string.
	LimitPx string
	// Sz is the order size as a decimal string.
	Sz string
	// ReduceOnly constrains the order to only reduce an existing position.
	ReduceOnly bool
	// OrderType selects limit or trigger order semantics; exactly one must be set.
	OrderType OrderType
	// Cloid is the optional 0x-prefixed 32-char hex client order id.
	Cloid string
}

// Builder describes an optional builder-fee recipient for the action.
type Builder struct {
	// Address is the builder's 0x-prefixed Ethereum address.
	Address string
	// Fee is the builder fee in tenths of a basis point (e.g. 10 == 1 bps).
	Fee int
}

// OrderRequest groups one or more orders with a grouping mode and optional builder.
type OrderRequest struct {
	Orders   []OrderParams
	Grouping string
	Builder  *Builder
}

// OrderResponse is the result of a successful Order call.
type OrderResponse struct {
	Statuses []OrderStatus
}

func buildOrderAction(req OrderRequest) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "order")

	orders := make([]any, len(req.Orders))
	for i, o := range req.Orders {
		om := msgpack.NewOrderedMap()
		om.Set("a", uint64(o.Asset))
		om.Set("b", o.IsBuy)
		om.Set("p", o.LimitPx)
		om.Set("s", o.Sz)
		om.Set("r", o.ReduceOnly)

		t := msgpack.NewOrderedMap()
		switch {
		case o.OrderType.Limit != nil:
			lim := msgpack.NewOrderedMap()
			lim.Set("tif", string(o.OrderType.Limit.Tif))
			t.Set("limit", lim)
		case o.OrderType.Trigger != nil:
			trg := msgpack.NewOrderedMap()
			trg.Set("isMarket", o.OrderType.Trigger.IsMarket)
			trg.Set("triggerPx", o.OrderType.Trigger.TriggerPx)
			trg.Set("tpsl", o.OrderType.Trigger.Tpsl)
			t.Set("trigger", trg)
		}
		om.Set("t", t)
		if o.Cloid != "" {
			om.Set("c", o.Cloid)
		}
		orders[i] = om
	}
	a.Set("orders", orders)
	a.Set("grouping", req.Grouping)
	if req.Builder != nil {
		b := msgpack.NewOrderedMap()
		b.Set("b", req.Builder.Address)
		b.Set("f", uint64(req.Builder.Fee))
		a.Set("builder", b)
	}
	return a
}

// Order submits one or more orders.
func (c *Client) Order(ctx context.Context, req OrderRequest) (*OrderResponse, error) {
	action := buildOrderAction(req)
	var inner OrderResponseInner
	if err := c.submitL1(ctx, action, &inner); err != nil {
		return nil, err
	}
	return &OrderResponse{Statuses: inner.Data.Statuses}, nil
}

// CancelParams identifies one perp order to cancel by oid.
type CancelParams struct {
	Asset uint32
	Oid   uint64
}

func buildCancelAction(cancels []CancelParams) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "cancel")
	arr := make([]any, len(cancels))
	for i, c := range cancels {
		om := msgpack.NewOrderedMap()
		om.Set("a", uint64(c.Asset))
		om.Set("o", c.Oid)
		arr[i] = om
	}
	a.Set("cancels", arr)
	return a
}

// CancelResponse is the inner response shape for cancel/cancelByCloid.
type CancelResponse struct {
	Statuses []string `json:"statuses"`
}

// Cancel cancels one or more orders by oid.
func (c *Client) Cancel(ctx context.Context, cancels []CancelParams) (*CancelResponse, error) {
	action := buildCancelAction(cancels)
	var inner struct {
		Type string         `json:"type"`
		Data CancelResponse `json:"data"`
	}
	if err := c.submitL1(ctx, action, &inner); err != nil {
		return nil, err
	}
	return &inner.Data, nil
}

// CancelByCloidParams identifies one order to cancel by client-supplied cloid.
type CancelByCloidParams struct {
	Asset uint32
	Cloid string
}

func buildCancelByCloidAction(cancels []CancelByCloidParams) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "cancelByCloid")
	arr := make([]any, len(cancels))
	for i, c := range cancels {
		om := msgpack.NewOrderedMap()
		om.Set("asset", uint64(c.Asset))
		om.Set("cloid", c.Cloid)
		arr[i] = om
	}
	a.Set("cancels", arr)
	return a
}

// CancelByCloid cancels one or more orders by client-supplied cloid.
func (c *Client) CancelByCloid(ctx context.Context, cancels []CancelByCloidParams) (*CancelResponse, error) {
	action := buildCancelByCloidAction(cancels)
	var inner struct {
		Type string         `json:"type"`
		Data CancelResponse `json:"data"`
	}
	if err := c.submitL1(ctx, action, &inner); err != nil {
		return nil, err
	}
	return &inner.Data, nil
}

// ModifyParams modifies one resting order to new price/size/order-type.
type ModifyParams struct {
	Oid      uint64
	NewOrder OrderParams
}

func buildOneOrder(o OrderParams) *msgpack.OrderedMap {
	om := msgpack.NewOrderedMap()
	om.Set("a", uint64(o.Asset))
	om.Set("b", o.IsBuy)
	om.Set("p", o.LimitPx)
	om.Set("s", o.Sz)
	om.Set("r", o.ReduceOnly)
	t := msgpack.NewOrderedMap()
	switch {
	case o.OrderType.Limit != nil:
		lim := msgpack.NewOrderedMap()
		lim.Set("tif", string(o.OrderType.Limit.Tif))
		t.Set("limit", lim)
	case o.OrderType.Trigger != nil:
		trg := msgpack.NewOrderedMap()
		trg.Set("isMarket", o.OrderType.Trigger.IsMarket)
		trg.Set("triggerPx", o.OrderType.Trigger.TriggerPx)
		trg.Set("tpsl", o.OrderType.Trigger.Tpsl)
		t.Set("trigger", trg)
	}
	om.Set("t", t)
	if o.Cloid != "" {
		om.Set("c", o.Cloid)
	}
	return om
}

func buildModifyAction(p ModifyParams) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "modify")
	a.Set("oid", p.Oid)
	a.Set("order", buildOneOrder(p.NewOrder))
	return a
}

// Modify modifies one resting order.
func (c *Client) Modify(ctx context.Context, p ModifyParams) error {
	return c.submitL1(ctx, buildModifyAction(p), nil)
}

func buildBatchModifyAction(modifies []ModifyParams) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "batchModify")
	arr := make([]any, len(modifies))
	for i, m := range modifies {
		om := msgpack.NewOrderedMap()
		om.Set("oid", m.Oid)
		om.Set("order", buildOneOrder(m.NewOrder))
		arr[i] = om
	}
	a.Set("modifies", arr)
	return a
}

// BatchModify modifies several resting orders atomically.
func (c *Client) BatchModify(ctx context.Context, modifies []ModifyParams) error {
	return c.submitL1(ctx, buildBatchModifyAction(modifies), nil)
}

// ScheduleCancel arms a dead-man's-switch. Pass deadlineMs == 0 to disarm.
func (c *Client) ScheduleCancel(ctx context.Context, deadlineMs int64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "scheduleCancel")
	if deadlineMs > 0 {
		a.Set("time", uint64(deadlineMs))
	}
	return c.submitL1(ctx, a, nil)
}

// UpdateLeverage changes the leverage setting on one perp asset.
func (c *Client) UpdateLeverage(ctx context.Context, asset uint32, isCross bool, leverage int) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "updateLeverage")
	a.Set("asset", uint64(asset))
	a.Set("isCross", isCross)
	a.Set("leverage", uint64(leverage))
	return c.submitL1(ctx, a, nil)
}

// UpdateIsolatedMargin adds (positive) or removes (negative) margin on an isolated perp position.
func (c *Client) UpdateIsolatedMargin(ctx context.Context, asset uint32, isBuy bool, ntli int64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "updateIsolatedMargin")
	a.Set("asset", uint64(asset))
	a.Set("isBuy", isBuy)
	a.Set("ntli", ntli)
	return c.submitL1(ctx, a, nil)
}

// TwapParams describes a TWAP order request.
type TwapParams struct {
	Asset      uint32
	IsBuy      bool
	Sz         string
	ReduceOnly bool
	Minutes    int
	Randomize  bool
}

func buildTwapOrderAction(p TwapParams) *msgpack.OrderedMap {
	a := msgpack.NewOrderedMap()
	a.Set("type", "twapOrder")
	t := msgpack.NewOrderedMap()
	t.Set("a", uint64(p.Asset))
	t.Set("b", p.IsBuy)
	t.Set("s", p.Sz)
	t.Set("r", p.ReduceOnly)
	t.Set("m", uint64(p.Minutes))
	t.Set("t", p.Randomize)
	a.Set("twap", t)
	return a
}

// TwapOrder submits a TWAP order.
func (c *Client) TwapOrder(ctx context.Context, p TwapParams) error {
	return c.submitL1(ctx, buildTwapOrderAction(p), nil)
}

// TwapCancel cancels a running TWAP by id.
func (c *Client) TwapCancel(ctx context.Context, asset uint32, twapID uint64) error {
	a := msgpack.NewOrderedMap()
	a.Set("type", "twapCancel")
	a.Set("a", uint64(asset))
	a.Set("t", twapID)
	return c.submitL1(ctx, a, nil)
}
