package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("AuthMiddleware", "Handler Start", "Starting AuthMiddleware", "")

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			headerErr := errors.New("AUTH HEADER NOT FOUND")
			logger.LogWarning("AuthMiddleware", "Missing Authorization Header", "", headerErr)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated Access: invalid user"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			headerErr := errors.New("INVALID AUTH HEADER")
			logger.LogWarning("AuthMiddleware", "Invalid Authorization Header", authHeader, headerErr)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated Access: invalid user"})
			c.Abort()
			return
		}

		_, err := jwt.ValidateToken(tokenString)
		if err != nil {
			logger.LogError("AuthMiddleware", "Token Validation Error", tokenString, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated Access: invalid user"})
			c.Abort()
			return
		}

		logger.LogInfo("AuthMiddleware", "Handler Success", "Authentication successful", "")
		c.Next()
	}
}
