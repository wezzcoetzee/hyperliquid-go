package exchange

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/sha3"

	"github.com/wezzcoetzee/hyperliquid/internal/msgpack"
)

// ActionHash returns keccak256(msgpack(action) || nonce_be8 || vault_marker[||addr20]).
// Per the TS SDK: vault marker is 0x00 when no vault; otherwise 0x01 followed by
// the 20-byte vault address. Currently does NOT support the optional
// expiresAfter trailer; callers that need it must extend this helper.
func ActionHash(action *msgpack.OrderedMap, nonce uint64, vault []byte) ([]byte, error) {
	if vault != nil && len(vault) != 20 {
		return nil, errors.New("exchange: vault must be 20 bytes or nil")
	}
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
		h.Write(vault)
	}
	return h.Sum(nil), nil
}
