package middlewares

import (
	"errors"
	"net/http"

	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/rate_limiter"

	"github.com/gin-gonic/gin"
)

func RateLimitingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rate_limiter.AllowRequest(c.ClientIP()) {
			logger.LogError("RateLimitingMiddleware", "Rate Limiting",
				map[string]interface{}{
					"clientIP": c.ClientIP(),
				}, errors.New("too many requests"))

			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		logger.LogInfo("RateLimitingMiddleware", "Rate Limiting", "Request allowed",
			map[string]interface{}{
				"clientIP": c.ClientIP(),
			},
		)
		c.Next()
	}
}
