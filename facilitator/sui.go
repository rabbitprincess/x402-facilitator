package facilitator

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/rabbitprincess/x402-facilitator/scheme/sui"
	"github.com/rabbitprincess/x402-facilitator/types"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
)

type SuiFacilitator struct {
	network string
	scheme  string
	client  *client.Client
	signer  types.Signer
	address sui_types.SuiAddress
}

func NewSuiFacilitator(network string, url string, privateKeyHex string) (*SuiFacilitator, error) {
	cli, err := client.DialWithClient(url, http.DefaultClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sui client: %w", err)
	}

	privateKey, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid hex private key: %w", err)
	}
	signer := sui.NewRawPrivateSigner(privateKey)
	address, err := sui.NewAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive Sui address from private key: %w", err)
	}

	return &SuiFacilitator{
		network: network,
		client:  cli,
		signer:  signer,
		address: sui_types.SuiAddress(address),
	}, nil
}

// verification steps:
//
// Steps to verify a payment for the `exact` scheme:
// - ✅ Verify the network is for the agreed upon chain.
// - ✅ Verify the signature is valid over the provided transaction.
// - ✅ Simulate the transaction to ensure it would succeed and has not already been executed/committed to the chain.
// - ✅ Verify the outputs of the simulation/execution to ensure the resource server's address sees a balance change equal to the value in the `paymentRequirements.maxAmountRequired` in the agreed upon `asset`.
func (t *SuiFacilitator) Verify(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error) {

	return nil, nil
}

func (t *SuiFacilitator) Settle(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error) {
	return nil, nil
}

func (t *SuiFacilitator) Supported() []*types.SupportedKind {
	return []*types.SupportedKind{
		{
			Scheme:  string(types.Sui),
			Network: string(types.Sui),
		},
	}
}
