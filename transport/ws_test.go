package transport

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestDefaultWS_RoundTrip(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		_, msg, _ := c.ReadMessage()
		_ = c.WriteMessage(websocket.TextMessage, []byte("echo:"+string(msg)))
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := DefaultWSDialer{}.Dial(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err := conn.WriteJSON(map[string]string{"hi": "there"}); err != nil {
		t.Fatal(err)
	}
	b, err := conn.ReadMessage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(b), "echo:") {
		t.Fatalf("got %q", b)
	}
}

func TestDefaultWS_ReadCancelledByCtx(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := upgrader.Upgrade(w, r, nil)
		_, _, _ = c.ReadMessage()
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	ctx, cancel := context.WithCancel(context.Background())

	conn, err := DefaultWSDialer{}.Dial(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	if _, err := conn.ReadMessage(ctx); err == nil {
		t.Fatal("expected error on ctx cancel")
	}
}
