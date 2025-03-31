package redis

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisHandler struct {
	client *redis.Client
	sync.RWMutex
}

func New(url string) (*RedisHandler, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	return &RedisHandler{client: client}, nil
}

func (rh *RedisHandler) Ping(ctx context.Context) (string, error) {
	result, err := rh.client.Ping(ctx).Result()
	return result, err
}

func (rh *RedisHandler) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := rh.client.Set(ctx, key, value, expiration).Err()
	return err
}

func (rh *RedisHandler) Get(ctx context.Context, key string) (string, error) {
	result, err := rh.client.Get(ctx, key).Result()
	return result, err
}

// Pipelined -
func (rh *RedisHandler) Pipelined(ctx context.Context, handler func(pipe redis.Pipeliner) error) ([]redis.Cmder, error) {
	return rh.client.Pipelined(ctx, handler)
}

func (rh *RedisHandler) Pipeline(ctx context.Context) redis.Pipeliner {
	return rh.client.Pipeline()
}
func (rh *RedisHandler) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return rh.client.HGetAll(ctx, key)
}

func (rh *RedisHandler) Close() error {
	rh.Lock()
	defer rh.Unlock()
	err := rh.client.Close()
	rh.client = nil
	return err
}
