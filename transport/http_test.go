package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultHTTP_PostJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/info" {
			t.Fatalf("unexpected req: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("missing content-type")
		}
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req["type"] != "meta" {
			t.Fatalf("body: %v", req)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"universe":[]}`))
	}))
	defer srv.Close()

	tr := NewDefaultHTTP(srv.URL, nil)
	var out struct {
		Universe []any `json:"universe"`
	}
	if err := tr.PostJSON(context.Background(), "/info", map[string]any{"type": "meta"}, &out); err != nil {
		t.Fatalf("PostJSON: %v", err)
	}
}

func TestDefaultHTTP_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("boom"))
	}))
	defer srv.Close()

	tr := NewDefaultHTTP(srv.URL, nil)
	var out any
	err := tr.PostJSON(context.Background(), "/info", map[string]any{}, &out)
	if err == nil {
		t.Fatal("expected error")
	}
}
