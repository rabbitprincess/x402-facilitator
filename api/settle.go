package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Settle handles settlement requests
// TODO: Implement actual settlement logic
func (s *server) Settle(c echo.Context) error {
	ctx := c.Request().Context()

	log.Ctx(ctx).Info().Msg("Settlement request received")

	return c.JSON(http.StatusNotImplemented, map[string]string{
		"message": "settlement not implemented yet",
	})
}
