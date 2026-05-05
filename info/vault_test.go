package info

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_VaultDetails(t *testing.T) {
	c := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := readReqType(t, r.Body); got != "vaultDetails" {
			t.Fatalf("type: %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(loadResponse(t, "vault_details.json"))
	}))
	out, err := c.VaultDetails(context.Background(), "0xabc", "")
	if err != nil {
		t.Fatal(err)
	}
	if out.Name != "VaultX" {
		t.Errorf("name: %q", out.Name)
	}
}

func TestClient_VaultSummaries(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "vaultSummaries", "vault_summaries.json"))
	out, err := c.VaultSummaries(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].VaultAddress != "0xabc" {
		t.Errorf("summaries: %+v", out)
	}
}

func TestClient_UserVaultEquities(t *testing.T) {
	c := newTestClient(t, userJSONHandler(t, "userVaultEquities", "user_vault_equities.json"))
	out, err := c.UserVaultEquities(context.Background(), testUser)
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].Equity != "500" {
		t.Errorf("equities: %+v", out)
	}
}
