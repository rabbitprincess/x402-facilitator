package evm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainInfo struct {
	ChainID        *big.Int
	TokenContracts map[string]*DomainConfig
}

var EVMChains = map[string]ChainInfo{
	"ethereum": {
		ChainID: big.NewInt(1),
		TokenContracts: map[string]*DomainConfig{
			"USDC": {
				Name:              "USD Coin",
				Version:           "2",
				ChainID:           1,
				VerifyingContract: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			},
		},
	},
	"base": {
		ChainID: big.NewInt(8453),
		TokenContracts: map[string]*DomainConfig{
			"USDC": {
				Name:              "USD Coin",
				Version:           "2",
				ChainID:           8453,
				VerifyingContract: common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"),
			},
		},
	},
	"base-sepolia": {
		ChainID: big.NewInt(84532),
		TokenContracts: map[string]*DomainConfig{
			"USDC": {
				Name:              "USDC",
				Version:           "2",
				ChainID:           84532,
				VerifyingContract: common.HexToAddress("0x036CbD53842c5426634e7929541eC2318f3dCF7e"),
			},
		},
	},
}

func GetChainID(chain string) *big.Int {
	chainInfo, ok := EVMChains[chain]
	if !ok {
		return nil
	}
	return chainInfo.ChainID
}

func GetDomainConfig(chain, token string) *DomainConfig {
	chainInfo, ok := EVMChains[chain]
	if !ok {
		return nil
	}
	tokenContract, ok := chainInfo.TokenContracts[token]
	if !ok {
		return nil
	}
	return tokenContract
}
