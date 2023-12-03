package transaction

import (
	"go-trx/domain/transaction/model"
	"go-trx/domain/transaction/service"
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

func (h *Handler) InsertTransaction(c echo.Context) error {

	var payload model.NewTransaction
	if err := c.Bind(&payload); err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return err
	}

	err := h.service.InsertTransaction(c.Request().Context(), payload)
	if err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return c.JSON(http.StatusBadRequest, response.SetResponse(http.StatusBadRequest, err.Error(), nil))
	}
	return c.JSON(http.StatusCreated, response.SetResponse(http.StatusCreated, "success", nil))

}
