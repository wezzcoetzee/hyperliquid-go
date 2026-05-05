package ws

// Trade is one trade event.
type Trade struct {
	Coin  string    `json:"coin"`
	Side  string    `json:"side"`
	Px    string    `json:"px"`
	Sz    string    `json:"sz"`
	Time  int64     `json:"time"`
	Hash  string    `json:"hash"`
	Tid   uint64    `json:"tid"`
	Users [2]string `json:"users,omitempty"`
}

// AllMidsEvent is the payload for the allMids channel.
type AllMidsEvent struct {
	Mids map[string]string `json:"mids"`
}

// L2Level is one price level in an L2BookEvent.
type L2Level struct {
	// Px is the price as a decimal string.
	Px string `json:"px"`
	// Sz is the aggregate size at this level as a decimal string.
	Sz string `json:"sz"`
	// N is the number of resting orders at this level.
	N int `json:"n"`
}

// L2BookEvent is the payload for the l2Book channel.
type L2BookEvent struct {
	Coin   string       `json:"coin"`
	Time   int64        `json:"time"`
	Levels [2][]L2Level `json:"levels"`
}

// Notification is the payload for the notification channel.
type Notification struct {
	Notification string `json:"notification"`
}

// CandleEvent is the payload for the candle channel.
type CandleEvent struct {
	// T is the candle open timestamp in milliseconds.
	T int64 `json:"t"`
	// C is the close price as a decimal string.
	C string `json:"c"`
	// H is the high price as a decimal string.
	H string `json:"h"`
	// L is the low price as a decimal string.
	L string `json:"l"`
	// O is the open price as a decimal string.
	O string `json:"o"`
	// V is the volume as a decimal string.
	V string `json:"v"`
	// N is the number of trades in the candle.
	N int `json:"n"`
	// I is the candle interval string (e.g. "1m", "1h").
	I string `json:"i"`
	// S is the coin/symbol for the candle.
	S string `json:"s"`
}

// BboEvent is the best-bid/best-offer event.
type BboEvent struct {
	Coin string      `json:"coin"`
	Time int64       `json:"time"`
	Bbo  [2]*L2Level `json:"bbo"`
}

// OrderUpdate is one order-status change.
type OrderUpdate struct {
	Order           map[string]any `json:"order"`
	Status          string         `json:"status"`
	StatusTimestamp int64          `json:"statusTimestamp"`
}

// UserEvent is one event in the user-events stream.
type UserEvent map[string]any

// Fill is one execution.
type Fill struct {
	Coin          string `json:"coin"`
	Px            string `json:"px"`
	Sz            string `json:"sz"`
	Side          string `json:"side"`
	Time          int64  `json:"time"`
	StartPosition string `json:"startPosition"`
	Dir           string `json:"dir"`
	ClosedPnl     string `json:"closedPnl"`
	Hash          string `json:"hash"`
	Oid           uint64 `json:"oid"`
	Cloid         string `json:"cloid,omitempty"`
	Crossed       bool   `json:"crossed"`
	Fee           string `json:"fee"`
	FeeToken      string `json:"feeToken"`
	Tid           uint64 `json:"tid"`
}

// UserFillsEvent is the payload for userFills.
type UserFillsEvent struct {
	User       string `json:"user"`
	Fills      []Fill `json:"fills"`
	IsSnapshot bool   `json:"isSnapshot,omitempty"`
}

// LedgerUpdate is one entry in the user funding / non-funding ledger.
type LedgerUpdate struct {
	Time  int64          `json:"time"`
	Hash  string         `json:"hash"`
	Delta map[string]any `json:"delta"`
}

// UserFundingsEvent is the payload for the userFundings channel.
type UserFundingsEvent struct {
	User       string         `json:"user"`
	Fundings   []LedgerUpdate `json:"fundings"`
	IsSnapshot bool           `json:"isSnapshot,omitempty"`
}

// UserLedgerEvent is the payload for the userNonFundingLedgerUpdates channel.
type UserLedgerEvent struct {
	User       string         `json:"user"`
	Updates    []LedgerUpdate `json:"nonFundingLedgerUpdates"`
	IsSnapshot bool           `json:"isSnapshot,omitempty"`
}

// ActiveAssetCtxEvent is the active-asset perp context.
type ActiveAssetCtxEvent struct {
	Coin string         `json:"coin"`
	Ctx  map[string]any `json:"ctx"`
}

// ActiveAssetDataEvent is per-user per-asset state.
type ActiveAssetDataEvent struct {
	User             string         `json:"user"`
	Coin             string         `json:"coin"`
	Leverage         map[string]any `json:"leverage"`
	MaxTradeSzs      []string       `json:"maxTradeSzs"`
	AvailableToTrade []string       `json:"availableToTrade"`
}

// TwapSliceFill is a single TWAP slice fill.
type TwapSliceFill struct {
	Fill   Fill   `json:"fill"`
	TwapID uint64 `json:"twapId"`
}

// TwapSliceFillsEvent is the payload for the twapSliceFills channel.
type TwapSliceFillsEvent struct {
	User           string          `json:"user"`
	TwapSliceFills []TwapSliceFill `json:"twapSliceFills"`
	IsSnapshot     bool            `json:"isSnapshot,omitempty"`
}

// TwapHistoryEvent is the payload for the twapHistory channel.
type TwapHistoryEvent struct {
	User       string           `json:"user"`
	History    []map[string]any `json:"history"`
	IsSnapshot bool             `json:"isSnapshot,omitempty"`
}
