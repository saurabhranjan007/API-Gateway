package handlers

import (
	"errors"
	"net/http"
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/internal/adapter/service"
	user "zeneye-gateway/internal/application/user"
	"zeneye-gateway/internal/dto"
	error "zeneye-gateway/pkg/error"
	"zeneye-gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("CreateUser", "Handler Start", "Starting CreateUser handler", "")

		var req user.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.LogError("CreateUser", "Binding JSON", "", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		if err := userController.CreateUser(req); err != nil {
			logger.LogError("CreateUser", "CreateUser Error", req, err)
			if err.Error() == "email already associated with another account" {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			return
		}

		logger.LogInfo("CreateUser", "Handler Success", "User created successfully", req)
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
	}
}

func EditUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("EditUser", "Handler Start", "Starting EditUser handler", "")

		var req user.EditUserRequest
		id := c.Param("id")
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.LogError("EditUser", "Binding JSON", "", err)
			error.NewErrorResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
			return
		}

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		if err := userController.EditUser(id, req); err != nil {
			logger.LogError("EditUser", "EditUser Error", req, err)
			if err.Error() == "user not found" {
				error.NewErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
			} else if err.Error() == "email already associated with another account" || err.Error() == "username already taken" {
				error.NewErrorResponse(c, http.StatusConflict, err.Error(), "")
			} else {
				error.NewErrorResponse(c, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		logger.LogInfo("EditUser", "Handler Success", "User updated successfully", req)
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("DeleteUser", "Handler Start", "Starting DeleteUser handler", "")

		id := c.Param("id")
		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		if err := userController.DeleteUser(id); err != nil {
			logger.LogError("DeleteUser", "DeleteUser Error", id, err)
			if err.Error() == "user not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		logger.LogInfo("DeleteUser", "Handler Success", "User deleted successfully", id)
		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("GetUser", "Handler Start", "Starting GetUser handler", "")

		id := c.Param("id")
		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		user, err := userController.GetUser(id)
		if err != nil {
			logger.LogError("GetUser", "GetUser Error", id, err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		logger.LogInfo("GetUser", "Handler Success", "User retrieved successfully", user)
		c.JSON(http.StatusOK, user)
	}
}

func ListUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("ListUsers", "Handler Start", "Starting ListUsers handler", "")

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		users, err := userController.ListUsers()
		if err != nil {
			logger.LogError("ListUsers", "ListUsers Error", "", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
			return
		}

		var userList []dto.UserListResponse
		for _, user := range users {
			userList = append(userList, dto.UserListResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Role:      user.Role,
				UserUUID:  user.UserUUID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			})
		}

		logger.LogInfo("ListUsers", "Handler Success", "Users retrieved successfully", userList)
		c.JSON(http.StatusOK, userList)
	}
}

func CheckSuperadmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("CheckSuperadmin", "Handler Start", "Starting CheckSuperadmin handler", "")

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)

		superadminExists, err := userService.IsSuperadminPresent()
		if err != nil {
			logger.LogError("CheckSuperadmin", "CheckSuperadmin Error", "", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for superadmin"})
			return
		}

		logger.LogInfo("CheckSuperadmin", "Handler Success", "Superadmin existence checked successfully", superadminExists)
		if superadminExists {
			c.JSON(http.StatusOK, gin.H{"superadmin_exists": true})
		} else {
			c.JSON(http.StatusOK, gin.H{"superadmin_exists": false})
		}
	}
}

func CreateSuperadmin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.LogInfo("CreateSuperadmin", "Handler Start", "Starting CreateSuperadmin handler", "")

		var req user.CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.LogError("CreateSuperadmin", "Binding JSON", nil, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Ensure the role is superadmin
		if req.Role != "superadmin" {
			err := errors.New("role must be superadmin")
			logger.LogError("CreateSuperadmin", "Invalid Role", req.Role, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be superadmin"})
			return
		}

		repo := postgres.NewUserRepository(db)
		userService := service.NewUserService(repo)
		userController := user.NewUserController(userService)

		// Check if superadmin already exists
		superadminExists, err := userService.IsSuperadminPresent()
		if err != nil {
			logger.LogError("CreateSuperadmin", "Check Superadmin Error", "", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for superadmin"})
			return
		}

		if superadminExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Superadmin already exists"})
			return
		}

		if err := userController.CreateUser(req); err != nil {
			logger.LogError("CreateSuperadmin", "Create Superadmin Error", req, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating superadmin"})
			return
		}

		logger.LogInfo("CreateSuperadmin", "Handler Success", "Superadmin created successfully", req)
		c.JSON(http.StatusCreated, gin.H{"message": "Superadmin created successfully"})
	}
}
