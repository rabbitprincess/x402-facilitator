package evm

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil)
}

func padAddress(addr Address) []byte {
	return append(make([]byte, 12), addr[:]...)
}

func padBigInt(n *big.Int) []byte {
	return leftPadBytes(n.Bytes(), 32)
}

func leftPadBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

// Utility to convert hex string to Address
func ParseAddress(hexStr string) (Address, error) {
	var a Address
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	if err != nil || len(b) != 20 {
		return a, errors.New("invalid address")
	}
	copy(a[:], b)
	return a, nil
}
