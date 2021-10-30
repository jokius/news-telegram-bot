package entity

import (
	"time"
)

type Group struct {
	ID           uint64    `gorm:"primaryKey"`
	UserID       uint64    `gorm:"not null;index"`
	SourceName   string    `gorm:"not null"`
	GroupName    string    `gorm:"not null"`
	LastUpdateAt time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}
