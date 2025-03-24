package user

import (
	"strconv"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/internal/domain/port"
	"zeneye-gateway/internal/dto"
	"zeneye-gateway/pkg/logger"
	validation "zeneye-gateway/pkg/validation"
)

type UserController struct {
	userService port.UserService
}

func NewUserController(userService port.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) CreateUser(req CreateUserRequest) error {
	logger.LogInfo("UserController", "CreateUser", "Validating user creation request", req)

	// Validate username
	if err := validation.ValidateUsername(req.Username); err != nil {
		logger.LogError("UserController", "CreateUser", req, err)
		return err
	}

	// Validate password
	if err := validation.ValidatePassword(req.Password); err != nil {
		logger.LogError("UserController", "CreateUser", req, err)
		return err
	}

	// Validate role
	if err := validation.ValidateRole(req.Role); err != nil {
		logger.LogError("UserController", "CreateUser", req, err)
		return err
	}

	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		logger.LogError("UserController", "CreateUser", req, err)
		return err
	}

	user := &entity.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

	err := c.userService.CreateUser(user)
	if err != nil {
		logger.LogError("UserController", "CreateUser", user, err)
		return err
	}

	logger.LogInfo("UserController", "CreateUser", "User created successfully", user)
	return nil
}

func (c *UserController) EditUser(id string, req EditUserRequest) error {
	logger.LogInfo("UserController", "EditUser", "Validating user edit request", req)

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.LogError("UserController", "EditUser", id, err)
		return err
	}

	// Validate username
	if err := validation.ValidateUsername(req.Username); err != nil {
		logger.LogError("UserController", "EditUser", req.Username, err)
		return err
	}

	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		logger.LogError("UserController", "EditUser", req.Email, err)
		return err
	}

	user := &entity.User{
		ID:       uint(userID),
		Username: req.Username,
		Email:    req.Email,
	}

	err = c.userService.EditUser(user)
	if err != nil {
		logger.LogError("UserController", "EditUser", user, err)
		return err
	}

	logger.LogInfo("UserController", "EditUser", "User edited successfully", user)
	return nil
}

func (c *UserController) DeleteUser(id string) error {
	logger.LogInfo("UserController", "DeleteUser", "Deleting user", id)

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.LogError("UserController", "DeleteUser", id, err)
		return err
	}

	err = c.userService.DeleteUser(uint(userID))
	if err != nil {
		logger.LogError("UserController", "DeleteUser", id, err)
		return err
	}

	logger.LogInfo("UserController", "DeleteUser", "User deleted successfully", id)
	return nil
}

func (c *UserController) GetUser(id string) (*dto.UserResponse, error) {
	logger.LogInfo("UserController", "GetUser", "Getting user", id)

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		logger.LogError("UserController", "GetUser", id, err)
		return nil, err
	}

	user, err := c.userService.GetUser(uint(userID))
	if err != nil {
		logger.LogError("UserController", "GetUser", id, err)
		return nil, err
	}

	logger.LogInfo("UserController", "GetUser", "User retrieved successfully", user)
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		UserUUID:  user.UserUUID,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (c *UserController) ListUsers() ([]*entity.User, error) {
	logger.LogInfo("UserController", "ListUsers", "Listing all users", "")

	users, err := c.userService.ListUsers()
	if err != nil {
		logger.LogError("UserController", "ListUsers", "", err)
		return nil, err
	}

	logger.LogInfo("UserController", "ListUsers", "All users listed successfully", users)
	return users, nil
}

func (c *UserController) Login(req LoginRequest) (*entity.User, error) {
	logger.LogInfo("UserController", "Login", "Authenticating user", req)

	user, err := c.userService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		logger.LogError("UserController", "Login", req, err)
		return nil, err
	}

	logger.LogInfo("UserController", "Login", "User authenticated successfully", user)
	return user, nil
}
