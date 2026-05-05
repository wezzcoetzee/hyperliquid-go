package exchange

import (
	"context"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
)

// Tif is a time-in-force qualifier for limit orders.
type Tif string

const (
	TifGtc Tif = "Gtc"
	TifIoc Tif = "Ioc"
	TifAlo Tif = "Alo"
)

// LimitOrder is the limit-order order-type leaf.
type LimitOrder struct {
	Tif Tif
}

// TriggerOrder is the trigger-order order-type leaf (TP/SL).
type TriggerOrder struct {
	IsMarket  bool
	TriggerPx string
	Tpsl      string
}

// OrderType picks one of Limit or Trigger. Exactly one must be non-nil.
type OrderType struct {
	Limit   *LimitOrder
	Trigger *TriggerOrder
}

// OrderParams describes a single order in an Order request.
type OrderParams struct {
	Asset      uint32
	IsBuy      bool
	LimitPx    string
	Sz         string
	ReduceOnly bool
	OrderType  OrderType
	Cloid      string
}

// Builder describes an optional builder-fee recipient for the action.
type Builder struct {
	Address string
	Fee     int
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
