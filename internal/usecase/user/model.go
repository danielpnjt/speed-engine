package user

import (
	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

// * Requests
type (
	RegisterRequest struct {
		Username        string `json:"username" validate:"required"`
		Email           string `json:"email" validate:"required"`
		Name            string `json:"name" validate:"required"`
		Password        string `json:"password" validate:"required"`
		ConfirmPassword string `json:"confirmPassword" validate:"required"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	FindAllRequest struct {
		constants.PaginationRequest
	}
)

// * Responses
type (
	DefaultResponse struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Errors  []string    `json:"errors"`
	}

	LoginResponseData struct {
		Token    string `json:"token"`
		ExpireAt int64  `json:"expireAt"`
	}

	GenerateAccessTokenResponseData struct {
		AccessToken string `json:"accessToken"`
		TokenType   string `json:"tokenType"`
		ExpiresAt   int64  `json:"expiresAt"`
	}

	GenerateAccessTokenResponse struct {
		DefaultResponse
		Data GenerateAccessTokenResponseData `json:"data"`
	}

	CreateResponse struct {
		constants.DefaultResponse
	}

	FindAllResponse struct {
		constants.DefaultResponse
	}
)
