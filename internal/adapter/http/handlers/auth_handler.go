package handlers

import (
	"net/http"
	"os"
	"strconv"
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/internal/adapter/service"
	"zeneye-gateway/internal/application/user"
	error "zeneye-gateway/pkg/error"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/logger"
	"zeneye-gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("auth_handler", "Login", "Login handler called", "")

		var req user.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.LogError("auth_handler", "Login", "Error binding JSON", err)
			error.NewErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		logger.LogInfo("auth_handler", "Login", "Login request received", req)

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		user, err := userController.Login(req)
		if err != nil {
			logger.LogError("auth_handler", "Login", "Invalid username or password", err)
			error.NewErrorResponse(c, http.StatusUnauthorized, "Invalid username or password", err.Error())
			return
		}

		token, err := jwt.GenerateToken(user.ID, user.Username, user.Role, user.UserUUID)
		if err != nil {
			logger.LogError("auth_handler", "Login", "Could not generate token", err)
			error.NewErrorResponse(c, http.StatusInternalServerError, "Could not generate token", err.Error())
			return
		}

		refreshToken, err := userService.GenerateRefreshToken(user.ID)
		if err != nil {
			logger.LogError("auth_handler", "Login", "Could not generate refresh token", err)
			error.NewErrorResponse(c, http.StatusInternalServerError, "Could not generate refresh token", err.Error())
			return
		}

		jwtExpiration, _ := strconv.Atoi(utils.GetEnv("JWT_EXPIRATION"))
		refreshTokenExpiration, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION"))

		c.Header("Authorization", "Bearer "+token)
		c.Header("X-Refresh-Token", refreshToken)
		c.Header("X-Token-Expires-In", strconv.Itoa(jwtExpiration*3600))                  // convert hours to seconds
		c.Header("X-Refresh-Token-Expires-In", strconv.Itoa(refreshTokenExpiration*3600)) // convert hours to seconds

		logger.LogInfo("auth_handler", "Login", "Login successful", user)

		c.JSON(http.StatusOK, gin.H{
			"message": "Login successful",
		})
	}
}

func RefreshToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("auth_handler", "RefreshToken", "Refresh Token Called", "")

		var req struct {
			RefreshToken string `json:"refresh_token" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.LogError("auth_handler", "RefreshToken", "Error binding JSON", err)
			error.NewErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)

		token, err := userService.RefreshAccessToken(req.RefreshToken)
		if err != nil {
			logger.LogError("auth_handler", "RefreshToken", "Invalid refresh token", err)
			error.NewErrorResponse(c, http.StatusUnauthorized, "Invalid refresh token", err.Error())
			return
		}

		// Set the tokens and expiration times in the headers
		jwtExpiration, _ := strconv.Atoi(utils.GetEnv("JWT_EXPIRATION"))

		c.Header("Authorization", "Bearer "+token)
		c.Header("X-Token-Expires-In", strconv.Itoa(jwtExpiration*3600)) // convert hours to seconds

		logger.LogInfo("auth_handler", "RefreshToken", "Token refreshed successfully", "")

		c.JSON(http.StatusOK, gin.H{
			"message": "Token refreshed successfully",
		})
	}
}
