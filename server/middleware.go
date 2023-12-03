package server

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go-trx/logger"
	"time"
)

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		ctx := req.Context()
		correlationID := c.Request().Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = uuid.New().String()
		}
		ctx = context.WithValue(ctx, "X-Correlation-ID", correlationID)
		logger.Info(ctx, "%s: %s", req.Method, req.URL.Path)
		defer func(start time.Time) {
			logger.Info(ctx, "%s: %s took %s", req.Method, req.URL.Path, time.Since(start))
		}(time.Now())
		req = req.WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}
