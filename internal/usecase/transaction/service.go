package transaction

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

type Service interface {
	Generate(ctx context.Context, req GenerateRequest) (res constants.DefaultResponse, err error)
	TopUp(ctx context.Context, reference string) (res constants.DefaultResponse, err error)
	Withdraw(ctx context.Context, req WithdrawRequest) (res constants.DefaultResponse, err error)
}
