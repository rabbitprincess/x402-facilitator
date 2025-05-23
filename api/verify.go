package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rabbitprincess/x402-facilitator/types"
)

var (
	// https://github.com/go-playground/validator?tab=readme-ov-file#special-notes
	// NOTE: If new to using validator it is highly recommended to initialize it using the WithRequiredStructEnabled option which is opt-in to new behaviour that will become the default behaviour in v11+. See documentation for more details.
	validate = validator.New(validator.WithRequiredStructEnabled())
)

// Verify handles verification requests
// Specification:https://github.com/coinbase/x402/tree/3895881f3d6c71fa060076958c8eabc139fcbe5a?tab=readme-ov-file#facilitator-types--interface
func (s *server) Verify(c echo.Context) error {
	ctx := c.Request().Context()

	// validate payment requirements
	requirement := &types.PaymentVerifyRequest{}
	if err := json.NewDecoder(c.Request().Body).Decode(requirement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed payment requirements")
	}
	if err := validate.Struct(requirement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received invalid payment requirements")
	}

	// validate payment payload
	payment := &types.PaymentPayload{}
	paymentDecoded, err := base64.StdEncoding.DecodeString(requirement.PaymentHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed Payment header")
	}
	if err := json.Unmarshal(paymentDecoded, payment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed Payment header")
	}
	if err := validate.Struct(payment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received invalid Payment header")
	}

	verified, err := s.facilitator.Verify(ctx, payment, &requirement.PaymentRequirements)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, verified)
}
