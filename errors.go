package hyperliquid

import "fmt"

// APIError represents a non-2xx HTTP response from the Hyperliquid API.
type APIError struct {
	Status int
	Body   string
	Type   string
}

func (e *APIError) Error() string {
	if e.Type == "" {
		return fmt.Sprintf("hyperliquid: api error %d: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("hyperliquid: api error %d (%s): %s", e.Status, e.Type, e.Body)
}

// ActionError represents an /exchange response with status:"err".
type ActionError struct {
	Action   string
	Response string
}

func (e *ActionError) Error() string {
	return fmt.Sprintf("hyperliquid: action %q rejected: %s", e.Action, e.Response)
}

// SignError wraps an error from the signing layer.
type SignError struct {
	Err error
}

func (e *SignError) Error() string {
	if e.Err == nil {
		return "hyperliquid: sign error"
	}
	return "hyperliquid: sign error: " + e.Err.Error()
}
func (e *SignError) Unwrap() error { return e.Err }

// WSError represents a WebSocket protocol error.
type WSError struct {
	Code   int
	Reason string
}

func (e *WSError) Error() string {
	return fmt.Sprintf("hyperliquid: websocket error %d: %s", e.Code, e.Reason)
}
