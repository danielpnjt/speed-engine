package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/danielpnjt/speed-engine/internal/pkg/utils"
	"github.com/danielpnjt/speed-engine/internal/usecase/user"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) SetUserService(service user.Service) *userHandler {
	h.userService = service
	return h
}

func (h *userHandler) Validate() *userHandler {
	if h.userService == nil {
		panic("userService is nil")
	}
	return h
}

func (h *userHandler) Register(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req user.RegisterRequest
	if err = utils.Validate(c, &req); err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}
	res, err := h.userService.Register(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) Login(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req user.LoginRequest
	if err = utils.Validate(c, &req); err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}
	res, err := h.userService.Login(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) Logout(c echo.Context) (err error) {
	ctx := c.Request().Context()

	res, err := h.userService.Logout(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) GetDetail(c echo.Context) (err error) {
	ctx := c.Request().Context()

	userID, err := cast.ToIntE(c.Param("userID"))
	if err != nil {
		slog.Error("failed convert id in param into int", err)
		err = errors.New("invalid request format")
		return
	}

	res, err := h.userService.GetDetail(ctx, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) GetAll(c echo.Context) (err error) {
	ctx := c.Request().Context()

	var req user.FindAllRequest
	if err = utils.Validate(c, &req); err != nil {
		return
	}
	res, err := h.userService.GetAll(ctx, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *userHandler) GetDetailPlayer(c echo.Context) (err error) {
	ctx := c.Request().Context()

	res, err := h.userService.GetDetailPlayer(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, res)
}
