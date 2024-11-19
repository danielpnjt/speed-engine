package payment

import (
	"context"
	"time"

	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

type paymentWrapper struct {
}

func NewPaymentWrapper() *paymentWrapper {
	return &paymentWrapper{}
}

func (w *paymentWrapper) CreateVA(ctx context.Context, req CreateVARequest) (resp CreateVAResponse, err error) {
	// xendit.Opt.SecretKey = config.GetString("xendit.secretKey")

	req.ExpectedAmount = 100000
	isSingleUse := true
	isClosed := true
	expiredDate := time.Now().Add(60 * time.Minute).UTC()
	adminFee := 4500

	// data := virtualaccount.CreateFixedVAParams{
	// 	ExternalID:     req.ExternalID,
	// 	BankCode:       req.BankCode,
	// 	Name:           req.Name,
	// 	IsSingleUse:    &isSingleUse,
	// 	IsClosed:       &isSingleUse,
	// 	ExpectedAmount: float64(int(req.ExpectedAmount) + adminFee),
	// 	ExpirationDate: &expiredDate,
	// }

	// va, err := virtualaccount.CreateFixedVA(&data)
	// if err != nil {
	// 	slog.ErrorContext(ctx, "failed to create VA", err)
	// 	err = fmt.Errorf("failed to create VA")
	// 	return
	// }

	mockResponseData := CreateVAResponseData{
		OwnerID:         "57b4e5181473eeb61c11f9b9",
		ExternalID:      req.ExternalID,
		BankCode:        "BNI",
		MerchantCode:    "8808",
		Name:            "Michael Chen",
		AccountNumber:   "8808999939380502",
		IsClosed:        &isClosed,
		ID:              "57f6fbf26b9f064272622aa6",
		IsSingleUse:     &isSingleUse,
		Status:          "PENDING",
		Currency:        "IDR",
		ExpirationDate:  &expiredDate,
		SuggestedAmount: float64(int(req.ExpectedAmount) + adminFee),
		ExpectedAmount:  float64(int(req.ExpectedAmount) + adminFee),
		Description:     "mock payment",
	}

	response := CreateVAResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    mockResponseData,
	}

	return response, nil
}

func (w *paymentWrapper) TopUp(ctx context.Context, req TopUpRequest) (resp TopUpResponse, err error) {
	// xendit.Opt.SecretKey = config.GetString("xendit.secretKey")

	// res, err := http.NewRequest("GET", "https://api.xendit.co/callback_virtual_accounts"+fmt.Sprintf(req.ExternalID), nil)
	// if err != nil {
	// 	fmt.Println("Error creating request:", err)
	// 	return
	// }

	// res.Header.Set("Content-Type", "application/json")
	// res.SetBasicAuth(xendit.Opt.SecretKey, "")

	// client := &http.Client{}
	// data, err := client.Do(res)
	// if err != nil {
	// 	fmt.Println("Error sending request:", err)
	// 	return
	// }
	// defer data.Body.Close()

	expectedAmount := 100000
	isSingleUse := true
	isClosed := true
	expiredDate := time.Now().Add(60 * time.Minute).UTC()
	adminFee := 4500

	mockResponseData := TopUpResponseData{
		OwnerID:         "57b4e5181473eeb61c11f9b9",
		ExternalID:      req.ExternalID,
		BankCode:        "BNI",
		MerchantCode:    "8808",
		Name:            "Michael Chen",
		AccountNumber:   "8808999939380502",
		IsClosed:        &isClosed,
		ID:              "57f6fbf26b9f064272622aa6",
		IsSingleUse:     &isSingleUse,
		Status:          "COMPLETED",
		Currency:        "IDR",
		ExpirationDate:  &expiredDate,
		SuggestedAmount: float64(int(expectedAmount) + adminFee),
		ExpectedAmount:  float64(int(expectedAmount) + adminFee),
		Description:     "mock payment",
	}

	response := TopUpResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    mockResponseData,
	}

	return response, nil
}

func (w *paymentWrapper) Withdraw(ctx context.Context, req WithdrawRequest) (resp WithdrawResponse, err error) {
	// xendit.Opt.SecretKey = config.GetString("xendit.secretKey")
	// bytesReq, err := json.Marshal(req)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// res, err := http.NewRequest("POST", "https://api.xendit.co/disbursements", bytes.NewReader(bytesReq))
	// if err != nil {
	// 	fmt.Println("Error creating request:", err)
	// 	return
	// }

	// res.Header.Set("Content-Type", "application/json")
	// res.SetBasicAuth(xendit.Opt.SecretKey, "")

	// client := &http.Client{}
	// data, err := client.Do(res)
	// if err != nil {
	// 	fmt.Println("Error sending request:", err)
	// 	return
	// }
	// defer data.Body.Close()

	mockResponseData := WithdrawResponseData{
		ID:                      "57f1ce05bb1a631a65eee662",
		ExternalID:              req.ExternalID,
		UserID:                  "5785e6334d7b410667d355c4",
		BankCode:                req.BankCode,
		AccountHolderName:       req.AccountHolderName,
		Amount:                  req.Amount,
		DisbursementDescription: req.Description,
		Status:                  "COMPLETED",
	}

	response := WithdrawResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    mockResponseData,
	}

	return response, nil
}
