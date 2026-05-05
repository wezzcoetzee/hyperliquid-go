package info

import (
	"context"
	"testing"
)

func TestClient_MaxBuilderFee(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "maxBuilderFee", "max_builder_fee.json"))
	got, err := c.MaxBuilderFee(context.Background(), testUser, "0xb")
	if err != nil {
		t.Fatal(err)
	}
	if got != 1 {
		t.Errorf("fee = %d", got)
	}
}

func TestClient_LegalCheck(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "legalCheck", "legal_check.json"))
	out, err := c.LegalCheck(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Accepted {
		t.Errorf("accepted: %v", out.Accepted)
	}
}

func TestClient_ExtraAgents(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "extraAgents", "extra_agents.json"))
	out, err := c.ExtraAgents(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Name != "agent1" {
		t.Errorf("agents: %+v", out)
	}
}
