// Example: read-only market data and account state.
//
// Run:
//   go run ./examples/info-only
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/wezzcoetzee/hyperliquid-go"
)

func main() {
	c, err := hyperliquid.New(hyperliquid.Config{Network: hyperliquid.Mainnet})
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	mids, err := c.Info.AllMids(ctx)
	if err != nil {
		log.Fatal(err)
	}
	b, _ := json.MarshalIndent(mids, "", "  ")
	fmt.Println("--- AllMids ---")
	fmt.Println(string(b))

	meta, err := c.Info.Meta(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n--- Meta ---\n%d perp assets\n", len(meta.Universe))
	for i, a := range meta.Universe {
		if i >= 5 {
			fmt.Println("  ...")
			break
		}
		fmt.Printf("  %s (szDecimals=%d, maxLev=%d)\n", a.Name, a.SzDecimals, a.MaxLeverage)
	}
}
