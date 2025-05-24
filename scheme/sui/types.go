package sui

import (
	"crypto/ed25519"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

const (
	PUBKEY_SIZE      = ed25519.PublicKeySize
	ADDRESS_SIZE     = 32
	ADDRESS_HEX_SIZE = 64

	SchemeEd25519 = 0
)

func NewAddressFromPrivateKey(privateKey []byte) ([]byte, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: %d, expected: %d", len(privateKey), ed25519.PrivateKeySize)
	}

	public := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
	msg := make([]byte, PUBKEY_SIZE+1)
	msg[0] = SchemeEd25519
	copy(msg[1:], public)
	hash := blake2b.Sum256(msg)
	address := hash[:ADDRESS_SIZE]
	return address, nil
}

// ExactEvmPayloadAuthorization represents the payload for an exact EVM payment ERC-3009
// authorization EIP-712 typed data message
type SuiPayload struct {
	Signature string `json:"signature"`
	// Transaction Transaction `json:"authorization"`
}
