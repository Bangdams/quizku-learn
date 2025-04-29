package entity

import (
	"time"
)

type Quiz struct {
	ID         uint      `gorm:"primaryKey"`
	CourseId   uint      `gorm:"not null"`
	ClassId    uint      `gorm:"not null"`
	QuestionId uint      `gorm:"not null"`
	Deadline   time.Time `gorm:"not null"`
	Status     string    `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
