package user

import (
	"context"

	"github.com/danielpnjt/speed-engine/internal/pkg/constants"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (res constants.DefaultResponse, err error)
	Login(ctx context.Context, req LoginRequest) (res constants.DefaultResponse, err error)
	Logout(ctx context.Context) (res constants.DefaultResponse, err error)
	GetDetail(ctx context.Context, userID int) (res constants.DefaultResponse, err error)
	GetAll(ctx context.Context, req FindAllRequest) (res constants.DefaultResponse, err error)
	GetDetailPlayer(ctx context.Context) (res constants.DefaultResponse, err error)
}
