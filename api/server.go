package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/rabbitprincess/x402-facilitator/api/middleware"
)

var _ http.Handler = (*server)(nil)

func NewServer() *server {
	e := &server{
		Echo: echo.New(),
	}

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.ErrorWrapper())
	e.Use(echomiddleware.RecoverWithConfig(echomiddleware.RecoverConfig{
		DisableErrorHandler: true,
	}))
	e.Use(echomiddleware.CORS())

	e.POST("/verify", e.Verify)
	e.POST("/settle", e.Settle)

	return e
}

type server struct {
	*echo.Echo
}
