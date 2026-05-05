package ws

import "testing"

func TestRegistry_AddRemove(t *testing.T) {
	r := newRegistry()
	id1 := r.add("trades:BTC", func([]byte) {})
	id2 := r.add("trades:BTC", func([]byte) {})

	if got := len(r.subscribers("trades:BTC")); got != 2 {
		t.Fatalf("got %d subscribers, want 2", got)
	}

	if last := r.remove("trades:BTC", id1); last {
		t.Fatal("expected lastGone=false, still 1 subscriber")
	}
	subs := r.subscribers("trades:BTC")
	if len(subs) != 1 || subs[0].id != id2 {
		t.Fatalf("after remove: %+v", subs)
	}

	if last := r.remove("trades:BTC", id2); !last {
		t.Fatal("expected lastGone=true after final remove")
	}
	if got := len(r.subscribers("trades:BTC")); got != 0 {
		t.Fatalf("expected 0 subs after lastGone, got %d", got)
	}
}

func TestRegistry_RecordAndSnapshot(t *testing.T) {
	r := newRegistry()
	r.add("trades:BTC", func([]byte) {})
	r.recordParams("trades:BTC", map[string]any{"type": "trades", "coin": "BTC"})
	r.add("allMids", func([]byte) {})
	r.recordParams("allMids", map[string]any{"type": "allMids"})

	snap := r.snapshot()
	if len(snap) != 2 {
		t.Fatalf("snapshot len = %d", len(snap))
	}
	if _, ok := snap["trades:BTC"]; !ok {
		t.Fatal("missing trades:BTC")
	}

	// Snapshot is a copy: mutating it must not affect the registry.
	snap["new"] = "x"
	if _, ok := r.keys["new"]; ok {
		t.Fatal("snapshot leaked into registry")
	}
}
