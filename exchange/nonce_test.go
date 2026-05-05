package exchange

import "testing"

func TestNonce_Monotonic(t *testing.T) {
	g := &nonceGen{}
	a := g.next()
	b := g.next()
	c := g.next()
	if !(a < b && b < c) {
		t.Fatalf("not monotonic: %d %d %d", a, b, c)
	}
}

func TestNonce_BumpsOnSameMs(t *testing.T) {
	g := &nonceGen{last: 1700000000000}
	first := g.next()
	if first <= 1700000000000 {
		t.Fatalf("expected first > seed, got %d", first)
	}
	second := g.next()
	if second <= first {
		t.Fatalf("expected second > first")
	}
}
