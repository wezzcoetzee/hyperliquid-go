package types

import (
	"encoding/hex"
	"errors"
	"strings"
)

// Address is a 20-byte Ethereum-style address. The zero value is the zero address.
type Address [20]byte

// ParseAddress parses a hex-encoded address with or without "0x" prefix, case-insensitive.
func ParseAddress(s string) (Address, error) {
	var a Address
	s = strings.TrimPrefix(strings.ToLower(s), "0x")
	if len(s) != 40 {
		return a, errors.New("invalid address length")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return a, err
	}
	copy(a[:], b)
	return a, nil
}

// Hex returns the lowercase 0x-prefixed hex representation.
func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

func (a Address) String() string { return a.Hex() }

// Side identifies an order direction. Hyperliquid encodes Buy as "B" and Sell as "A".
type Side string

const (
	Buy  Side = "B"
	Sell Side = "A"
)

func (s Side) String() string { return string(s) }

// Tif is a time-in-force qualifier for limit orders.
type Tif string

const (
	TifGtc Tif = "Gtc"
	TifIoc Tif = "Ioc"
	TifAlo Tif = "Alo"
)

// Cloid is a client-supplied 16-byte order id, hex-encoded with 0x prefix.
type Cloid string
