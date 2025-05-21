package facilitator

import (
	"fmt"

	"github.com/rabbitprincess/x402-facilitator/types"
)

type Facilitator interface {
	Verify(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error)
	Settle(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error)
	Supported() []*types.SupportedKind
}

func NewFacilitator(scheme types.Scheme, url string, privateKeyHex string) (Facilitator, error) {
	switch scheme {
	case types.EVM:
		return NewEVMFacilitator(url, privateKeyHex)
	default:
		return nil, fmt.Errorf("unsupporsed scheme: %s", scheme)
	}
}
