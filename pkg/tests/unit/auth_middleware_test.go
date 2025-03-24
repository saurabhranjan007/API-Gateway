package unit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"zeneye-gateway/internal/adapter/http/middlewares"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {

	logger.InitLogger() // ensure logger
	defer logger.SyncLogger()

	router := gin.New()
	router.Use(middlewares.AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		logger.LogInfo("TestAuthMiddleware", "Handler", "Accessed /test endpoint", "")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	token, _ := jwt.GenerateToken(1)                 // Pass a uint ID
	req.Header.Set("Authorization", "Bearer "+token) // Add "Bearer" prefix

	logger.LogInfo("TestAuthMiddleware", "Test", "Sending request", map[string]interface{}{
		"Method":  req.Method,
		"URL":     req.URL,
		"Headers": req.Header,
	})

	router.ServeHTTP(w, req)

	logger.LogInfo("TestAuthMiddleware", "Test", "Received response", map[string]interface{}{
		"StatusCode": w.Code,
		"Body":       w.Body.String(),
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
