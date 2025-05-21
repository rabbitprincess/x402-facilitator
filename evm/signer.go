package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func SignEip3009(auth *Authorization, domain *DomainConfig, signer Signer) (string, error) {
	sig, err := signer(HashEip3009(auth, domain))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sig), nil
}

func HashEip3009(auth *Authorization, domain *DomainConfig) []byte {
	domainSeparator := domain.ToMessageHash()
	messageHash := auth.ToMessageHash()

	// Final EIP-712 hash
	var prefix = []byte{0x19, 0x01}
	return Keccak256(
		append(prefix, append(domainSeparator, messageHash...)...),
	)
}

type Signer func(digest []byte) ([]byte, error)

func NewRawPrivateSigner(privateKey []byte) Signer {
	return func(digest []byte) ([]byte, error) {
		privKey := secp256k1.PrivKeyFromBytes(privateKey)

		sig, err := Sign(digest, privKey.ToECDSA())
		if err != nil {
			return nil, err
		}
		return sig, nil
	}
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
