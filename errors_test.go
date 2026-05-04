package hyperliquid

import (
	"errors"
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
	if got := e.Error(); got == "" {
		t.Error("Error() returned empty")
	}
}

func TestActionError_Wraps(t *testing.T) {
	var e error = &ActionError{Action: "order", Response: "bad tick size"}
	var aerr *ActionError
	if !errors.As(e, &aerr) {
		t.Fatal("expected ActionError")
	}
}
