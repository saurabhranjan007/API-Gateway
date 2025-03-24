package error

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string, details string) {
	c.JSON(statusCode, ErrorResponse{
		Message: message,
		Details: details,
	})
	c.Abort()
}
