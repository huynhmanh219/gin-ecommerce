package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) CacheClient {
	return &RedisCache{client:client}
}

func (r *RedisCache) Get(ctx context.Context, key string)(string, error){
	val, err := r.client.Get(ctx,key).Result()
	if errors.Is(err,redis.Nil){
		return "",nil
	}
	
	if err != nil {
		return "",nil
	}

	return val,nil
}

func (r *RedisCache) Set(ctx context.Context,key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx,key,value,ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context,keys ... string) error{
	if len(keys) == 0 {
		return nil
	}
	return r.client.Del(ctx,keys...).Err()
}

func (r *RedisCache) Exists(ctx context.Context,keys ... string)(int64,error){
	return r.client.Exists(ctx,keys...).Result()
}

func (r *RedisCache) Increment(ctx context.Context, key string) int64 {
	return r.client.Incr(ctx,key).Val()
}