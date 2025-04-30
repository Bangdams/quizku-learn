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
	Class     Class `gorm:"foreignKey:class_id;references:id"`
	User      *User `gorm:"foreignKey:user_id;references:id"`
}
