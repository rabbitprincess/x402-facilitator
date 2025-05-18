package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			r := c.Request()
			rid := GetRequestID(r.Context())
			ctx := r.Context()
			ctx = log.With().Str("request_id", rid).Logger().WithContext(ctx)
			c.SetRequest(r.WithContext(ctx))

			start := time.Now()
			err := next(c)

			var evt *zerolog.Event
			if err != nil {
				evt = log.Ctx(ctx).Error().Err(err)
			} else {
				evt = log.Ctx(ctx).Info()
			}

			evt.
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("ip", c.RealIP()).
				Str("user_agent", r.UserAgent()).
				Str("status", strconv.Itoa(c.Response().Status)).
				Dur("latency", time.Since(start)).
				Send()

			return err

		}
	}
}
