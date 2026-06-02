package middleware

import (
	"time"

	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		ctx := c.Request.Context()
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		fields := []zap.Field{
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.Duration("latency", latency),
		}

		if statusCode >= 500 {
			logger.Error(ctx, "request", fields...)
		} else if statusCode >= 400 {
			logger.Warn(ctx, "request", fields...)
		} else {
			logger.Info(ctx, "request", fields...)
		}
	}
}
