package facilitator

import "github.com/rabbitprincess/x402-facilitator/types"

type Facilitator interface {
	Verify(payment types.PaymentPayload, req types.PaymentRequirements) (types.PaymentVerifyResponse, error)
	Settle(payment types.PaymentPayload, req types.PaymentRequirements) (types.PaymentSettleResponse, error)
	Supported() []types.SupportedKind
}
