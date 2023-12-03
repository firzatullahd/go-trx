package health_check

import (
	"github.com/labstack/echo/v4"
	"go-trx/domain/health_check/service"
	"go-trx/logger"
	"go-trx/utils/response"
	"net/http"
)

type HealthCheckHandler struct {
	service service.Service
}

func NewHealthCheckHandler(service service.Service) *HealthCheckHandler {
	return &HealthCheckHandler{
		service: service,
	}
}

func (h *HealthCheckHandler) Ping(c echo.Context) error {
	err := h.service.Ping(c.Request().Context())
	if err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return err
	}
	return c.JSON(http.StatusOK, response.SetResponse(http.StatusOK, "success", nil))

}
