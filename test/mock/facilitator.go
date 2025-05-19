package mock

import (
	"github.com/rabbitprincess/x402-facilitator/types"
)

// Facilitator implements the Facilitator interface for testing
type Facilitator struct{}

func (m *Facilitator) Verify(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error) {
	return &types.PaymentVerifyResponse{
		IsValid: true,
	}, nil
}

func (m *Facilitator) Settle(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error) {
	return &types.PaymentSettleResponse{
		Success: true,
	}, nil
}

func (m *Facilitator) Supported() []*types.SupportedKind {
	return []*types.SupportedKind{}
}
