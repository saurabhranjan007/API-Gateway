package handlers

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AdminManagementHandler(c *gin.Context) {
	logger.LogInfo("AdminManagementHandler", "Handler Start", "Starting AdminManagementHandler", "")

	log.Println("AdminManagementHandler Called ....")

	// Target Service base URL from the environment variables
	baseURL := utils.GetEnv("ADMIN_MANAGEMENT_SERVICE_URL")

	// Parse the base URL
	target, err := url.Parse(baseURL)
	if err != nil {
		logger.LogError("AdminManagementHandler", "Parsing URL", baseURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing base URL"})
		return
	}

	// Replace "/admin-management" from the request path
	trimmedPath := strings.Replace(c.Request.URL.Path, "/admin-management", "", 1)
	target.Path = strings.TrimSuffix(target.Path, "/") + trimmedPath

	// New request to the target service
	req, err := http.NewRequest(c.Request.Method, target.String(), c.Request.Body)

	if err != nil {
		logger.LogError("AdminManagementHandler", "Creating Request", target.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request to target service"})
		return
	}

	// Copy headers from the original request
	for header, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogError("AdminManagementHandler", "Requesting Target Service", target.String(), err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error contacting target service"})
		return
	}
	defer resp.Body.Close()

	// Copy the response from the target service to the client
	for header, values := range resp.Header {
		for _, value := range values {
			c.Header(header, value)
		}
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)

	logger.LogInfo("AdminManagementHandler", "Handler Success", "Request successfully redirected", "")

	log.Println("AdminManagementHandler Finished!!")
}
