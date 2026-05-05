package ws

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestClient_Trades_LiveServer(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		_, _, _ = c.ReadMessage()
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"channel":"trades","data":[{"coin":"BTC","side":"B","px":"30000","sz":"0.1","time":1,"hash":"0x","tid":1}]}`))
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	c := &Client{URL: url}
	defer c.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	var got []Trade
	sub, err := c.Trades(context.Background(), "BTC", func(t []Trade) {
		got = t
		wg.Done()
	})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Unsubscribe(context.Background())

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("no event")
	}
	if len(got) != 1 || got[0].Coin != "BTC" {
		t.Fatalf("got %+v", got)
	}
}

func TestClient_AllMids_LiveServer(t *testing.T) {
	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, _ := upgrader.Upgrade(w, r, nil)
		defer c.Close()
		_, _, _ = c.ReadMessage()
		_ = c.WriteMessage(websocket.TextMessage, []byte(`{"channel":"allMids","data":{"mids":{"BTC":"30000"}}}`))
		time.Sleep(200 * time.Millisecond)
	}))
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	c := &Client{URL: url}
	defer c.Close()

	got := make(chan AllMidsEvent, 1)
	if _, err := c.AllMids(context.Background(), func(e AllMidsEvent) { got <- e }); err != nil {
		t.Fatal(err)
	}

	select {
	case e := <-got:
		if e.Mids["BTC"] != "30000" {
			t.Fatalf("BTC: %v", e.Mids)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no event")
	}
}
