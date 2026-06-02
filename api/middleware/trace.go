package middleware

import (
	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		// Inject into gin context
		c.Set("trace_id", traceID)
		// Inject into request context for downstream (services, ws)
		c.Request = c.Request.WithContext(logger.WithTraceID(c.Request.Context(), traceID))
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}

func GetTraceID(c *gin.Context) string {
	if v, ok := c.Get("trace_id"); ok {
		return v.(string)
	}
	return ""
}
