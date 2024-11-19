package handler

import (
	"net/http"

	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"github.com/danielpnjt/speed-engine/internal/usecase/transaction"
	"github.com/labstack/echo/v4"
)

type transactionHandler struct {
	transactionService transaction.Service
}

func NewTransactionHandler() *transactionHandler {
	return &transactionHandler{}
}

func (h *transactionHandler) SetTransactionService(service transaction.Service) *transactionHandler {
	h.transactionService = service
	return h
}

func (h *transactionHandler) Validate() *transactionHandler {
	if h.transactionService == nil {
		panic("transactionService is nil")
	}
	return h
}

func (h *transactionHandler) Generate(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req transaction.GenerateRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}
	res, err := h.transactionService.Generate(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *transactionHandler) Withdraw(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req transaction.WithdrawRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}
	res, err := h.transactionService.Withdraw(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}
