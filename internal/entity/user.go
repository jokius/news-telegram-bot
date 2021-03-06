package entity

import (
	"time"
)

type User struct {
	ID         uint64    `gorm:"primaryKey"`
	TelegramID uint64    `gorm:"not null;index"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}
