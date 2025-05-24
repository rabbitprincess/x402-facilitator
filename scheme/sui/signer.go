package sui

import (
	"crypto/ed25519"
	"fmt"

	"github.com/rabbitprincess/x402-facilitator/types"
)

func NewRawPrivateSigner(privateKey []byte) types.Signer {
	return func(digest []byte) ([]byte, error) {
		if len(privateKey) != ed25519.PrivateKeySize {
			return nil, fmt.Errorf("invalid private key length: %d, expected: %d", len(privateKey), ed25519.PrivateKeySize)
		}
		signature := ed25519.Sign(privateKey, digest)
		return signature, nil
	}
}
