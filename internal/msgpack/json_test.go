package msgpack

import (
	"strings"
	"testing"
)

func TestMarshalOrderedJSON_PreservesOrder(t *testing.T) {
	m := NewOrderedMap()
	m.Set("type", "order")
	m.Set("orders", []any{})
	m.Set("grouping", "na")

	got, err := MarshalOrderedJSON(m)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"type":"order","orders":[],"grouping":"na"}`
	if string(got) != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestMarshalOrderedJSON_Nested(t *testing.T) {
	inner := NewOrderedMap()
	inner.Set("a", uint64(1))
	inner.Set("b", uint64(2))

	outer := NewOrderedMap()
	outer.Set("orders", []any{inner})
	outer.Set("flag", true)

	got, _ := MarshalOrderedJSON(outer)
	if !strings.Contains(string(got), `"orders":[{"a":1,"b":2}]`) {
		t.Fatalf("nested order broken: %s", got)
	}
}
