package redis

import (
	"context"
	"time"
)

type Wrapper interface {
	Set(ctx context.Context, key string, expirationTime time.Duration, req interface{}) (err error)
	Get(ctx context.Context, key string) (resp interface{}, err error)
	GetTTL(ctx context.Context, key string) (resp time.Duration, err error)
	Delete(ctx context.Context, key string) (err error)
}
