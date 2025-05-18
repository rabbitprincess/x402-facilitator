package evm

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1"
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
		v := []byte{sigCompact[0] + 27} // v = recoveryID + 27

		return append(append(r, s...), v...), nil
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
