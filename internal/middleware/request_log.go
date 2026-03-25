package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		start := time.Now()
		requestID := uuid.NewString()
		ctx.Set("request_id",requestID)

		ctx.Next()

		latency := time.Since(start)
		status := ctx.Writer.Status()

		fields := []zap.Field{
			zap.String("request_id",requestID),
			zap.String("path",ctx.Request.URL.Path),
			zap.String("method",ctx.Request.Method),
			zap.Int("status",status),
			zap.Duration("latency",latency),
		}
		
		if status >= 500 {
			logger.Error("request_failed",fields...)
			return
		}

		if status >= 400 {
			logger.Warn("request_warning",fields...)
			return
		}

		logger.Info("request_completed",fields...)
	}
}