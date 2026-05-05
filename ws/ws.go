// Package ws provides multiplexed WebSocket subscriptions for Hyperliquid.
//
// Construct via hyperliquid.New; the Client exposes one method per channel.
// Each subscribe call dials lazily on first use, registers a handler, and
// returns a Subscription handle whose Unsubscribe method tears down the
// callback (and the server-side subscription if no other subscribers remain).
//
// The connection auto-reconnects with exponential backoff and replays all
// active subscriptions on reconnect. A single read goroutine multiplexes
// messages to per-key handlers.
package ws

import (
	"context"
	"sync"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Client is the WebSocket subscription client. Construct via hyperliquid.New.
//
// Client is safe for concurrent use. The connection is established lazily on
// the first subscribe call.
type Client struct {
	URL    string
	Dialer transport.WSDialer

	once     sync.Once
	conn     *Conn
	registry *registry
}

func (c *Client) ensureStarted(ctx context.Context) error {
	c.once.Do(func() {
		if c.Dialer == nil {
			c.Dialer = transport.DefaultWSDialer{}
		}
		c.registry = newRegistry()
		c.conn = newConn(context.Background(), c.URL, c.Dialer, c.registry)
	})
	return c.conn.ensure(ctx)
}

// Close shuts down the underlying connection. After Close, subscribe calls
// return an error.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// Subscription is the handle returned by every subscribe method. Call
// Unsubscribe to remove the callback.
type Subscription struct {
	c   *Client
	key string
	id  int64
	sub map[string]any
}

// Unsubscribe removes the callback. If this was the last subscriber for the
// key, the server-side subscription is also cancelled.
func (s *Subscription) Unsubscribe(ctx context.Context) error {
	last := s.c.registry.remove(s.key, s.id)
	if last {
		return s.c.conn.WriteJSON(map[string]any{"method": "unsubscribe", "subscription": s.sub})
	}
	return nil
}

func (c *Client) subscribe(ctx context.Context, key string, params map[string]any, cb func([]byte)) (*Subscription, error) {
	if err := c.ensureStarted(ctx); err != nil {
		return nil, err
	}
	id := c.registry.add(key, cb)
	c.registry.recordParams(key, params)
	if err := c.conn.WriteJSON(map[string]any{"method": "subscribe", "subscription": params}); err != nil {
		c.registry.remove(key, id)
		return nil, err
	}
	return &Subscription{c: c, key: key, id: id, sub: params}, nil
}
