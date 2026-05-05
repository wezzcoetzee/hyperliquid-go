package exchange

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
)

// ActionHash returns the L1 action hash used by Hyperliquid's signing scheme.
//
//	keccak256( msgpack(action) || nonce_be8
//	           || (vault ? 0x01 || addr20 : 0x00)
//	           || (expiresAfter ? 0x00 || expiresAfter_be8 : ε) )
//
// The expiresAfter trailer is OMITTED entirely when nil — not zero-padded —
// matching the TS SDK's createL1ActionHash helper.
func ActionHash(action *msgpack.OrderedMap, nonce uint64, vault *[20]byte, expiresAfter *uint64) ([]byte, error) {
	encoded, err := msgpack.Encode(action)
	if err != nil {
		return nil, err
	}

	var nonceBytes [8]byte
	binary.BigEndian.PutUint64(nonceBytes[:], nonce)

	h := sha3.NewLegacyKeccak256()
	h.Write(encoded)
	h.Write(nonceBytes[:])
	if vault == nil {
		h.Write([]byte{0x00})
	} else {
		h.Write([]byte{0x01})
		h.Write(vault[:])
	}
	if expiresAfter != nil {
		var ea [8]byte
		binary.BigEndian.PutUint64(ea[:], *expiresAfter)
		h.Write([]byte{0x00})
		h.Write(ea[:])
	}
	return h.Sum(nil), nil
}
