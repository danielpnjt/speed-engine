package handler

import (
	"github.com/danielpnjt/speed-engine/internal/usecase/admin"
	"github.com/labstack/echo/v4"
)

type adminHandler struct {
	adminService admin.Service
}

func NewAdminHandler() *adminHandler {
	return &adminHandler{}
}

func (h *adminHandler) SetAdminService(service admin.Service) *adminHandler {
	h.adminService = service
	return h
}

func (h *adminHandler) Validate() *adminHandler {
	if h.adminService == nil {
		panic("adminService is nil")
	}
	return h
}

func (h *adminHandler) GetAll(c echo.Context) (err error) {
	return
}

func (h *adminHandler) GetByUserID(c echo.Context) (err error) {
	return
}
