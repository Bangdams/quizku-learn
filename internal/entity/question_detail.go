package entity

import (
	"time"
)

type QuestionDetail struct {
	ID         uint   `gorm:"primaryKey"`
	QuestionId uint   `gorm:"not null"`
	Question   string `gorm:"not null"`
	AnswerId   uint   `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
