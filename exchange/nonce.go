package exchange

import (
	"sync"
	"time"
)

// nonceGen produces strictly-monotonic millisecond timestamps. Hyperliquid
// requires nonces to be strictly increasing per signer; in burst conditions
// where two calls arrive within the same millisecond, we increment the
// previous nonce by one rather than reusing.
type nonceGen struct {
	mu   sync.Mutex
	last uint64
}

func (g *nonceGen) next() uint64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := uint64(time.Now().UnixMilli())
	if now <= g.last {
		now = g.last + 1
	}
	g.last = now
	return now
}
