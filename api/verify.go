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
// TODO: Implement actual verification logic
func (s *server) Verify(c echo.Context) error {
	ctx := c.Request().Context()

	// validate X-PAYMENT header
	payment := &types.PaymentPayload{}
	paymentHeader := c.Request().Header.Get("X-PAYMENT")
	if paymentHeader == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "X-PAYMENT header is required")
	}
	paymentDecoded, err := base64.StdEncoding.DecodeString(paymentHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed X-PAYMENT header")
	}
	if err := json.Unmarshal(paymentDecoded, payment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed X-PAYMENT header")
	}
	if err := validate.Struct(payment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received invalid X-PAYMENT header")
	}

	// validate payment requirements
	requirement := &types.PaymentRequirements{}
	if err := json.NewDecoder(c.Request().Body).Decode(requirement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received malformed payment requirements")
	}
	if err := validate.Struct(requirement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Received invalid payment requirements")
	}

	verified, err := s.facilitator.Verify(ctx, payment, requirement)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, verified)
}
