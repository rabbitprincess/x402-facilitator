package tron

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/rabbitprincess/x402-facilitator/scheme/tron/pb"
	"github.com/rabbitprincess/x402-facilitator/types"
	"google.golang.org/protobuf/proto"
)

func NewRawPrivateSigner(privateKey []byte) types.Signer {
	return func(digest []byte) ([]byte, error) {
		privKey := secp256k1.PrivKeyFromBytes(privateKey)

		sig := ecdsa.SignCompact(privKey, digest, false)
		return sig, nil
	}
}

func HashTransaction(tx *Transaction) ([]byte, error) {

	bytes, err := hex.DecodeString(txStr)
	var trans pb.Transaction
	err = proto.Unmarshal(bytes, &trans)
	if err != nil {
		return pb.Transaction{}, err
	}
	return trans, nil
	rawData, err := proto.Marshal(trans.GetRawData())
	if err != nil {
		return "", err
	}
	s256 := sha256.New()
	s256.Write(rawData)
	hash := s256.Sum(nil)
	return hex.EncodeToString(hash), nil
}
