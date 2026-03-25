package cache

import (
	"context"
	"time"
)

type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)

	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	Delete(ctx context.Context,keys ...string) error

	Exists(ctx context.Context, keys ...string) (int64,error)

	Increment(ctx context.Context, key string) int64
}