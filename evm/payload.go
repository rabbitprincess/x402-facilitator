package evm

import (
	"fmt"
	"math/big"
)

func GeneratePayload(chain, token, from, to string, value *big.Int, signer Signer) (*EVMPayload, error) {
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
