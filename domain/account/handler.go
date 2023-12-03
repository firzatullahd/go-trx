package account

import (
	"go-trx/domain/account/service"
	"go-trx/logger"
	"go-trx/utils/response"
	"net/http"
	"strconv"

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

func (h *Handler) GetAccountBalance(c echo.Context) error {

	strUserID := c.Param("userID")
	userID, err := strconv.ParseUint(strUserID, 10, 64)
	if err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return err
	}

	account, err := h.service.AccountBalance(c.Request().Context(), userID)
	if err != nil {
		logger.Error(c.Request().Context(), err.Error())
		return c.JSON(http.StatusBadRequest, response.SetResponse(http.StatusBadRequest, err.Error(), nil))
	}
	return c.JSON(http.StatusOK, response.SetResponse(http.StatusOK, "success", account))

}
