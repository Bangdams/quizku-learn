package entity

import (
	"time"
)

type Class struct {
	ID        uint `gorm:"primaryKey"`
	Name      uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
