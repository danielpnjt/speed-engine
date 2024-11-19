package handler

import (
	"net/http"

	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"github.com/danielpnjt/speed-engine/internal/usecase/bank"
	"github.com/labstack/echo/v4"
)

type bankHandler struct {
	bankService bank.Service
}

func NewBankHandler() *bankHandler {
	return &bankHandler{}
}

func (h *bankHandler) SetBankService(service bank.Service) *bankHandler {
	h.bankService = service
	return h
}

func (h *bankHandler) Validate() *bankHandler {
	if h.bankService == nil {
		panic("bankService is nil")
	}
	return h
}

func (h *bankHandler) SubmitBank(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req bank.SubmitBank
	if err = utils.Validate(c, &req); err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}
	res, err := h.bankService.SubmitBank(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *bankHandler) FindAll(c echo.Context) (err error) {
	ctx := c.Request().Context()

	res, err := h.bankService.FindAll(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}
