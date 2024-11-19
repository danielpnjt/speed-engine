package payment

import "time"

type (
	CreateVARequest struct {
		ExternalID     string  `json:"external_id" validate:"required"`
		ExpectedAmount float64 `json:"expected_amount,omitempty"`
	}

	TopUpRequest struct {
		ExternalID string `json:"external_id" validate:"required"`
	}

	WithdrawRequest struct {
		ExternalID        string `json:"external_id" validate:"required"`
		BankCode          string `json:"bank_code" validate:"required"`
		AccountHolderName string `json:"account_holder_name" validate:"required"`
		AccountNumber     string `json:"account_number" validate:"required"`
		Amount            int    `json:"amount,required"`
		Description       string `json:"description,omitempty"`
	}
)

type (
	CreateVAResponse struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Data    CreateVAResponseData `json:"data"`
	}

	CreateVAResponseData struct {
		OwnerID         string     `json:"owner_id"`
		ExternalID      string     `json:"external_id"`
		BankCode        string     `json:"bank_code"`
		MerchantCode    string     `json:"merchant_code"`
		Name            string     `json:"name"`
		AccountNumber   string     `json:"account_number"`
		IsClosed        *bool      `json:"is_closed"`
		ID              string     `json:"id"`
		IsSingleUse     *bool      `json:"is_single_use"`
		Status          string     `json:"status"`
		Currency        string     `json:"currency"`
		ExpirationDate  *time.Time `json:"expiration_date"`
		SuggestedAmount float64    `json:"suggested_amount,omitempty"`
		ExpectedAmount  float64    `json:"expected_amount,omitempty"`
		Description     string     `json:"description,omitempty"`
	}

	TopUpResponse struct {
		Status  string            `json:"status"`
		Message string            `json:"message"`
		Data    TopUpResponseData `json:"data"`
	}

	TopUpResponseData struct {
		OwnerID         string     `json:"owner_id"`
		ExternalID      string     `json:"external_id"`
		BankCode        string     `json:"bank_code"`
		MerchantCode    string     `json:"merchant_code"`
		Name            string     `json:"name"`
		AccountNumber   string     `json:"account_number"`
		IsClosed        *bool      `json:"is_closed"`
		ID              string     `json:"id"`
		IsSingleUse     *bool      `json:"is_single_use"`
		Status          string     `json:"status"`
		Currency        string     `json:"currency"`
		ExpirationDate  *time.Time `json:"expiration_date"`
		SuggestedAmount float64    `json:"suggested_amount,omitempty"`
		ExpectedAmount  float64    `json:"expected_amount,omitempty"`
		Description     string     `json:"description,omitempty"`
	}

	WithdrawResponse struct {
		Status  string               `json:"status"`
		Message string               `json:"message"`
		Data    WithdrawResponseData `json:"data"`
	}

	WithdrawResponseData struct {
		ID                      string `json:"id"`
		ExternalID              string `json:"external_id"`
		UserID                  string `json:"user_id"`
		BankCode                string `json:"bank_code"`
		AccountHolderName       string `json:"account_holder_name"`
		Amount                  int    `json:"amount"`
		DisbursementDescription string `json:"disbursement_description"`
		Status                  string `json:"status"`
	}
)
