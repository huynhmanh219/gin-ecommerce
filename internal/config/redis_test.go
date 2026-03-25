package config

import (
	"context"
	"testing"
	"time"
)

func TestRedisConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test Init Redis
	client, err := InitRedis(ctx, "localhost", "6379")
	if err != nil {
		t.Fatalf("Failed to connect Redis: %v", err)
	}
	defer client.Close()

	// Test Ping
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	t.Logf("Ping response: %s", pong)
}

func TestRedisSetGet(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := InitRedis(ctx, "localhost", "6379")
	if err != nil {
		t.Fatalf("Failed to connect Redis: %v", err)
	}
	defer client.Close()

	// Test Set
	key := "test:key"
	value := "test_value"
	err = client.Set(ctx, key, value, 1*time.Minute).Err()
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	t.Logf("✓ Set %s = %s", key, value)

	// Test Get
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if result != value {
		t.Fatalf("Expected %s, got %s", value, result)
	}
	t.Logf("✓ Get %s = %s", key, result)

	// Test Delete
	err = client.Del(ctx, key).Err()
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	t.Logf("✓ Deleted %s", key)
}

func TestRedisIncrement(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := InitRedis(ctx, "localhost", "6379")
	if err != nil {
		t.Fatalf("Failed to connect Redis: %v", err)
	}
	defer client.Close()

	// Test Increment (cho cache metrics)
	key := "cache:hits"
	err = client.Del(ctx, key).Err() // Clear trước
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Increment 3 lần
	for i := 1; i <= 3; i++ {
		val, err := client.Incr(ctx, key).Result()
		if err != nil {
			t.Fatalf("Incr failed: %v", err)
		}
		t.Logf("✓ Increment %s: %d", key, val)
	}

	// Verify final value
	final, _ := client.Get(ctx, key).Int64()
	if final != 3 {
		t.Fatalf("Expected 3, got %d", final)
	}
	t.Logf("✓ Final value: %d", final)
}
