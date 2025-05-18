package facilitator

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rabbitprincess/x402-facilitator/evm"
	"github.com/rabbitprincess/x402-facilitator/evm/eip3009"
	"github.com/rabbitprincess/x402-facilitator/types"
)

var _ Facilitator = (*EVMFacilitator)(nil)

type EVMFacilitator struct {
	scheme string
	client *ethclient.Client
	signer evm.Signer
}

func NewEVMFacilitator(scheme string, url string, privateKeyHex string) (*EVMFacilitator, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	signer := evm.NewRawPrivateSigner(privateKey)

	return &EVMFacilitator{
		scheme: scheme,
		client: client,
		signer: signer,
	}, nil
}

// verification steps:
//   - ✅ verify payload format
//   - ✅ verify payload version
//   - ✅ verify usdc address is correct for the chain
//   - ✅ verify permit signature
//   - ✅ verify deadline
//   - verify nonce is current
//   - ✅ verify client has enough funds to cover paymentRequirements.maxAmountRequired
//   - ✅ verify value in payload is enough to cover paymentRequirements.maxAmountRequired
//   - check min amount is above some threshold we think is reasonable for covering gas
//   - verify resource is not already paid for (next version)
func (t *EVMFacilitator) Verify(payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error) {
	// Step 1: Payload format
	var evmPayload *evm.EVMPayload
	if err := json.Unmarshal([]byte(payload.Payload), evmPayload); err != nil {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: fmt.Sprintf("Invalid payload format: %v", err),
		}, nil
	}

	// Step 2: Scheme verification
	if payload.Scheme != t.scheme || req.Scheme != t.scheme {
		return &types.PaymentVerifyResponse{
			IsValid: false,
			InvalidReason: fmt.Sprintf("Incompatible payload scheme. payload: %s, paymentRequirements: %s, supported: %s",
				payload.Scheme, req.Scheme, t.scheme),
			Payer: evmPayload.Authorization.From.String(),
		}, nil
	}

	// Step 3: Network info
	chainID, ok := evm.GetChainID(payload.Network)
	if !ok {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: "invalid_network",
			Payer:         evmPayload.Authorization.From.String(),
		}, nil
	}
	_ = chainID

	// Step 4: Verify signature (EIP-712)

	// Step 5: Validate payTo

	// Step 6: Deadline check

	// Step 7: TODO: Nonce freshness check (optional in v1)

	// Step 8: Check ERC20 balance

	// Step 9: Check value in permit matches requirement

	// Step 10: TODO: Check minimum payment threshold (e.g. for gas overhead)

	// Step 11: TODO: Check if resource already paid (next version)

	// ✅ All checks passed
	return &types.PaymentVerifyResponse{
		IsValid: true,
		Payer:   evmPayload.Authorization.From.String(),
	}, nil
}

func (t *EVMFacilitator) Settle(payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error) {
	var evmPayload *evm.EVMPayload
	if err := json.Unmarshal([]byte(payload.Payload), evmPayload); err != nil {
		return &types.PaymentSettleResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid payload format: %v", err),
		}, nil
	}

	networkID, ok := evm.GetChainID(req.Network)
	if !ok {
		return &types.PaymentSettleResponse{
			Success: false,
			Error:   "invalid network",
		}, nil
	}

	contractAddress := common.HexToAddress(req.Asset)
	contract, err := eip3009.NewEip3009(contractAddress, t.client)
	if err != nil {
		return nil, fmt.Errorf("contract bind failed: %w", err)
	}
	sig, err := hex.DecodeString(evmPayload.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	r, s, v, err := evm.ParseSignature(sig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse signature: %w", err)
	}

	tx, err := contract.TransferWithAuthorization(
		nil,
		evmPayload.Authorization.From,
		evmPayload.Authorization.To,
		evmPayload.Authorization.Value,
		evmPayload.Authorization.ValidAfter,
		evmPayload.Authorization.ValidBefore,
		evmPayload.Authorization.Nonce,
		v, r, s,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer with authorization %w", err)
	}

	err = t.client.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}
	txHash := tx.Hash().Hex()

	return &types.PaymentSettleResponse{
		Success:   true,
		TxHash:    txHash,
		NetworkId: fmt.Sprintf("%d", networkID),
	}, nil
}

func (t *EVMFacilitator) Supported() []*types.SupportedKind {
	// TODO add scheme
	return nil
}
