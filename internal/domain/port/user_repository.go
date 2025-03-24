package port

import "zeneye-gateway/internal/domain/entity"

type UserRepository interface {
	CreateUser(user *entity.User) error
	EditUser(user *entity.User) error
	DeleteUser(id uint) error
	GetUser(id uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetAllUsers() ([]*entity.User, error)

	IsSuperadminPresent() (bool, error)
	CreateRefreshToken(token *entity.RefreshToken) error
	GetRefreshToken(token string) (*entity.RefreshToken, error)
	DeleteRefreshToken(token string) error
	IsEmailExists(email string) (bool, error)
}
