package bank

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

type Service interface {
	SubmitBank(ctx context.Context, req SubmitBank) (res constants.DefaultResponse, err error)
	FindAll(ctx context.Context) (res constants.DefaultResponse, err error)
}
