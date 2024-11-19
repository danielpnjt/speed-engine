package admin

import (
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

// * Requests
type ()

// * Responses
type (
	DefaultResponse struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Errors  []string    `json:"errors"`
	}

	CreateResponse struct {
		constants.DefaultResponse
	}

	FindAllResponse struct {
		constants.DefaultResponse
	}
)
