package transport

import (
	"context"
	"sync"

	"github.com/gorilla/websocket"
)

// WSConn is the minimum surface a Hyperliquid WebSocket connection needs.
// The default implementation is goroutine-safe for writes via WriteJSON;
// reads are single-threaded by convention (only the connection's read goroutine
// should call ReadMessage).
type WSConn interface {
	WriteJSON(v any) error
	ReadMessage(ctx context.Context) ([]byte, error)
	Close() error
}

// WSDialer dials a WebSocket. The default implementation uses gorilla/websocket.
type WSDialer interface {
	Dial(ctx context.Context, url string) (WSConn, error)
}

// DefaultWSDialer is the gorilla-backed WSDialer.
type DefaultWSDialer struct{}

// Dial opens a WebSocket connection to url, honoring ctx for the handshake.
func (DefaultWSDialer) Dial(ctx context.Context, url string) (WSConn, error) {
	d := websocket.DefaultDialer
	c, _, err := d.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return &defaultWSConn{c: c}, nil
}

type defaultWSConn struct {
	c   *websocket.Conn
	wmu sync.Mutex
}

func (w *defaultWSConn) WriteJSON(v any) error {
	w.wmu.Lock()
	defer w.wmu.Unlock()
	return w.c.WriteJSON(v)
}

// ReadMessage blocks until a text/binary message arrives or ctx is cancelled.
// On cancellation the underlying connection is closed (gorilla doesn't expose
// a per-read deadline directly tied to context, so we close to unblock the
// read goroutine).
func (w *defaultWSConn) ReadMessage(ctx context.Context) ([]byte, error) {
	type res struct {
		b   []byte
		err error
	}
	ch := make(chan res, 1)
	go func() {
		_, b, err := w.c.ReadMessage()
		ch <- res{b, err}
	}()
	select {
	case r := <-ch:
		return r.b, r.err
	case <-ctx.Done():
		_ = w.c.Close()
		return nil, ctx.Err()
	}
}

func (w *defaultWSConn) Close() error { return w.c.Close() }
