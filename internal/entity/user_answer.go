package entity

import (
	"time"
)

type UserAnswer struct {
	ID         uint `gorm:"primaryKey"`
	AnswerId   uint `gorm:"not null"`
	QuestionId uint `gorm:"not null"`
	UserId     uint `gorm:"not null"`
	QuizzId    uint `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
