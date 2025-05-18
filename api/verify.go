package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// Verify handles verification requests
// TODO: Implement actual verification logic
func (s *server) Verify(c echo.Context) error {
	return fmt.Errorf("verification not implemented yet")
}
