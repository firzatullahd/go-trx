package server

import (
	"github.com/labstack/echo/v4"
	hc "go-trx/domain/health_check"
)

func InitializeRouter(e *echo.Echo, hcHandler *hc.HealthCheckHandler) *echo.Echo {
	e.Use(Logging)
	e.GET("/api/health", hcHandler.Ping)

	return e
}
