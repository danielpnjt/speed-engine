package payment

import "context"

type Wrapper interface {
	CreateVA(ctx context.Context, req CreateVARequest) (resp CreateVAResponse, err error)
	TopUp(ctx context.Context, req TopUpRequest) (resp TopUpResponse, err error)
	Withdraw(ctx context.Context, req WithdrawRequest) (resp WithdrawResponse, err error)
}
