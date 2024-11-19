package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisWrapper struct {
	redis *redis.Client
}

func NewRedisConnection(redis *redis.Client) *redisWrapper {
	return &redisWrapper{redis}
}

func (r *redisWrapper) Set(ctx context.Context, key string, expirationTime time.Duration, req interface{}) (err error) {
	client := r.redis
	json, err := json.Marshal(req)
	if err != nil {
		err = fmt.Errorf("redis marshal error: %w", err)
		return
	}
	err = client.Set(ctx, key, json, expirationTime).Err()
	if err != nil {
		err = fmt.Errorf("redis set error: %w", err)
		return
	}
	return
}

func (r *redisWrapper) Get(ctx context.Context, key string) (resp interface{}, err error) {
	client := r.redis
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return
	}

	if json.Unmarshal([]byte(val), &resp) == nil {
		return
	}

	return val, nil
}

func (r *redisWrapper) Delete(ctx context.Context, key string) (err error) {
	client := r.redis
	_, err = client.Del(ctx, key).Result()
	if err != nil {
		return
	}
	return
}
func (r *redisWrapper) GetTTL(ctx context.Context, key string) (resp time.Duration, err error) {
	client := r.redis
	resp, err = client.TTL(ctx, key).Result()
	if err != nil {
		return
	}
	return
}
