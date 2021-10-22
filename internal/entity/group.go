package entity

import (
	"time"
)

type Group struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	GroupLink    string    `gorm:"not null"`
	LastUpdateAt time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	User
}
