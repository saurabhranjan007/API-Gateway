package middlewares

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/loadbalancer"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

// MicroserviceRoutingMiddleware handles routing requests to the appropriate microservice.
func MicroserviceRoutingMiddleware(userRepo *postgres.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			logger.LogError("MicroserviceRoutingMiddleware", "Authorization Check", "Authorization header missing", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Extract user information from the token
		claims, err := jwt.ValidateToken(strings.TrimPrefix(tokenString, "Bearer "))
		if err != nil {
			logger.LogError("MicroserviceRoutingMiddleware", "JWT Extraction", "Error extracting user info from JWT", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Additional user details from db
		user, err := userRepo.GetUser(claims.UserID)
		if err != nil {
			logger.LogError("MicroserviceRoutingMiddleware", "Fetch User Details", "Error fetching user details from DB", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Extracted user information in the request header
		c.Request.Header.Set("X-Username", user.Username)
		c.Request.Header.Set("X-User-Role", user.Role)
		c.Request.Header.Set("X-User-UUID", user.UserUUID)

		// Determine the microservice URL for the incoming request.
		microservice := loadbalancer.RouteRequest(c.Request)
		if microservice == "" {
			logger.LogError("MicroserviceRoutingMiddleware", "Route Request", "Unable to route request", errors.New("unable to route request"))
			c.JSON(http.StatusBadGateway, gin.H{"error": "Unable to route request"})
			c.Abort()
			return
		}

		logger.LogInfo("MicroserviceRoutingMiddleware", "Routing", "Routing request to microservice", map[string]interface{}{
			"path":         c.Request.URL.Path,
			"microservice": microservice,
		})

		// Parse the microservice URL.
		target, err := url.Parse(microservice)
		if err != nil {
			logger.LogError("MicroserviceRoutingMiddleware", "Parse URL", "Error parsing microservice URL", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Request path for the target microservice.
		originalPath := c.Request.URL.Path
		prefixes := []string{
			"/admin-management",
			"/agent",
			"/compliance",
			"/config",
			"/notify",
			"/bot-detection",
			"/waf",
			"/breach",
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(originalPath, prefix) {
				target.Path = strings.TrimPrefix(originalPath, prefix)
				break
			}
		}

		if strings.HasPrefix(target.Path, "/") {
			target.Path = "/" + strings.TrimPrefix(target.Path, "/")
		}

		// New request to the target service
		newReq, err := http.NewRequest(c.Request.Method, target.String(), c.Request.Body)
		if err != nil {
			logger.LogError("MicroserviceRoutingMiddleware", "Creating Request", "Error creating new request", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Copy headers from the original request and add custom headers
		for header, values := range c.Request.Header {
			for _, value := range values {
				newReq.Header.Add(header, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			logger.LogError("MicroserviceRoutingMiddleware", "Request to Target", "Error contacting target service", err)
			c.JSON(http.StatusBadGateway, gin.H{"error": "Error contacting target service"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Copy the response back to the client
		for header, values := range resp.Header {
			for _, value := range values {
				c.Header(header, value)
			}
			c.Status(resp.StatusCode)
			io.Copy(c.Writer, resp.Body)
			c.Abort()
		}
	}
}
