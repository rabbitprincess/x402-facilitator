package facilitator

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/rabbitprincess/x402-facilitator/evm"
	"github.com/rabbitprincess/x402-facilitator/types"
	"github.com/stretchr/testify/require"
)

const (
	EVMUrl      = "https://sepolia.base.org"
	PrivateKey  = "437770403392c618d3e256200c0c28b679a3c2aa73ce55bbf370445f56d72bd2"
	X402Version = 1
	Network     = "base-sepolia"
	Token       = "USDC"
)

func TestEVMVerify(t *testing.T) {
	facilitator, err := NewEVMFacilitator(types.EVM, EVMUrl, PrivateKey)
	require.NoError(t, err)

	_ = facilitator
}

func TestEVMSettle(t *testing.T) {
	facilitator, err := NewEVMFacilitator(types.EVM, EVMUrl, PrivateKey)
	require.NoError(t, err)

	privKey, err := hex.DecodeString("ef635a97985a15029dea94e7cca2589c688b59bc98b8cd4a90d857048fcbeec1")
	require.NoError(t, err)
	evmPayload, err := evm.GeneratePayload(Network, Token,
		"", "", big.NewInt(10000), evm.NewRawPrivateSigner(privKey))
	require.NoError(t, err)
	evmPayloadJson, err := json.Marshal(evmPayload)
	require.NoError(t, err)

	domainConfig := evm.GetDomainConfig(Network, Token)

	payload := &types.PaymentPayload{
		X402Version: X402Version,
		Scheme:      string(types.EVM),
		Network:     Network,
		Payload:     evmPayloadJson,
	}

	req := &types.PaymentRequirements{
		Scheme:  string(types.EVM),
		Network: Network,
		Asset:   domainConfig.VerifyingContract.String(),
	}

	res, err := facilitator.Settle(payload, req)
	require.NoError(t, err)

	jsonRes, err := json.MarshalIndent(res, "", "\t")
	require.NoError(t, err)
	fmt.Println(string(jsonRes))
}
