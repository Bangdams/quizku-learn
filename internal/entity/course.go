package entity

import (
	"time"
)

type Course struct {
	ID        uint `gorm:"primaryKey"`
	Name      uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
