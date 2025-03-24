package middlewares

import (
	"time"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		logger.LogInfo("LoggingMiddleware", "Request Start", "Request started",
			map[string]interface{}{
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
				"client_ip": c.ClientIP(),
			},
		)

		c.Next()

		duration := time.Since(start)

		logger.LogInfo("LoggingMiddleware", "Request End", "Request completed",
			map[string]interface{}{
				"method":    c.Request.Method,
				"path":      c.Request.URL.Path,
				"status":    c.Writer.Status(),
				"duration":  duration,
				"client_ip": c.ClientIP(),
			},
		)
	}
}
