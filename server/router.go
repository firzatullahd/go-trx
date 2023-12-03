package server

import (
	a "go-trx/domain/account"
	hc "go-trx/domain/health_check"
	t "go-trx/domain/transaction"

	"github.com/labstack/echo/v4"
)

func InitializeRouter(e *echo.Echo, hcHandler *hc.Handler, a *a.Handler, t *t.Handler) *echo.Echo {
	e.Use(Logging)
	e.GET("/api/health", hcHandler.Ping)
	e.GET("/api/account/:userID", a.GetAccountBalance)
	e.POST("/api/transaction", t.InsertTransaction)

	return e
}
