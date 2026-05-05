package info

import "encoding/json"

// ClearinghouseState is the perp-account snapshot for one user.
type ClearinghouseState struct {
	AssetPositions             []AssetPosition `json:"assetPositions"`
	CrossMaintenanceMarginUsed string          `json:"crossMaintenanceMarginUsed"`
	CrossMarginSummary         MarginSummary   `json:"crossMarginSummary"`
	MarginSummary              MarginSummary   `json:"marginSummary"`
	Time                       int64           `json:"time"`
	Withdrawable               string          `json:"withdrawable"`
}

// AssetPosition is one open perp position.
type AssetPosition struct {
	Type     string   `json:"type"`
	Position Position `json:"position"`
}

// Position is the inner position payload.
type Position struct {
	Coin           string   `json:"coin"`
	EntryPx        string   `json:"entryPx"`
	Leverage       Leverage `json:"leverage"`
	LiquidationPx  string   `json:"liquidationPx"`
	MarginUsed     string   `json:"marginUsed"`
	MaxLeverage    int      `json:"maxLeverage"`
	PositionValue  string   `json:"positionValue"`
	ReturnOnEquity string   `json:"returnOnEquity"`
	Szi            string   `json:"szi"`
	UnrealizedPnl  string   `json:"unrealizedPnl"`
}

// Leverage describes the per-position leverage setting.
type Leverage struct {
	Type   string `json:"type"`
	Value  int    `json:"value"`
	RawUsd string `json:"rawUsd,omitempty"`
}

// MarginSummary aggregates account-level margin info.
type MarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalMarginUsed string `json:"totalMarginUsed"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
}

// SpotBalance is one spot-account token balance.
type SpotBalance struct {
	Coin     string `json:"coin"`
	Token    int    `json:"token"`
	Hold     string `json:"hold"`
	Total    string `json:"total"`
	EntryNtl string `json:"entryNtl"`
}

// SpotClearinghouseState is the spot-account snapshot for one user.
type SpotClearinghouseState struct {
	Balances []SpotBalance `json:"balances"`
}

// OpenOrder is one resting order.
type OpenOrder struct {
	Coin      string `json:"coin"`
	LimitPx   string `json:"limitPx"`
	// Oid is the numeric order id assigned by the exchange.
	Oid  uint64 `json:"oid"`
	Side string `json:"side"`
	// Sz is the remaining unfilled size as a decimal string.
	Sz        string `json:"sz"`
	Timestamp int64  `json:"timestamp"`
}

// FrontendOpenOrder is OpenOrder enriched with order-type metadata.
type FrontendOpenOrder struct {
	OpenOrder
	OrigSz           string  `json:"origSz"`
	Cloid            *string `json:"cloid,omitempty"`
	OrderType        string  `json:"orderType"`
	ReduceOnly       bool    `json:"reduceOnly"`
	TriggerCondition string  `json:"triggerCondition,omitempty"`
	TriggerPx        string  `json:"triggerPx,omitempty"`
	IsTrigger        bool    `json:"isTrigger,omitempty"`
	IsPositionTpsl   bool    `json:"isPositionTpsl,omitempty"`
}

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

// LedgerEntry is one row of the user funding / non-funding ledger.
type LedgerEntry struct {
	Time  int64           `json:"time"`
	Hash  string          `json:"hash"`
	Delta json.RawMessage `json:"delta"`
}

// UserRateLimit summarizes the caller's current rate-limit budget.
type UserRateLimit struct {
	// CumVlm is the cumulative notional volume used for rate-limit tier calculation.
	CumVlm string `json:"cumVlm"`
	// NRequestsUsed is the number of API requests consumed in the current window.
	NRequestsUsed int `json:"nRequestsUsed"`
	// NRequestsCap is the maximum number of requests allowed in the current window.
	NRequestsCap int `json:"nRequestsCap"`
}

// OrderStatusResponse mirrors the orderStatus endpoint return shape.
type OrderStatusResponse struct {
	Status string          `json:"status"`
	Order  json.RawMessage `json:"order,omitempty"`
}

// HistoricalOrder is one historical order row.
type HistoricalOrder struct {
	Order           FrontendOpenOrder `json:"order"`
	Status          string            `json:"status"`
	StatusTimestamp int64             `json:"statusTimestamp"`
}

// SubAccount describes one sub-account belonging to the master.
type SubAccount struct {
	Name               string                 `json:"name"`
	SubAccountUser     string                 `json:"subAccountUser"`
	Master             string                 `json:"master"`
	ClearinghouseState ClearinghouseState     `json:"clearinghouseState"`
	SpotState          SpotClearinghouseState `json:"spotState"`
}

// ReferralState reports referral details for a user.
type ReferralState struct {
	ReferredBy       *json.RawMessage `json:"referredBy"`
	CumVlm           string           `json:"cumVlm"`
	UnclaimedRewards string           `json:"unclaimedRewards"`
	ClaimedRewards   string           `json:"claimedRewards"`
	BuilderRewards   string           `json:"builderRewards"`
	ReferrerState    json.RawMessage  `json:"referrerState"`
	RewardHistory    json.RawMessage  `json:"rewardHistory"`
}

// PortfolioPeriod is one window-period of a user's portfolio.
type PortfolioPeriod struct {
	AccountValueHistory [][2]any `json:"accountValueHistory"`
	PnlHistory          [][2]any `json:"pnlHistory"`
	Vlm                 string   `json:"vlm"`
}

// TwapState describes one active or historical TWAP order.
type TwapState struct {
	State  json.RawMessage `json:"state"`
	Status struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	} `json:"status"`
	Time int64 `json:"time"`
}
