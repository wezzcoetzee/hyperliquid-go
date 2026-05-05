package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClient_Post_RoundTrip(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, raw, err := c.ReadMessage()
		if err != nil {
			return
		}
		var req struct {
			Method  string `json:"method"`
			ID      int64  `json:"id"`
			Request struct {
				Type    string `json:"type"`
				Payload any    `json:"payload"`
			} `json:"request"`
		}
		if err := json.Unmarshal(raw, &req); err != nil {
			return
		}

		reply := map[string]any{
			"channel": "post",
			"id":      req.ID,
			"data": map[string]any{
				"type":     req.Request.Type,
				"response": map[string]any{"status": "ok"},
			},
		}
		_ = c.WriteJSON(reply)
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	c := &Client{URL: url}
	defer c.Close()

	result, err := c.Post(context.Background(), "info", map[string]any{"type": "meta"})
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(result, &got); err != nil {
		t.Fatal(err)
	}
	if got["status"] != "ok" {
		t.Fatalf("unexpected response: %v", got)
	}
}

func TestClient_Post_ErrorResponse(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		_, raw, err := c.ReadMessage()
		if err != nil {
			return
		}
		var req struct {
			ID int64 `json:"id"`
		}
		_ = json.Unmarshal(raw, &req)

		reply := map[string]any{
			"channel": "post",
			"id":      req.ID,
			"data": map[string]any{
				"type":  "error",
				"error": "bad request",
			},
		}
		_ = c.WriteJSON(reply)
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	c := &Client{URL: url}
	defer c.Close()

	_, err := c.Post(context.Background(), "info", nil)
	if err == nil || err.Error() != "bad request" {
		t.Fatalf("expected 'bad request' error, got %v", err)
	}
}

func TestClient_Post_Timeout(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		_, _, _ = c.ReadMessage()
		time.Sleep(5 * time.Second)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	c := &Client{URL: url}
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := c.Post(ctx, "info", nil)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
