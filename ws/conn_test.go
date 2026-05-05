package ws

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wezzcoetzee/hyperliquid/transport"
)

type fakeWSConn struct {
	in      chan []byte
	closed  atomic.Bool
	written chan []byte
}

func newFakeWSConn() *fakeWSConn {
	return &fakeWSConn{in: make(chan []byte, 8), written: make(chan []byte, 8)}
}

func (f *fakeWSConn) WriteJSON(v any) error {
	if f.closed.Load() {
		return errors.New("closed")
	}
	select {
	case f.written <- []byte("write"):
	default:
	}
	return nil
}

func (f *fakeWSConn) ReadMessage(ctx context.Context) ([]byte, error) {
	select {
	case b, ok := <-f.in:
		if !ok {
			return nil, errors.New("eof")
		}
		return b, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (f *fakeWSConn) Close() error {
	if f.closed.CompareAndSwap(false, true) {
		close(f.in)
	}
	return nil
}

type fakeDialer struct {
	conns chan *fakeWSConn
}

func newFakeDialer() *fakeDialer {
	return &fakeDialer{conns: make(chan *fakeWSConn, 4)}
}

func (f *fakeDialer) Dial(ctx context.Context, url string) (transport.WSConn, error) {
	c := newFakeWSConn()
	f.conns <- c
	return c, nil
}

func TestConn_DispatchesByChannel(t *testing.T) {
	d := newFakeDialer()
	r := newRegistry()

	got := make(chan []byte, 1)
	r.add("allMids", func(b []byte) { got <- b })

	c := newConn(context.Background(), "ws://x", d, r)
	if err := c.ensure(context.Background()); err != nil {
		t.Fatal(err)
	}

	conn := <-d.conns
	conn.in <- []byte(`{"channel":"allMids","data":{"BTC":"30000"}}`)

	select {
	case b := <-got:
		if string(b) != `{"BTC":"30000"}` {
			t.Errorf("got %s", b)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no dispatch")
	}
	c.Close()
}

func TestConn_DispatchesByKey(t *testing.T) {
	d := newFakeDialer()
	r := newRegistry()

	got := make(chan []byte, 1)
	r.add("trades:BTC", func(b []byte) { got <- b })

	c := newConn(context.Background(), "ws://x", d, r)
	_ = c.ensure(context.Background())

	conn := <-d.conns
	conn.in <- []byte(`{"channel":"trades","data":[{"coin":"BTC"}]}`)

	select {
	case <-got:
	case <-time.After(2 * time.Second):
		t.Fatal("no dispatch to keyed subscriber")
	}
	c.Close()
}

func TestConn_Reconnects(t *testing.T) {
	d := newFakeDialer()
	r := newRegistry()
	r.add("trades:BTC", func([]byte) {})
	r.recordParams("trades:BTC", map[string]any{"type": "trades", "coin": "BTC"})

	c := newConn(context.Background(), "ws://x", d, r)
	c.reconnectMin = 10 * time.Millisecond
	c.reconnectMax = 10 * time.Millisecond
	if err := c.ensure(context.Background()); err != nil {
		t.Fatal(err)
	}

	first := <-d.conns
	first.Close()

	select {
	case <-d.conns:
	case <-time.After(2 * time.Second):
		t.Fatal("did not reconnect")
	}
	c.Close()
}
