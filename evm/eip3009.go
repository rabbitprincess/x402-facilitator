package evm

import (
	"encoding/hex"
)

func SignEip3009(auth *Authorization, domain *DomainConfig, signer Signer) (string, error) {
	domainSeparator := domain.ToMessageHash()
	messageHash := auth.ToMessageHash()

	// Final EIP-712 hash
	var prefix = []byte{0x19, 0x01}
	hashBytes := Keccak256(
		append(prefix, append(domainSeparator, messageHash...)...),
	)

	sig, err := signer(hashBytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}
