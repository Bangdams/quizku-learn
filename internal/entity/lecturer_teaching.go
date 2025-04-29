package entity

import (
	"time"
)

type LecturerTeaching struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"not null;unique"`
	UserId    uint   `gorm:"not null;unique"`
	ClassId   uint   `gorm:"not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
