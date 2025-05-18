package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorWrapper() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			dummy := &echo.HTTPError{}
			if err != nil && !errors.As(err, &dummy) {
				err = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return err
		}
	}
}
