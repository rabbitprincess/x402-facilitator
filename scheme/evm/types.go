package evm

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"

	"github.com/rabbitprincess/x402-facilitator/types"
)

func NewEVMPayload(chain, token, from, to string, value *big.Int, signer types.Signer) (*EVMPayload, error) {
	authorization := NewAuthorization(from, to, value)
	domain := GetDomainConfig(chain, token)
	if domain == nil {
		return nil, fmt.Errorf("domain config not found for chain %s and token %s", chain, token)
	}
	signature, err := SignEip3009(authorization, domain, signer)
	if err != nil {
		return nil, err
	}
	return &EVMPayload{
		Signature:     signature,
		Authorization: authorization,
	}, nil

}

// ExactEvmPayloadAuthorization represents the payload for an exact EVM payment ERC-3009
// authorization EIP-712 typed data message
type EVMPayload struct {
	Signature     string         `json:"signature"`
	Authorization *Authorization `json:"authorization"`
}

func NewAuthorization(from, to string, value *big.Int) *Authorization {
	now := time.Now().Unix()
	authorization := &Authorization{
		From:        common.HexToAddress(from),
		To:          common.HexToAddress(to),
		Value:       value,
		ValidAfter:  big.NewInt(0),
		ValidBefore: big.NewInt(now + 3600), // 1 hour
		Nonce:       GenerateEIP3009Nonce(),
	}
	return authorization
}

type Authorization struct {
	From        common.Address
	To          common.Address
	Value       *big.Int
	ValidAfter  *big.Int
	ValidBefore *big.Int
	Nonce       [32]byte
}

func (a Authorization) ToMessageHash() []byte {
	hash := Keccak256([]byte(
		"TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"),
	)

	encoded := bytes.Join([][]byte{
		hash,
		padAddress(a.From),
		padAddress(a.To),
		padBigInt(a.Value),
		padBigInt(a.ValidAfter),
		padBigInt(a.ValidBefore),
		a.Nonce[:], // already 32 bytes
	}, nil)
	return Keccak256(encoded)
}

func NewDomainConfig(name, version string, chainID int64, verifyingContract string) *DomainConfig {
	return &DomainConfig{
		Name:              name,
		Version:           version,
		ChainID:           chainID,
		VerifyingContract: common.HexToAddress(verifyingContract),
	}
}

type DomainConfig struct {
	Name              string
	Version           string
	ChainID           int64
	VerifyingContract common.Address
}

func (d DomainConfig) ToMessageHash() []byte {
	hash := Keccak256([]byte(
		"EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"),
	)
	nameHash := Keccak256([]byte(d.Name))
	versionHash := Keccak256([]byte(d.Version))
	chainID := padBigInt(big.NewInt(d.ChainID))
	contract := padAddress(d.VerifyingContract)

	return Keccak256(
		hash,
		nameHash,
		versionHash,
		chainID,
		contract,
	)
}

func GetAddrssFromPrivateKey(privateKey []byte) (common.Address, error) {
	if len(privateKey) != 32 {
		return common.Address{}, errors.New("invalid private key length")
	}

	privKey := secp256k1.PrivKeyFromBytes(privateKey)
	uncompressed := privKey.PubKey().SerializeUncompressed()
	address := common.BytesToAddress(Keccak256(uncompressed[1:])[12:])

	return address, nil
}

func GenerateEIP3009Nonce() [32]byte {
	var nonce [32]byte
	rand.Read(nonce[:])
	return nonce
}

func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil)
}

func padAddress(addr common.Address) []byte {
	return append(make([]byte, 12), addr[:]...)
}

func padBigInt(n *big.Int) []byte {
	return leftPadBytes(n.Bytes(), 32)
}

func leftPadBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b
	}
	padded := make([]byte, size)
	copy(padded[size-len(b):], b)
	return padded
}

// Utility to convert hex string to Address
func ParseAddress(hexStr string) (common.Address, error) {
	var a common.Address
	hexStr = strings.TrimPrefix(hexStr, "0x")
	b, err := hex.DecodeString(hexStr)
	if err != nil || len(b) != 20 {
		return a, errors.New("invalid address")
	}
	copy(a[:], b)
	return a, nil
}

func ParseSignature(sig []byte) (r, s [32]byte, v uint8, err error) {
	if len(sig) != 65 {
		return r, s, 0, errors.New("invalid signature length")
	}
	copy(r[:], sig[0:32])
	copy(s[:], sig[32:64])
	v = sig[64]
	if v < 27 { // normalize v to 27 or 28
		v += 27
	}
	if v != 27 && v != 28 {
		return r, s, 0, errors.New("invalid signature v value")
	}
	return r, s, v, nil
}
