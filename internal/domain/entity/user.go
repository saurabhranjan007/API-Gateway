package entity

import (
	"time"

	"zeneye-gateway/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	UserUUID  string    `gorm:"size:36;not null;uniqueIndex"`
	Username  string    `gorm:"size:32;not null"`
	Password  string    `gorm:"size:128;not null"`
	Email     string    `gorm:"size:128;not null"`
	Role      string    `gorm:"size:32;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	logger.LogInfo("UserEntity", "BeforeCreate", "Generating UUID for new user", "")

	u.UserUUID = uuid.New().String()
	logger.LogInfo("UserEntity", "BeforeCreate", "Generated UUID for new user", u.UserUUID)
	return
}
