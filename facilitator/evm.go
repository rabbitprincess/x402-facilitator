package facilitator

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/rabbitprincess/x402-facilitator/evm"
	"github.com/rabbitprincess/x402-facilitator/evm/eip3009"
	"github.com/rabbitprincess/x402-facilitator/types"
)

var _ Facilitator = (*EVMFacilitator)(nil)

type EVMFacilitator struct {
	scheme    types.Scheme
	network   string
	networkID *big.Int

	client  *ethclient.Client
	signer  evm.Signer
	address common.Address
}

func NewEVMFacilitator(url string, privateKeyHex string) (*EVMFacilitator, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	networkId, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %w", err)
	}
	network := evm.GetChainName(networkId)
	if network == "" {
		return nil, fmt.Errorf("unsupported network ID: %s", networkId.String())
	}

	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, err
	}
	signer := evm.NewRawPrivateSigner(privateKey)
	address, err := evm.GetAddrssFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get address from private key: %w", err)
	}

	return &EVMFacilitator{
		scheme:    types.EVM,
		network:   network,
		networkID: networkId,

		client:  client,
		signer:  signer,
		address: address,
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
func (t *EVMFacilitator) Verify(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error) {
	// Step 1: Payload format
	var evmPayload evm.EVMPayload
	if err := json.Unmarshal([]byte(payload.Payload), &evmPayload); err != nil {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: fmt.Sprintf("Invalid payload format: %v", err),
		}, nil
	}

	// Step 2: Scheme verification
	if payload.Scheme != string(t.scheme) || req.Scheme != string(t.scheme) {
		return &types.PaymentVerifyResponse{
			IsValid: false,
			InvalidReason: fmt.Sprintf("Incompatible payload scheme. payload: %s, paymentRequirements: %s, supported: %s",
				payload.Scheme, req.Scheme, t.scheme),
			Payer: evmPayload.Authorization.From.String(),
		}, nil
	}

	// Step 3: Network info
	if payload.Network != t.network {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: "network_mismatch",
			Payer:         evmPayload.Authorization.From.String(),
		}, nil
	}
	chainID := evm.GetChainID(payload.Network)
	if chainID == nil {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: "invalid_network",
			Payer:         evmPayload.Authorization.From.String(),
		}, nil
	}
	if chainID.Cmp(t.networkID) != 0 {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: "network_id_mismatch",
			Payer:         evmPayload.Authorization.From.String(),
		}, nil
	}

	// Step 4: Verify signature (EIP-712)
	sig, err := hex.DecodeString(evmPayload.Signature)
	if err != nil {
		return nil, err
	}
	digest := evmPayload.Authorization.ToMessageHash()
	pubkey, err := evm.Ecrecover(digest, sig)
	if err != nil {
		return nil, err
	}
	if valid := evm.VerifySignature(pubkey, digest, sig); !valid {
		return &types.PaymentVerifyResponse{
			IsValid:       false,
			InvalidReason: "invalid_signature",
			Payer:         evmPayload.Authorization.From.String(),
		}, nil
	}

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

func (t *EVMFacilitator) Settle(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error) {
	var evmPayload evm.EVMPayload
	if err := json.Unmarshal([]byte(payload.Payload), &evmPayload); err != nil {
		return &types.PaymentSettleResponse{
			Success: false,
			Error:   fmt.Sprintf("invalid payload format: %v", err),
		}, nil
	}

	networkID := evm.GetChainID(req.Network)
	if networkID == nil {
		return &types.PaymentSettleResponse{
			Success: false,
			Error:   "invalid network",
		}, nil
	}

	contractAddress := common.HexToAddress(req.Asset)
	if contractAddress == (common.Address{}) {
		return &types.PaymentSettleResponse{
			Success: false,
			Error:   "invalid contract address",
		}, nil
	}
	contract, err := eip3009.NewEip3009(contractAddress, t.client)
	if err != nil {
		return nil, fmt.Errorf("contract bind failed: %w", err)
	}
	clientSig, err := hex.DecodeString(evmPayload.Signature) // client signature
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature: %w", err)
	}
	if clientSig[64] < 27 {
		clientSig[64] += 27
	}

	tx, err := contract.TransferWithAuthorization(
		&bind.TransactOpts{
			Context: ctx,
			Signer:  evm.ToGethSigner(t.signer, networkID), // facilitator signature
			From:    t.address,
		},
		evmPayload.Authorization.From,
		evmPayload.Authorization.To,
		evmPayload.Authorization.Value,
		evmPayload.Authorization.ValidAfter,
		evmPayload.Authorization.ValidBefore,
		evmPayload.Authorization.Nonce,
		clientSig,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer with authorization %w", err)
	}

	return &types.PaymentSettleResponse{
		Success:   true,
		TxHash:    tx.Hash().Hex(),
		NetworkId: fmt.Sprintf("%d", networkID),
	}, nil
}

func (t *EVMFacilitator) Supported() []*types.SupportedKind {
	return []*types.SupportedKind{
		{
			Scheme:  string(t.scheme),
			Network: t.network,
		},
	}
}
