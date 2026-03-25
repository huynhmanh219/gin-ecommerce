package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CacheMetricsHandler(redisClient *redis.Client, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Get as string first, then parse
		hitsStr, _ := redisClient.Get(ctx, "cache:hits").Result()
		missesStr, _ := redisClient.Get(ctx, "cache:misses").Result()

		hits := int64(0)
		misses := int64(0)

		// Parse safely
		if hitsStr != "" {
			if h, err := strconv.ParseInt(hitsStr, 10, 64); err == nil {
				hits = h
			} else {
				logger.Warn("failed to parse cache:hits", zap.String("value", hitsStr), zap.Error(err))
			}
		}

		if missesStr != "" {
			if m, err := strconv.ParseInt(missesStr, 10, 64); err == nil {
				misses = m
			} else {
				logger.Warn("failed to parse cache:misses", zap.String("value", missesStr), zap.Error(err))
			}
		}

		total := hits + misses
		hitRatio := 0.0

		if total > 0 {
			hitRatio = float64(hits) / float64(total) * 100
		}

		info := redisClient.Info(ctx, "memory").Val()

		logger.Info("cache_metrics",
			zap.Int64("hits", hits),
			zap.Int64("misses", misses),
			zap.Float64("hit_ratio_percent", hitRatio),
		)

		c.JSON(200, gin.H{
			"cache_hits":        hits,
			"cache_misses":      misses,
			"hit_ratio_percent": hitRatio,
			"redis_info":        info,
		})
	}
}
