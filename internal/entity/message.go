package entity

import (
	"time"
)

type Message struct {
	ID        uint64    `gorm:"primaryKey"`
	GroupID   uint64    `gorm:"not null;index"`
	MessageID uint64    `gorm:"not null"`
	Source    string    `gorm:"not null"`
	MessageAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	User      User      `gorm:"foreignKey:GroupID"`
}
