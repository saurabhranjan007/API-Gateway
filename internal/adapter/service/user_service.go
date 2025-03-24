package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/internal/domain/port"
	"zeneye-gateway/pkg/jwt"
	"zeneye-gateway/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// var refreshTokenExpiration, _ = time.ParseDuration(utils.GetEnv("REFRESH_TOKEN_EXPIRATION"))
var refreshTokenExpiration, _ = time.ParseDuration(os.Getenv("REFRESH_TOKEN_EXPIRATION"))

type UserService struct {
	repo port.UserRepository
}

func NewUserService(repo port.UserRepository) port.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *entity.User) error {
	logger.LogInfo("UserService", "CreateUser", "Creating a new user", user)

	// Check if a superadmin already exists
	if user.Role == "superadmin" {
		superadminExists, err := s.repo.IsSuperadminPresent()
		if err != nil {
			logger.LogError("UserService", "Checking if a superadmin already exists", user, err)
			return err
		}
		if superadminExists {
			err = errors.New("superadmin already exists")
			logger.LogError("UserService", "CreateUser", user, err)
			return err
		}
	}

	// Check if email already exists
	emailExists, err := s.repo.IsEmailExists(user.Email)
	if err != nil {
		logger.LogError("UserService", "Checking if email already exists", user, err)
		return err
	}
	if emailExists {
		err = errors.New("email already associated with another account")
		logger.LogError("UserService", "CreateUser", user, err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.LogError("UserService", "CreateUser", user, err)
		return err
	}
	user.Password = string(hashedPassword)

	user.UserUUID = uuid.New().String() // Generate UserUUID before creating a user

	err = s.repo.CreateUser(user)
	if err != nil {
		logger.LogError("UserService", "CreateUser", user, err)
		return err
	}

	logger.LogInfo("UserService", "CreateUser", "User created successfully", user)
	return nil
}

func (s *UserService) EditUser(user *entity.User) error {
	logger.LogInfo("UserService", "EditUser", "Editing user", user)

	existingUser, err := s.repo.GetUser(user.ID)
	if err != nil {
		logger.LogError("UserService", "EditUser", user, err)
		return errors.New("user not found")
	}

	// If email already exists
	emailExists, err := s.repo.IsEmailExists(user.Email)
	if err != nil {
		logger.LogError("UserService", "Checking if email already exists", user, err)
		return err
	}
	if emailExists && existingUser.Email != user.Email {
		err = errors.New("email already associated with another account")
		logger.LogError("UserService", "EditUser", user, err)
		return err
	}

	// If username already exists
	userWithSameUsername, err := s.repo.GetUserByUsername(user.Username)
	if err == nil && userWithSameUsername.ID != user.ID {
		err = errors.New("username already taken")
		logger.LogError("UserService", "EditUser", user, err)
		return err
	}

	existingUser.Username = user.Username
	existingUser.Email = user.Email
	err = s.repo.EditUser(existingUser)
	if err != nil {
		logger.LogError("UserService", "EditUser", user, err)
		return err
	}

	logger.LogInfo("UserService", "EditUser", "User updated successfully", existingUser)
	return nil
}

func (s *UserService) DeleteUser(id uint) error {
	logger.LogInfo("UserService", "DeleteUser", "Deleting user", id)

	_, err := s.repo.GetUser(id)
	if err != nil {
		logger.LogError("UserService", "DeleteUser", id, err)
		return errors.New("user not found")
	}

	err = s.repo.DeleteUser(id)
	if err != nil {
		logger.LogError("UserService", "DeleteUser", id, err)
		return err
	}

	logger.LogInfo("UserService", "DeleteUser", "User deleted successfully", id)
	return nil
}

func (s *UserService) GetUser(id uint) (*entity.User, error) {
	logger.LogInfo("UserService", "GetUser", "Getting user", id)

	user, err := s.repo.GetUser(id)
	if err != nil {
		logger.LogError("UserService", "GetUser", id, err)
		return nil, err
	}

	logger.LogInfo("UserService", "GetUser", "User retrieved successfully", user)
	return user, nil
}

func (s *UserService) ListUsers() ([]*entity.User, error) {
	logger.LogInfo("UserService", "ListUsers", "Listing all users", "")

	users, err := s.repo.GetAllUsers()
	if err != nil {
		logger.LogError("UserService", "ListUsers", "", err)
		return nil, err
	}

	logger.LogInfo("UserService", "ListUsers", "All users retrieved successfully", users)
	return users, nil
}

func (s *UserService) AuthenticateUser(username, password string) (*entity.User, error) {
	logger.LogInfo("UserService", "AuthenticateUser", "Authenticating user", username)

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		logger.LogError("UserService", "AuthenticateUser", username, err)
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.LogError("UserService", "AuthenticateUser", username, err)
		return nil, errors.New("invalid username or password")
	}

	logger.LogInfo("UserService", "AuthenticateUser", "User authenticated successfully", user)
	return user, nil
}

func (s *UserService) IsSuperadminPresent() (bool, error) {
	logger.LogInfo("UserService", "IsSuperadminPresent", "Checking if superadmin is present", "")

	isPresent, err := s.repo.IsSuperadminPresent()
	if err != nil {
		logger.LogError("UserService", "IsSuperadminPresent", "", err)
		return false, err
	}

	logger.LogInfo("UserService", "IsSuperadminPresent", "Superadmin presence checked", isPresent)
	return isPresent, nil
}

func (s *UserService) GenerateRefreshToken(userID uint) (string, error) {
	logger.LogInfo("UserService", "GenerateRefreshToken", "Generating refresh token", userID)

	refreshToken := make([]byte, 32)
	if _, err := rand.Read(refreshToken); err != nil {
		logger.LogError("UserService", "GenerateRefreshToken", userID, err)
		return "", err
	}

	tokenString := hex.EncodeToString(refreshToken)
	token := &entity.RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: time.Now().Add(refreshTokenExpiration), // from .env
	}

	if err := s.repo.CreateRefreshToken(token); err != nil {
		logger.LogError("UserService", "GenerateRefreshToken", userID, err)
		return "", err
	}

	logger.LogInfo("UserService", "GenerateRefreshToken", "Refresh token generated successfully", tokenString)
	return tokenString, nil
}

func (s *UserService) RefreshAccessToken(refreshToken string) (string, error) {
	logger.LogInfo("UserService", "RefreshAccessToken", "Refreshing access token", refreshToken)

	token, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil {
		logger.LogError("UserService", "RefreshAccessToken", refreshToken, err)
		return "", errors.New("invalid refresh token")
	}

	if token.ExpiresAt.Before(time.Now()) {
		err = errors.New("refresh token has expired")
		logger.LogError("UserService", "RefreshAccessToken", refreshToken, err)
		return "", err
	}

	// Fetch user details using the UserID from the refresh token
	user, err := s.repo.GetUser(token.UserID)
	if err != nil {
		logger.LogError("UserService", "GetUser", token.UserID, err)
		return "", err
	}

	// Generate new access token with additional user details
	newToken, err := jwt.GenerateToken(user.ID, user.Username, user.Role, user.UserUUID)
	if err != nil {
		logger.LogError("UserService", "RefreshAccessToken", user.ID, err)
		return "", err
	}

	logger.LogInfo("UserService", "RefreshAccessToken", "Access token refreshed successfully", newToken)
	return newToken, nil
}
