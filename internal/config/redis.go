package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, host, port string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", host, port),
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis connection failed: %v", err)
		// return nil, err
	}

	log.Printf("Redis connected: %s", pong)
	return client, nil
}

func CloseRedis(client *redis.Client) error {
	return client.Close()
}
