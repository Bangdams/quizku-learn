package entity

import (
	"time"
)

type UserClass struct {
	ID        uint `gorm:"primaryKey"`
	UserId    uint `gorm:"not null"`
	ClassId   uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
