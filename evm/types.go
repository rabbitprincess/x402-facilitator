package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// ExactEvmPayloadAuthorization represents the payload for an exact EVM payment ERC-3009
// authorization EIP-712 typed data message
type EVMPayload struct {
	Signature     string         `json:"signature"`
	Authorization *Authorization `json:"authorization"`
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
	from := padAddress(a.From)
	to := padAddress(a.To)
	value := padBigInt(a.Value)
	validAfter := padBigInt(a.ValidAfter)
	validBefore := padBigInt(a.ValidBefore)
	nonce := a.Nonce[:]

	return Keccak256(
		hash,
		from,
		to,
		value,
		validAfter,
		validBefore,
		nonce,
	)
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
