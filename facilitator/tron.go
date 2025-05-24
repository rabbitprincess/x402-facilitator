package facilitator

import (
	"context"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/rabbitprincess/x402-facilitator/types"
)

type TronFacilitator struct {
	scheme   types.Scheme
	client   *client.GrpcClient
	feePayer string
}

func NewTronFacilitator(url string, privateKeyHex string) (*TronFacilitator, error) {
	c := client.NewGrpcClient(url)
	if err := c.Start(); err != nil {
		return nil, err
	}

	return &TronFacilitator{
		scheme: types.Tron,
		client: c,
	}, nil
}

func (t *TronFacilitator) Verify(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error) {

	return nil, nil
}

func (t *TronFacilitator) Settle(ctx context.Context, payload *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error) {
	t.client.DelegateResource(t.feePayer)
	t.client.TRC20TransferFrom()

	res, err := t.client.Broadcast()
	if err != nil {
		return nil, err
	}

	return &types.PaymentSettleResponse{
		TxID: res.Txid,
	}, nil
}

func (t *TronFacilitator) Supported() []*types.SupportedKind {
	return []*types.SupportedKind{
		{
			Scheme:  string(types.Tron),
			Network: string(types.Tron),
		},
	}
}
