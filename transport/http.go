// Package transport provides pluggable HTTP and WebSocket interfaces used by
// the Hyperliquid client. Default implementations sit behind these interfaces
// so callers can swap in custom plumbing for testing or proxying.
package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTP is the minimum surface the SDK needs to make a JSON POST.
type HTTP interface {
	PostJSON(ctx context.Context, path string, body any, out any) error
}

// DefaultHTTP is a net/http-backed HTTP. Use NewDefaultHTTP to construct.
type DefaultHTTP struct {
	baseURL string
	client  *http.Client
}

// NewDefaultHTTP returns a DefaultHTTP. If client is nil, a *http.Client with a
// 30-second timeout is used.
func NewDefaultHTTP(baseURL string, client *http.Client) *DefaultHTTP {
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	return &DefaultHTTP{baseURL: baseURL, client: client}
}

// PostJSON marshals body to JSON, POSTs it to baseURL+path, and decodes the
// response into out (skipped if out is nil). Non-2xx responses return a
// *TransportAPIError so the root package can wrap them with richer context.
func (h *DefaultHTTP) PostJSON(ctx context.Context, path string, body any, out any) error {
	buf, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", h.baseURL+path, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 != 2 {
		return &TransportAPIError{Status: resp.StatusCode, Body: string(respBody)}
	}
	if out == nil {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode: %w (body=%s)", err, string(respBody))
	}
	return nil
}

// TransportAPIError is returned for non-2xx HTTP responses. The root package
// wraps this into hyperliquid.APIError to avoid an import cycle.
type TransportAPIError struct {
	Status int
	Body   string
}

func (e *TransportAPIError) Error() string {
	return fmt.Sprintf("http %d: %s", e.Status, e.Body)
}
