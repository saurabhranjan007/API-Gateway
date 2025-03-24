package entity

import "time"

type Session struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"size:256;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time
}
