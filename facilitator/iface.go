package facilitator

import (
	"context"
	"fmt"

	"github.com/rabbitprincess/x402-facilitator/types"
)

type Facilitator interface {
	Verify(ctx context.Context, payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error)
	Settle(ctx context.Context, payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error)
	Supported() []*types.SupportedKind
}

func NewFacilitator(scheme types.Scheme, url string, privateKeyHex string) (Facilitator, error) {
	switch scheme {
	case types.EVM:
		return NewEVMFacilitator(url, privateKeyHex)
	case types.Solana:
		return NewSolanaFacilitator(url, privateKeyHex)
	case types.Sui:
		return NewSuiFacilitator(url, privateKeyHex)
	case types.Tron:
		return NewTronFacilitator(url, privateKeyHex)
	default:
		return nil, fmt.Errorf("unsupporsed scheme: %s", scheme)
	}
}
