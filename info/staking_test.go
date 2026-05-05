package info

import (
	"context"
	"testing"
)

func TestClient_Delegations(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "delegations", "delegations.json"))
	out, err := c.Delegations(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Validator != "0xval" {
		t.Errorf("delegations: %+v", out)
	}
}

func TestClient_DelegatorSummary(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "delegatorSummary", "delegator_summary.json"))
	out, err := c.DelegatorSummary(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if out.Delegated != "100" {
		t.Errorf("delegated: %q", out.Delegated)
	}
}

func TestClient_DelegatorHistory(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "delegatorHistory", "delegator_history.json"))
	out, err := c.DelegatorHistory(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 {
		t.Errorf("len = %d", len(out))
	}
}

func TestClient_DelegatorRewards(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "delegatorRewards", "delegator_rewards.json"))
	out, err := c.DelegatorRewards(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].TotalAmount != "1.5" {
		t.Errorf("rewards: %+v", out)
	}
}

func TestClient_ValidatorSummaries(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "validatorSummaries", "validator_summaries.json"))
	out, err := c.ValidatorSummaries(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Validator != "0xval" {
		t.Errorf("validators: %+v", out)
	}
}
