package evm

import "github.com/ethereum/go-ethereum/common"

type ChainInfo struct {
	ChainID        uint64
	TokenContracts map[string]common.Address // 예: "USDC" -> 0xA0b...
}

var EVMChains = map[string]ChainInfo{
	"ethereum": {
		ChainID: 1,
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"),
			"DAI":  common.HexToAddress("0x6b175474e89094c44da98b954eedeac495271d0f"),
		},
	},
	"base": {
		ChainID: 8453,
		TokenContracts: map[string]common.Address{
			"USDC": common.HexToAddress("0xd9aaEC86B65d86f6a7b5B1b0c42FFA531710b6CA"),
		},
	},
}

func GetChainID(chain string) (uint64, bool) {
	chainInfo, ok := EVMChains[chain]
	if !ok {
		return 0, false
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
