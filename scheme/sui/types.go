package sui

import (
	"crypto/ed25519"
	"fmt"

	"golang.org/x/crypto/blake2b"

	"github.com/mr-tron/base58"
	"github.com/rabbitprincess/x402-facilitator/types"
)

const (
	PUBKEY_SIZE      = ed25519.PublicKeySize
	ADDRESS_SIZE     = 32
	ADDRESS_HEX_SIZE = 64

	SchemeEd25519 = 0
)

func NewSuiPayload(sign types.Signer) (*SuiPayload, error) {
	digest := make([]byte, 1+PUBKEY_SIZE)
	signature, err := sign(digest)
	if err != nil {
		return nil, fmt.Errorf("failed to sign payload: %w", err)
	}

	return &SuiPayload{
		Signature: base58.Encode(signature),
		// Transaction: Transaction{},
	}, nil
}

type SuiPayload struct {
	Signature   string         `json:"signature"`
	Transaction SuiTransaction `json:"transaction"`
}

func NewSponsoredTransaction(from, to, sponser string) *SuiTransaction {
	return &SuiTransaction{}
}

type SuiTransaction struct {
}

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
