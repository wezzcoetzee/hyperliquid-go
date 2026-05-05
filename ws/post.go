package ws

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

const postTimeout = 10 * time.Second

// PostResult is the decoded response payload from a WS post call.
type PostResult struct {
	Response json.RawMessage
}

type pendingPost struct {
	ch chan postResponse
}

type postResponse struct {
	data json.RawMessage
	err  error
}

// postHub correlates outgoing post requests with their responses by id.
type postHub struct {
	nextID  atomic.Int64
	mu      sync.Mutex
	pending map[int64]*pendingPost
}

func newPostHub() *postHub {
	return &postHub{pending: make(map[int64]*pendingPost)}
}

func (h *postHub) register(id int64) *pendingPost {
	p := &pendingPost{ch: make(chan postResponse, 1)}
	h.mu.Lock()
	h.pending[id] = p
	h.mu.Unlock()
	return p
}

func (h *postHub) deliver(id int64, data json.RawMessage, err error) {
	h.mu.Lock()
	p, ok := h.pending[id]
	if ok {
		delete(h.pending, id)
	}
	h.mu.Unlock()
	if ok {
		p.ch <- postResponse{data: data, err: err}
	}
}

// Post sends a request over the WebSocket connection using the Hyperliquid
// post protocol and returns the response payload. The method and body are
// wrapped in {"method":"post","id":N,"request":{"method":method,"payload":body}}.
//
// This is a lower-level primitive. The caller is responsible for constructing
// the correct request payload. Responses are correlated by the numeric id
// assigned per call.
func (c *Client) Post(ctx context.Context, method string, body any) (json.RawMessage, error) {
	if err := c.ensureStarted(ctx); err != nil {
		return nil, err
	}

	c.postOnce.Do(func() { c.posts = newPostHub() })

	id := c.posts.nextID.Add(1)
	pending := c.posts.register(id)

	msg := map[string]any{
		"method": "post",
		"id":     id,
		"request": map[string]any{
			"type":    method,
			"payload": body,
		},
	}
	if err := c.conn.WriteJSON(msg); err != nil {
		c.posts.deliver(id, nil, err)
		return nil, err
	}

	timeout := postTimeout
	if dl, ok := ctx.Deadline(); ok {
		if remaining := time.Until(dl); remaining < timeout {
			timeout = remaining
		}
	}

	select {
	case resp := <-pending.ch:
		return resp.data, resp.err
	case <-time.After(timeout):
		c.posts.deliver(id, nil, nil)
		return nil, errors.New("ws: post timeout")
	case <-ctx.Done():
		c.posts.deliver(id, nil, nil)
		return nil, ctx.Err()
	}
}

// deliverPostResponse is called by the read loop when a "post" channel message arrives.
func (c *Client) deliverPostResponse(id int64, data json.RawMessage, errMsg string) {
	if c.posts == nil {
		return
	}
	if errMsg != "" {
		c.posts.deliver(id, nil, errors.New(errMsg))
		return
	}
	c.posts.deliver(id, data, nil)
}
