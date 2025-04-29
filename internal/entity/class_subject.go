package entity

import (
	"time"
)

type ClassSubject struct {
	ID        uint `gorm:"primaryKey"`
	ClassId   uint `gorm:"not null"`
	CourseId  uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
