package sui

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/mr-tron/base58"
	"golang.org/x/crypto/blake2b"

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

func Hash(txBytes string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		return "", err
	}
	data := make([]byte, len("TransactionData::")+len(b))
	copy(data, "TransactionData::")
	copy(data[len("TransactionData::"):], b)
	hash := blake2b.Sum256(data)
	return base58.Encode(hash[:]), nil
}
