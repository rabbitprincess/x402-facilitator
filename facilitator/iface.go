package facilitator

import (
	"fmt"

	"github.com/rabbitprincess/x402-facilitator/types"
	"github.com/rs/zerolog"
)

type Facilitator interface {
	Verify(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentVerifyResponse, error)
	Settle(payment *types.PaymentPayload, req *types.PaymentRequirements) (*types.PaymentSettleResponse, error)
	Supported() []*types.SupportedKind
}

func NewFacilitator(log *zerolog.Logger, scheme types.Scheme, url string, privateKeyHex string) (Facilitator, error) {
	switch scheme {
	case types.EVM:
		return NewEVMFacilitator(log, scheme, url, privateKeyHex)
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", scheme)
	}
}
