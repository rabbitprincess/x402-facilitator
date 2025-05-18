package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rabbitprincess/x402-facilitator/api/middleware"
)

// server represents the HTTP server for the API
type server struct {
	*echo.Echo
}

// Ensure server implements http.Handler
var _ http.Handler = (*server)(nil)

// NewServer creates and configures a new API server
func NewServer() *server {
	s := &server{
		Echo: echo.New(),
	}

	// Register middleware
	s.Use(middleware.RequestID())
	s.Use(middleware.Logger())
	s.Use(middleware.ErrorWrapper())
	s.Use(echomiddleware.RecoverWithConfig(echomiddleware.RecoverConfig{
		DisableErrorHandler: true,
	}))
	s.Use(echomiddleware.CORS())

	// Register routes
	s.POST("/verify", s.Verify)
	s.POST("/settle", s.Settle)

	return s
}
