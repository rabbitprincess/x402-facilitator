package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var requestIDKey = &struct{}{}

func GetRequestID(ctx context.Context) string {
	rid, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}
	return rid
}

func RequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Request().Header.Get("X-Request-ID")
			if rid == "" {
				rid = uuid.New().String()
			}

			ctx := c.Request().Context()
			// set request id context
			ctx = context.WithValue(ctx, requestIDKey, rid)

			// replace request context
			c.SetRequest(
				c.Request().WithContext(ctx),
			)

			// set response header
			c.Response().Header().Set("X-Request-ID", rid)

			return next(c)
		}
	}
}
