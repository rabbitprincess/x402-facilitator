package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v2"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
		return "s", err
	}

	return hex.EncodeToString(sig), nil
}

type Signer func(digest []byte) ([]byte, error)

func NewRawPrivateSigner(privateKey []byte) Signer {
	return func(digest []byte) ([]byte, error) {
		privKey, _ := secp256k1.PrivKeyFromBytes(privateKey)

		sig, err := privKey.Sign(digest)
		if err != nil {
			return nil, err
		}

		r := leftPadBytes(sig.R.Bytes(), 32)
		s := leftPadBytes(sig.S.Bytes(), 32)
		sigCompact, err := secp256k1.SignCompact(privKey, digest, false)
		if err != nil || len(sigCompact) != 65 {
			return nil, errors.New("failed to compute recovery ID")
		}
		v := []byte{sigCompact[0]}

		return append(append(r, s...), v...), nil
	}
}

func VerifySignature(digest []byte, signature []byte) (bool, error) {
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	v := signature[64]
	if v >= 35 { // Adjust v for recovery (should be 27 or 28)
		v = byte((v-35)%2 + 27)
	}
	if v != 27 && v != 28 {
		return false, errors.New("invalid recovery ID in signature")
	}
	signature[64] = v - 27 // convert to 0 or 1 for recovery
	compactSig := append([]byte{v}, signature[0:64]...)
	pubKey, _, err := secp256k1.RecoverCompact(compactSig, digest)
	if err != nil {
		return false, fmt.Errorf("failed to recover public key: %w", err)
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])

	// verify the signature
	valid := ecdsa.Verify(pubKey.ToECDSA(), digest, r, s)
	return valid, nil
}

func ToGethSigner(signer Signer, chainID *big.Int) bind.SignerFn {
	return func(_ common.Address, tx *types.Transaction) (*types.Transaction, error) {
		signerObj := types.LatestSignerForChainID(chainID)
		digest := signerObj.Hash(tx).Bytes()

		sig, err := signer(digest)
		if err != nil {
			return nil, err
		} else if len(sig) != 65 {
			return nil, fmt.Errorf("invalid signature length: %d", len(sig))
		}

		return tx.WithSignature(signerObj, sig)
	}
}
