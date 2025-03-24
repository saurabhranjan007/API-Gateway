package postgres

import (
	"zeneye-gateway/internal/domain/entity"
	"zeneye-gateway/internal/domain/port"
	"zeneye-gateway/pkg/logger"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		logger.LogError("UserRepository", "CreateUser", user, err)
	} else {
		logger.LogInfo("UserRepository", "CreateUser", "User created successfully", user)
	}
	return err
}

func (r *UserRepository) EditUser(user *entity.User) error {
	err := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}).Error
	if err != nil {
		logger.LogError("UserRepository", "EditUser", user, err)
	} else {
		logger.LogInfo("UserRepository", "EditUser", "User updated successfully", user)
	}
	return err
}

func (r *UserRepository) DeleteUser(id uint) error {
	err := r.db.Delete(&entity.User{}, id).Error
	if err != nil {
		logger.LogError("UserRepository", "DeleteUser", id, err)
	} else {
		logger.LogInfo("UserRepository", "DeleteUser", "User deleted successfully", id)
	}
	return err
}

func (r *UserRepository) GetUser(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		logger.LogError("UserRepository", "GetUser", id, err)
		return nil, err
	}
	logger.LogInfo("UserRepository", "GetUser", "User retrieved successfully", user)
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.LogError("UserRepository", "GetUserByUsername", username, err)
		return nil, err
	}
	logger.LogInfo("UserRepository", "GetUserByUsername", "User retrieved successfully", user)
	return &user, nil
}

func (r *UserRepository) DeleteRefreshToken(token string) error {
	err := r.db.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
	if err != nil {
		logger.LogError("UserRepository", "DeleteRefreshToken", token, err)
	} else {
		logger.LogInfo("UserRepository", "DeleteRefreshToken", "Refresh token deleted successfully", token)
	}
	return err
}

func (r *UserRepository) IsSuperadminPresent() (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("role = ?", "superadmin").Count(&count).Error
	if err != nil {
		logger.LogError("UserRepository", "IsSuperadminPresent", "", err)
		return false, err
	}
	logger.LogInfo("UserRepository", "IsSuperadminPresent", "Superadmin presence checked", count)
	return count > 0, nil
}

func (r *UserRepository) IsEmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		logger.LogError("UserRepository", "IsEmailExists", email, err)
		return false, err
	}
	logger.LogInfo("UserRepository", "IsEmailExists", "Email existence checked", count)
	return count > 0, nil
}

func (r *UserRepository) CreateRefreshToken(token *entity.RefreshToken) error {
	err := r.db.Create(token).Error
	if err != nil {
		logger.LogError("UserRepository", "CreateRefreshToken", token, err)
	} else {
		logger.LogInfo("UserRepository", "CreateRefreshToken", "Refresh token created successfully", token)
	}
	return err
}

func (r *UserRepository) GetRefreshToken(refreshToken string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken
	err := r.db.Where("token = ?", refreshToken).First(&token).Error
	if err != nil {
		logger.LogError("UserRepository", "GetRefreshToken", refreshToken, err)
		return nil, err
	}
	logger.LogInfo("UserRepository", "GetRefreshToken", "Refresh token retrieved successfully", token)
	return &token, nil
}

func (r *UserRepository) GetAllUsers() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Find(&users).Error
	if err != nil {
		logger.LogError("UserRepository", "GetAllUsers", "Retrieving all users", err)
		return nil, err
	}
	logger.LogInfo("UserRepository", "GetAllUsers", "All users retrieved successfully", users)
	return users, nil
}
