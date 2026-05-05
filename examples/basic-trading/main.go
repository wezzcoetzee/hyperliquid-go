// Example: place + cancel a way-out-of-market limit order on testnet.
//
// Run (with a funded testnet key):
//   HL_TESTNET_KEY=0x... go run ./examples/basic-trading
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wezzcoetzee/hyperliquid-go"
	"github.com/wezzcoetzee/hyperliquid-go/exchange"
	"github.com/wezzcoetzee/hyperliquid-go/signer/privkey"
)

func main() {
	pk := os.Getenv("HL_TESTNET_KEY")
	if pk == "" {
		log.Fatal("set HL_TESTNET_KEY (testnet only)")
	}
	signer, err := privkey.New(pk)
	if err != nil {
		log.Fatal(err)
	}
	c, err := hyperliquid.New(hyperliquid.Config{
		Network: hyperliquid.Testnet,
		Signer:  signer,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	resp, err := c.Exchange.Order(ctx, exchange.OrderRequest{
		Orders: []exchange.OrderParams{{
			Asset: 0, IsBuy: true, LimitPx: "1", Sz: "0.001",
			OrderType: exchange.OrderType{Limit: &exchange.LimitOrder{Tif: exchange.TifGtc}},
		}},
		Grouping: "na",
	})
	if err != nil {
		log.Fatal(err)
	}
	if resting := resp.Statuses[0].Resting; resting != nil {
		fmt.Println("placed oid:", resting.Oid)
		if _, err := c.Exchange.Cancel(ctx, []exchange.CancelParams{{Asset: 0, Oid: resting.Oid}}); err != nil {
			log.Fatalf("cancel: %v", err)
		}
		fmt.Println("cancelled.")
	} else {
		fmt.Printf("no resting order: %+v\n", resp.Statuses)
	}
}
