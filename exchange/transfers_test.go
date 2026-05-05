package exchange

import (
	"context"
	"testing"

	"github.com/wezzcoetzee/hyperliquid/signer/privkey"
)

func newPK(t *testing.T) *privkey.PrivKey {
	t.Helper()
	pk, err := privkey.New("0x0000000000000000000000000000000000000000000000000000000000000001")
	if err != nil {
		t.Fatal(err)
	}
	return pk
}

func okResponse() any {
	return map[string]any{"status": "ok", "response": map[string]any{"type": "default", "data": map[string]any{}}}
}

func TestClient_UsdSend(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "usdSend" || action["signatureChainId"] != "0xa4b1" || action["hyperliquidChain"] != "Mainnet" {
			t.Fatalf("action: %v", action)
		}
		if action["destination"] != "0x0000000000000000000000000000000000000002" || action["amount"] != "10" {
			t.Fatalf("body: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.UsdSend(context.Background(), "0x0000000000000000000000000000000000000002", "10"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SpotSend(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "spotSend" || action["token"] != "USDC:0x..." {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.SpotSend(context.Background(), "0xdest", "USDC:0x...", "5"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_Withdraw3(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "withdraw3" || action["amount"] != "100" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.Withdraw3(context.Background(), "0xdest", "100"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_UsdClassTransfer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "usdClassTransfer" || action["toPerp"] != true {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.UsdClassTransfer(context.Background(), "10", true); err != nil {
		t.Fatal(err)
	}
}

func TestClient_TokenDelegate(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "tokenDelegate" || action["isUndelegate"] != false {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.TokenDelegate(context.Background(), "0x0000000000000000000000000000000000000003", 1000, false); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SubAccountTransfer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "subAccountTransfer" || action["isDeposit"] != true {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.SubAccountTransfer(context.Background(), "0xsub", true, 100); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SubAccountSpotTransfer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "subAccountSpotTransfer" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.SubAccountSpotTransfer(context.Background(), "0xsub", false, "USDC", "1"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_VaultTransfer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "vaultTransfer" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.VaultTransfer(context.Background(), "0xv", true, 100); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ApproveAgent(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "approveAgent" || action["agentName"] != "trader1" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.ApproveAgent(context.Background(), "0x0000000000000000000000000000000000000004", "trader1"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ApproveBuilderFee(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "approveBuilderFee" || action["maxFeeRate"] != "0.001%" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.ApproveBuilderFee(context.Background(), "0x0000000000000000000000000000000000000005", "0.001%"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateSubAccount(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "createSubAccount" || action["name"] != "sub1" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.CreateSubAccount(context.Background(), "sub1"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SubAccountModify(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "subAccountModify" {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.SubAccountModify(context.Background(), "0xsub", "renamed"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_SetReferrer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		if req["action"].(map[string]any)["type"] != "setReferrer" {
			t.Fatalf("type: %v", req)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.SetReferrer(context.Background(), "ABC"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_RegisterReferrer(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		if req["action"].(map[string]any)["type"] != "registerReferrer" {
			t.Fatalf("type: %v", req)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.RegisterReferrer(context.Background(), "ABC"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreateVault(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		if req["action"].(map[string]any)["type"] != "createVault" {
			t.Fatalf("type: %v", req)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.CreateVault(context.Background(), CreateVaultRequest{Name: "V", Description: "", InitialUsd: 100}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_VaultModify(t *testing.T) {
	allow := true
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		action := req["action"].(map[string]any)
		if action["type"] != "vaultModify" || action["allowDeposits"] != true {
			t.Fatalf("action: %v", action)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.VaultModify(context.Background(), VaultModifyRequest{VaultAddress: "0xv", AllowDeposits: &allow}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_VaultDistribute(t *testing.T) {
	c, srv := newExchangeTestClient(t, newPK(t), SourceMainnet, 42161, func(req map[string]any) any {
		if req["action"].(map[string]any)["type"] != "vaultDistribute" {
			t.Fatalf("type: %v", req)
		}
		return okResponse()
	})
	defer srv.Close()

	if err := c.VaultDistribute(context.Background(), "0xv", 50); err != nil {
		t.Fatal(err)
	}
}
