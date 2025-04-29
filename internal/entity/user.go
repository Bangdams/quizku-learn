package entity

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Role      string `gorm:"not null"`
	Image     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
