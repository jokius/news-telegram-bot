package entity

import (
	"time"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	TelegramID uint `gorm:"index"`
	VkToken    string
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}
