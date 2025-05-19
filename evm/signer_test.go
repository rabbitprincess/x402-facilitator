package evm

import (
	"testing"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/stretchr/testify/require"
)

func TestSignVerify(t *testing.T) {
	// Generate a random private key
	privKey, err := secp256k1.GeneratePrivateKey()
	require.NoError(t, err)

	// Create a signer using the private key
	signer := NewRawPrivateSigner(privKey.Serialize())

	// Create a message to sign
	message := []byte("Hello, World!")

	// Sign the message
	signature, err := signer(message)
	require.NoError(t, err)

	// Verify the signature
	valid, err := VerifySignature(message, signature)
	require.NoError(t, err)
	require.True(t, valid, "signature verification failed")

}
