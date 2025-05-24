package tron

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/base58"
	"github.com/rabbitprincess/x402-facilitator/types"
	"golang.org/x/crypto/sha3"
)

func NewTronPayloed(chain, token, from, to string, value *big.Int, signer types.Signer) (*TronPayload, error) {
	authorization := NewAuthorization(from, to, value)
	domain := GetDomainConfig(chain, token)
	if domain == nil {
		return nil, fmt.Errorf("domain config not found for chain %s and token %s", chain, token)
	}
	signature, err := SignEip3009(authorization, domain, signer)
	if err != nil {
		return nil, err
	}
	return &TronPayload{
		Signature:     signature,
		Authorization: authorization,
	}, nil
}

type TronPayload struct {
	Signature   string       `json:"signature"`
	Transaction *Transaction `json:"transaction"`
}

const (
	Network = 0x41
)

type Address []byte

func GetAddressByPubkey(publicKey *btcec.PublicKey) string {
	pubKey := publicKey.SerializeUncompressed()
	h := sha3.NewLegacyKeccak256()
	h.Write(pubKey[1:])
	hash := h.Sum(nil)[12:]
	return base58.CheckEncode(hash, Network)
}

func GetAddressHash(address string) (Address, error) {
	to, v, err := base58.CheckDecode(address)
	if err != nil {
		return nil, err
	}
	var bs []byte
	bs = append(bs, v)
	bs = append(bs, to...)
	return bs, nil
}

type Transaction struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Nonce    *big.Int `json:"nonce"`
	GasLimit *big.Int `json:"gasLimit"`
	GasPrice *big.Int `json:"gasPrice"`
	Value    float64  `json:"value"`
	Data     []byte   `json:"data"`
	Fee      *big.Int `json:"fee"`
}

type TronTransaction struct {
	Transaction
	RefBlockBytes string   `json:"ref_block_bytes"`
	RefBlockHash  string   `json:"ref_block_hash"`
	RefBlockNum   *big.Int `json:"ref_block_number"`
	Timestamp     *big.Int `json:"timestamp"`
	Expiration    *big.Int `json:"expiration"`
}

type TronTokenTransaction struct {
	TronTransaction
	AssetName       string `json:"asset"`
	ContractAddress string `json:"contractAddress"`
	FeeLimit        int64  `json:"feelimit"`
	Trc             string `json:"trc"`
}
