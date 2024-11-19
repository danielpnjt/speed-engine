package transaction

import (
	"time"

	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

// * Requests
type (
	GenerateRequest struct {
		Amount float64 `json:"amount" validate:"required"`
	}

	TopUpRequest struct {
		ReferenceID string `json:"referenceId" validate:"required"`
	}

	WithdrawRequest struct {
		BankID int `json:"bankId" validate:"required"`
		Amount int `json:"amount" validate:"required"`
	}
)

// * Responses
type (
	GenerateResponseData struct {
		VA          string    `json:"va"`
		ReferenceID string    `json:"referenceId"`
		ExpiredAt   time.Time `json:"expiredAt"`
	}
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
