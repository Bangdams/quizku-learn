package entity

import (
	"time"
)

type Answer struct {
	ID               uint   `gorm:"primaryKey"`
	QuestionDetailId uint   `gorm:"not null"`
	Answer           string `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
