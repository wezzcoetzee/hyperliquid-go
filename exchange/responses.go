package exchange

import "encoding/json"

// StatusResponse is the envelope returned by /exchange.
type StatusResponse struct {
	Status   string          `json:"status"`
	Response json.RawMessage `json:"response"`
}

// OrderResponseInner is the shape inside StatusResponse.Response for the
// "order" action.
type OrderResponseInner struct {
	Type string `json:"type"`
	Data struct {
		Statuses []OrderStatus `json:"statuses"`
	} `json:"data"`
}

// OrderStatus is one element of the Order response. Exactly one of Resting
// or Filled is non-nil on success; Error is set on per-order failure.
type OrderStatus struct {
	Resting *RestingOrder `json:"resting,omitempty"`
	Filled  *FilledOrder  `json:"filled,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// RestingOrder is the resting-order success payload.
type RestingOrder struct {
	Oid   uint64 `json:"oid"`
	Cloid string `json:"cloid,omitempty"`
}

// FilledOrder is the filled-order success payload.
type FilledOrder struct {
	Oid     uint64 `json:"oid"`
	TotalSz string `json:"totalSz"`
	AvgPx   string `json:"avgPx"`
	Cloid   string `json:"cloid,omitempty"`
}

// ActionRejected indicates a successful HTTP response whose body had
// status:"err". The Hyperliquid API uses this for per-action validation
// failures (e.g., bad tick size, insufficient margin).
type ActionRejected struct {
	// Action is the action type string that was rejected (e.g. "order").
	Action string
	// Response is the raw rejection message returned by the API.
	Response string
}

// Error implements the error interface.
func (e *ActionRejected) Error() string {
	return "exchange: action " + e.Action + " rejected: " + e.Response
}
