package hyperliquid

import "fmt"

// APIError represents a non-2xx HTTP response from the Hyperliquid API.
type APIError struct {
	// Status is the HTTP response status code.
	Status int
	// Body is the raw response body text.
	Body string
	// Type is an optional machine-readable error category, if provided by the API.
	Type string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Type == "" {
		return fmt.Sprintf("hyperliquid: api error %d: %s", e.Status, e.Body)
	}
	return fmt.Sprintf("hyperliquid: api error %d (%s): %s", e.Status, e.Type, e.Body)
}

// ActionError represents an /exchange response with status:"err".
type ActionError struct {
	// Action is the action type string that was rejected (e.g. "order").
	Action string
	// Response is the raw rejection message returned by the API.
	Response string
}

// Error implements the error interface.
func (e *ActionError) Error() string {
	return fmt.Sprintf("hyperliquid: action %q rejected: %s", e.Action, e.Response)
}

// SignError wraps an error from the signing layer.
type SignError struct {
	// Err is the underlying signing error, accessible via errors.Unwrap.
	Err error
}

// Error implements the error interface.
func (e *SignError) Error() string {
	if e.Err == nil {
		return "hyperliquid: sign error"
	}
	return "hyperliquid: sign error: " + e.Err.Error()
}
// Unwrap returns the underlying signing error for errors.As/Is unwrapping.
func (e *SignError) Unwrap() error { return e.Err }

// WSError represents a WebSocket protocol error.
type WSError struct {
	// Code is the WebSocket close code.
	Code int
	// Reason is the human-readable close reason from the peer.
	Reason string
}

// Error implements the error interface.
func (e *WSError) Error() string {
	return fmt.Sprintf("hyperliquid: websocket error %d: %s", e.Code, e.Reason)
}
