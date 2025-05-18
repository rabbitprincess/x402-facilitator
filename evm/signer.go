package evm

import (
	"errors"

	"github.com/decred/dcrd/dcrec/secp256k1"
)

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
