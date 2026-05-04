package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestDefaultHTTP_4xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		_, _ = w.Write([]byte("nope"))
	}))
	defer srv.Close()

	tr := NewDefaultHTTP(srv.URL, nil)
	err := tr.PostJSON(context.Background(), "/info", map[string]any{}, nil)
	var apiErr *TransportAPIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected TransportAPIError, got %v", err)
	}
	if apiErr.Status != 422 {
		t.Errorf("Status = %d", apiErr.Status)
	}
	if apiErr.Body != "nope" {
		t.Errorf("Body = %q", apiErr.Body)
	}
}

func TestDefaultHTTP_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html>oops</html>"))
	}))
	defer srv.Close()

	tr := NewDefaultHTTP(srv.URL, nil)
	var out map[string]any
	err := tr.PostJSON(context.Background(), "/info", map[string]any{}, &out)
	if err == nil {
		t.Fatal("expected decode error")
	}
	if !strings.Contains(err.Error(), "decode:") {
		t.Errorf("expected 'decode:' prefix, got %v", err)
	}
}

func TestDefaultHTTP_ContextCanceled(t *testing.T) {
	block := make(chan struct{})
	defer close(block)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-block
	}))
	defer srv.Close()

	tr := NewDefaultHTTP(srv.URL, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := tr.PostJSON(ctx, "/info", map[string]any{}, nil)
	if err == nil {
		t.Fatal("expected error from canceled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
