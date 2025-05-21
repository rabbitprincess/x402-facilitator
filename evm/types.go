package evm

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

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
