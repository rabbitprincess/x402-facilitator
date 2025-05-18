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
			"USDC": common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
		},
	},
	"base": {
		ChainID: big.NewInt(8453),
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0xd9aaEC86B65d86f6a7b5B1b0c42FFA531710b6CA"),
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
