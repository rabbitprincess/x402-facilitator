package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainInfo struct {
	ChainID        *big.Int
	TokenContracts map[string]common.Address
}

var EVMChains = map[string]ChainInfo{
	"ethereum": {
		ChainID: big.NewInt(1),
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
		},
	},
	"base": {
		ChainID: big.NewInt(8453),
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
		},
	},
	"base-sepolia": {
		ChainID: big.NewInt(167002),
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
		},
	},
}

func GetChainID(chain string) (*big.Int, bool) {
	chainInfo, ok := EVMChains[chain]
	if !ok {
		return nil, false
	}
	return chainInfo.ChainID, true
}

func GetTokenContract(chain, token string) (common.Address, bool) {
	chainInfo, ok := EVMChains[chain]
	if !ok {
		return common.Address{}, false
	}
	tokenContract, ok := chainInfo.TokenContracts[token]
	if !ok {
		return common.Address{}, false
	}
	return tokenContract, true
}
