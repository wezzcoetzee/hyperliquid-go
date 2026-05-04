package hyperliquid

import (
	"errors"
	"strings"
	"testing"
)

func TestAPIError_ErrorsAs(t *testing.T) {
	var e error = &APIError{Status: 500, Body: "boom", Type: "internal"}
	var apiErr *APIError
	if !errors.As(e, &apiErr) {
		t.Fatal("expected APIError via errors.As")
	}
	if apiErr.Status != 500 {
		t.Errorf("status = %d", apiErr.Status)
	}
	if got := e.Error(); !strings.Contains(got, "hyperliquid:") || !strings.Contains(got, "internal") || !strings.Contains(got, "boom") {
		t.Errorf("Error() = %q", got)
	}
}

func TestAPIError_EmptyType(t *testing.T) {
	e := &APIError{Status: 502, Body: "bad gateway"}
	got := e.Error()
	if strings.Contains(got, "()") {
		t.Errorf("empty parens should be omitted: %q", got)
	}
	if !strings.Contains(got, "502") || !strings.Contains(got, "bad gateway") {
		t.Errorf("Error() = %q", got)
	}
}

func TestActionError_ErrorsAs(t *testing.T) {
	var e error = &ActionError{Action: "order", Response: "bad tick size"}
	var aerr *ActionError
	if !errors.As(e, &aerr) {
		t.Fatal("expected ActionError")
	}
	if aerr.Action != "order" {
		t.Errorf("Action = %q", aerr.Action)
	}
}

var errSentinel = errors.New("sentinel")

func TestSignError_Unwrap(t *testing.T) {
	e := &SignError{Err: errSentinel}
	if !errors.Is(e, errSentinel) {
		t.Fatal("errors.Is should reach sentinel via Unwrap")
	}
	if !strings.Contains(e.Error(), "sentinel") {
		t.Errorf("Error() = %q", e.Error())
	}
}

func TestSignError_NilErr(t *testing.T) {
	e := &SignError{}
	got := e.Error() // must not panic
	if got == "" {
		t.Error("Error() returned empty")
	}
}

func TestWSError_ErrorsAs(t *testing.T) {
	var e error = &WSError{Code: 1011, Reason: "internal"}
	var wsErr *WSError
	if !errors.As(e, &wsErr) {
		t.Fatal("expected WSError")
	}
	if wsErr.Code != 1011 {
		t.Errorf("Code = %d", wsErr.Code)
	}
	if !strings.Contains(e.Error(), "1011") {
		t.Errorf("Error() = %q", e.Error())
	}
}
