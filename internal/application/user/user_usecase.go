package user

import (
	"zeneye-gateway/internal/adapter/repository/postgres"
	"zeneye-gateway/internal/adapter/service"
	"zeneye-gateway/internal/dto"
	"zeneye-gateway/pkg/logger"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreateUser(db *gorm.DB, req CreateUserRequest) error {
	logger.LogInfo("UserUsecase", "CreateUser", "Validating create user request", req)

	if err := validate.Struct(req); err != nil {
		logger.LogError("UserUsecase", "CreateUser", req, err)
		return err
	}

	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userController := NewUserController(userService)

	err := userController.CreateUser(req)
	if err != nil {
		logger.LogError("UserUsecase", "CreateUser", req, err)
		return err
	}

	logger.LogInfo("UserUsecase", "CreateUser", "User created successfully", req)
	return nil
}

func EditUser(db *gorm.DB, id string, req EditUserRequest) error {
	logger.LogInfo("UserUsecase", "EditUser", "Validating edit user request", req)

	if err := validate.Struct(req); err != nil {
		logger.LogError("UserUsecase", "EditUser", req, err)
		return err
	}

	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userController := NewUserController(userService)

	err := userController.EditUser(id, req)
	if err != nil {
		logger.LogError("UserUsecase", "EditUser", req, err)
		return err
	}

	logger.LogInfo("UserUsecase", "EditUser", "User edited successfully", req)
	return nil
}

func DeleteUser(db *gorm.DB, id string) error {
	logger.LogInfo("UserUsecase", "DeleteUser", "Deleting user", id)

	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userController := NewUserController(userService)

	err := userController.DeleteUser(id)
	if err != nil {
		logger.LogError("UserUsecase", "DeleteUser", id, err)
		return err
	}

	logger.LogInfo("UserUsecase", "DeleteUser", "User deleted successfully", id)
	return nil
}

func GetUser(db *gorm.DB, id string) (*dto.UserResponse, error) {
	logger.LogInfo("UserUsecase", "GetUser", "Getting user", id)

	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userController := NewUserController(userService)

	user, err := userController.GetUser(id)
	if err != nil {
		logger.LogError("UserUsecase", "GetUser", id, err)
		return nil, err
	}

	logger.LogInfo("UserUsecase", "GetUser", "User retrieved successfully", user)
	return user, nil
}

func ListUsers(db *gorm.DB) ([]*dto.UserListResponse, error) {
	logger.LogInfo("UserUsecase", "ListUsers", "Listing all users")

	repo := postgres.NewUserRepository(db)
	userService := service.NewUserService(repo)
	userController := NewUserController(userService)

	users, err := userController.ListUsers()
	if err != nil {
		logger.LogError("UserUsecase", "ListUsers", nil, err)
		return nil, err
	}

	var userList []*dto.UserListResponse
	for _, user := range users {
		userList = append(userList, &dto.UserListResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			UserUUID:  user.UserUUID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	logger.LogInfo("UserUsecase", "ListUsers", "All users listed successfully", userList)
	return userList, nil
}
