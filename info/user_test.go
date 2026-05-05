package info

import (
	"context"
	"net/http"
	"testing"
)

const testUser = "0x0000000000000000000000000000000000000001"

func userJSONHandler(t *testing.T, expectedType, fixture string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != expectedType {
			t.Fatalf("type: %q want %q", got, expectedType)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, fixture))
	}
}

func TestClient_ClearinghouseState(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "clearinghouseState", "clearinghouse_state.json"))
	out, err := c.ClearinghouseState(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if out.Withdrawable != "1000" {
		t.Errorf("Withdrawable = %q", out.Withdrawable)
	}
}

func TestClient_SpotClearinghouseState(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "spotClearinghouseState", "spot_clearinghouse_state.json"))
	out, err := c.SpotClearinghouseState(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Balances) != 1 || out.Balances[0].Coin != "USDC" {
		t.Errorf("balances: %+v", out.Balances)
	}
}

func TestClient_OpenOrders(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "openOrders", "open_orders.json"))
	out, err := c.OpenOrders(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Oid != 12345 {
		t.Errorf("orders: %+v", out)
	}
}

func TestClient_FrontendOpenOrders(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "frontendOpenOrders", "frontend_open_orders.json"))
	out, err := c.FrontendOpenOrders(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].OrderType != "Limit" {
		t.Errorf("frontend: %+v", out)
	}
}

func TestClient_UserFills(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userFills", "user_fills.json"))
	out, err := c.UserFills(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Px != "30000" {
		t.Errorf("fills: %+v", out)
	}
}

func TestClient_UserFillsByTime(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userFillsByTime", "user_fills.json"))
	out, err := c.UserFillsByTime(context.Background(), testUser, 1700000000000, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Errorf("len = %d", len(out))
	}
}

func TestClient_UserFunding(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userFunding", "user_funding.json"))
	out, err := c.UserFunding(context.Background(), testUser, 1700000000000, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Errorf("len = %d", len(out))
	}
}

func TestClient_UserNonFundingLedgerUpdates(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userNonFundingLedgerUpdates", "user_non_funding_ledger.json"))
	out, err := c.UserNonFundingLedgerUpdates(context.Background(), testUser, 1700000000000, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Errorf("len = %d", len(out))
	}
}

func TestClient_UserRateLimit(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userRateLimit", "user_rate_limit.json"))
	out, err := c.UserRateLimit(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if out.NRequestsCap != 1200 {
		t.Errorf("cap = %d", out.NRequestsCap)
	}
}

func TestClient_OrderStatusByOID(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "orderStatus", "order_status.json"))
	out, err := c.OrderStatusByOID(context.Background(), testUser, 12345)
	if err != nil {
		t.Fatal(err)
	}
	if out.Status != "order" {
		t.Errorf("status: %q", out.Status)
	}
}

func TestClient_HistoricalOrders(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "historicalOrders", "historical_orders.json"))
	out, err := c.HistoricalOrders(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Status != "filled" {
		t.Errorf("orders: %+v", out)
	}
}

func TestClient_TwapHistory(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "twapHistory", "twap_history.json"))
	out, err := c.TwapHistory(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Errorf("len = %d", len(out))
	}
}

func TestClient_SubAccounts(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "subAccounts", "sub_accounts.json"))
	out, err := c.SubAccounts(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Name != "sub1" {
		t.Errorf("subs: %+v", out)
	}
}

func TestClient_Referral(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "referral", "referral.json"))
	out, err := c.Referral(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if out.CumVlm != "0" {
		t.Errorf("CumVlm: %q", out.CumVlm)
	}
}

func TestClient_PortfolioPeriods(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "portfolio", "portfolio.json"))
	out, err := c.PortfolioPeriods(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := out["day"]; !ok {
		t.Errorf("missing day: %+v", out)
	}
	if _, ok := out["week"]; !ok {
		t.Errorf("missing week: %+v", out)
	}
}
