package entity

import (
	"time"
)

type Question struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	QuestionCount uint   `gorm:"not null"`
	Duration      uint   `gorm:"not null"`
	CourseId      uint   `gorm:"not null"`
	UserId        uint   `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
