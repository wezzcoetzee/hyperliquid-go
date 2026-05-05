// Example: stream BTC trades from the Hyperliquid mainnet websocket.
//
// Run:
//   go run ./examples/ws-subscriptions
//
// Press Ctrl-C to stop.
package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/wezzcoetzee/hyperliquid-go"
	"github.com/wezzcoetzee/hyperliquid-go/ws"
)

func main() {
	c, err := hyperliquid.New(hyperliquid.Config{Network: hyperliquid.Mainnet})
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	sub, err := c.Subscriptions.Trades(ctx, "BTC", func(trades []ws.Trade) {
		for _, t := range trades {
			fmt.Printf("%s %s @ %s\n", t.Side, t.Sz, t.Px)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Unsubscribe(context.Background())

	<-ctx.Done()
	fmt.Println("\ngoodbye.")
}
