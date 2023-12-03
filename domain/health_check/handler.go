package health_check

import (
	"go-trx/domain/health_check/service"
	"go-trx/logger"
	"go-trx/utils/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Ping(c echo.Context) error {
	err := h.service.Ping(c.Request().Context())
	if err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return err
	}
	return c.JSON(http.StatusOK, response.SetResponse(http.StatusOK, "success", nil))

}
