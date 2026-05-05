package ws

import "sync"

// subscriber pairs a callback with the unique id used to remove it.
type subscriber struct {
	id int64
	cb func([]byte)
}

// registry tracks active subscriptions. It maps a subscription key
// (channel-derived, e.g. "trades:BTC") to all callbacks currently registered
// against it. The original subscribe-payload is stored separately so the
// connection can replay every subscription on reconnect.
type registry struct {
	mu    sync.RWMutex
	next  int64
	byKey map[string][]subscriber
	keys  map[string]any // subscription params keyed by the same key, for replay
}

func newRegistry() *registry {
	return &registry{byKey: map[string][]subscriber{}, keys: map[string]any{}}
}

// add appends a subscriber under key and returns its id. The caller is
// expected to call remove(key, id) on Unsubscribe.
func (r *registry) add(key string, cb func([]byte)) int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	id := r.next
	r.byKey[key] = append(r.byKey[key], subscriber{id: id, cb: cb})
	return id
}

// remove drops the subscriber with id under key. Returns true if no
// subscribers remain for the key (so the caller can send a server-side
// unsubscribe and forget the replay params).
func (r *registry) remove(key string, id int64) (lastGone bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	subs := r.byKey[key]
	out := subs[:0]
	for _, s := range subs {
		if s.id != id {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		delete(r.byKey, key)
		delete(r.keys, key)
		return true
	}
	r.byKey[key] = out
	return false
}

// subscribers returns a snapshot of subscribers under key. Safe to range
// without holding the lock.
func (r *registry) subscribers(key string) []subscriber {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]subscriber, len(r.byKey[key]))
	copy(out, r.byKey[key])
	return out
}

// recordParams stores the subscribe-payload under key so it can be replayed
// on reconnect.
func (r *registry) recordParams(key string, params any) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.keys[key] = params
}

// snapshot returns a copy of all stored params for replay.
func (r *registry) snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]any, len(r.keys))
	for k, v := range r.keys {
		out[k] = v
	}
	return out
}

// allByKey returns a copy of every subscriber list, used by Conn.dispatch
// when fanning out to keyed subscribers.
func (r *registry) allByKey() map[string][]subscriber {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string][]subscriber, len(r.byKey))
	for k, v := range r.byKey {
		cp := make([]subscriber, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}
