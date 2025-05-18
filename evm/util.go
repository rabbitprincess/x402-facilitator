package evm

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

func GenerateEIP3009Nonce() (string, error) {
	var nonce [32]byte
	_, err := rand.Read(nonce[:])
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}
	nonceHex := hex.EncodeToString(nonce[:])
	return nonceHex, nil
}

func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil)
}

func padAddress(addr common.Address) []byte {
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
func ParseAddress(hexStr string) (common.Address, error) {
	var a common.Address
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	if err != nil || len(b) != 20 {
		return a, errors.New("invalid address")
	}
	copy(a[:], b)
	return a, nil
}

func ParseSignature(sig []byte) (r, s [32]byte, v uint8, err error) {
	if len(sig) != 65 {
		return r, s, 0, errors.New("invalid signature length")
	}
	copy(r[:], sig[0:32])
	copy(s[:], sig[32:64])
	v = sig[64]
	if v < 27 { // normalize v to 27 or 28
		v += 27
	}
	if v != 27 && v != 28 {
		return r, s, 0, errors.New("invalid signature v value")
	}
	return r, s, v, nil
}
