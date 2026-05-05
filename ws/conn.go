package ws

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

// Conn is a single multiplexed WebSocket connection. It owns one read
// goroutine, one ping goroutine, and a write mutex (provided by WSConn).
// Subscriptions are tracked in registry; the dispatch loop fans incoming
// messages to subscribers based on the channel field.
type Conn struct {
	url      string
	dialer   transport.WSDialer
	registry *registry

	mu     sync.Mutex
	conn   transport.WSConn
	closed bool

	pingPeriod   time.Duration
	reconnectMin time.Duration
	reconnectMax time.Duration

	rootCtx context.Context
	cancel  context.CancelFunc
}

func newConn(parent context.Context, url string, d transport.WSDialer, r *registry) *Conn {
	ctx, cancel := context.WithCancel(parent)
	return &Conn{
		url:          url,
		dialer:       d,
		registry:     r,
		pingPeriod:   30 * time.Second,
		reconnectMin: 500 * time.Millisecond,
		reconnectMax: 30 * time.Second,
		rootCtx:      ctx,
		cancel:       cancel,
	}
}

// ensure dials and starts the read/ping goroutines if not already connected.
// Safe to call multiple times.
func (c *Conn) ensure(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return errors.New("ws: conn closed")
	}
	if c.conn != nil {
		return nil
	}
	conn, err := c.dialer.Dial(ctx, c.url)
	if err != nil {
		return err
	}
	c.conn = conn
	go c.readLoop(conn)
	go c.pingLoop()
	go c.replaySubs()
	return nil
}

func (c *Conn) readLoop(conn transport.WSConn) {
	for {
		msg, err := conn.ReadMessage(c.rootCtx)
		if err != nil {
			c.handleDisconnect(conn)
			return
		}
		c.dispatch(msg)
	}
}

func (c *Conn) dispatch(msg []byte) {
	var env struct {
		Channel string          `json:"channel"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(msg, &env); err != nil {
		return
	}
	for _, sub := range c.registry.subscribers(env.Channel) {
		sub.cb(env.Data)
	}
	prefix := env.Channel + ":"
	for k, subs := range c.registry.allByKey() {
		if strings.HasPrefix(k, prefix) {
			for _, s := range subs {
				s.cb(env.Data)
			}
		}
	}
}

func (c *Conn) pingLoop() {
	t := time.NewTicker(c.pingPeriod)
	defer t.Stop()
	for {
		select {
		case <-c.rootCtx.Done():
			return
		case <-t.C:
			c.mu.Lock()
			conn := c.conn
			c.mu.Unlock()
			if conn == nil {
				return
			}
			_ = conn.WriteJSON(map[string]string{"method": "ping"})
		}
	}
}

func (c *Conn) handleDisconnect(stale transport.WSConn) {
	c.mu.Lock()
	if c.conn == stale {
		_ = c.conn.Close()
		c.conn = nil
	}
	closed := c.closed
	c.mu.Unlock()
	if closed {
		return
	}
	go c.reconnect()
}

func (c *Conn) reconnect() {
	backoff := c.reconnectMin
	for {
		select {
		case <-c.rootCtx.Done():
			return
		case <-time.After(backoff):
		}
		if err := c.ensure(c.rootCtx); err == nil {
			return
		}
		backoff *= 2
		if backoff > c.reconnectMax {
			backoff = c.reconnectMax
		}
	}
}

func (c *Conn) replaySubs() {
	for _, params := range c.registry.snapshot() {
		_ = c.WriteJSON(map[string]any{"method": "subscribe", "subscription": params})
	}
}

// WriteJSON sends a JSON-encoded message. Returns an error if not connected.
func (c *Conn) WriteJSON(v any) error {
	c.mu.Lock()
	conn := c.conn
	c.mu.Unlock()
	if conn == nil {
		return errors.New("ws: not connected")
	}
	return conn.WriteJSON(v)
}

// Close shuts down the connection and stops all goroutines.
func (c *Conn) Close() {
	c.mu.Lock()
	c.closed = true
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()
	c.cancel()
}
