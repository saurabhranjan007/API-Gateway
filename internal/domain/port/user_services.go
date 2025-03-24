package port

import "zeneye-gateway/internal/domain/entity"

type UserService interface {
	CreateUser(user *entity.User) error
	EditUser(user *entity.User) error
	DeleteUser(id uint) error
	GetUser(id uint) (*entity.User, error)
	ListUsers() ([]*entity.User, error)

	AuthenticateUser(username, password string) (*entity.User, error)
	IsSuperadminPresent() (bool, error)
	GenerateRefreshToken(userID uint) (string, error)
	RefreshAccessToken(refreshToken string) (string, error)
}
