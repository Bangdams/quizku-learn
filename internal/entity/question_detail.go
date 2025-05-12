package entity

import (
	"time"
)

type QuestionDetail struct {
	ID            uint   `gorm:"primaryKey"`
	QuestionId    uint   `gorm:"not null"`
	QuestionText  string `gorm:"not null"`
	CorrectAnswer string `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Answers       []Answer `gorm:"foreignKey:question_detail_id;references:id"`
	Question      Question `gorm:"foreignKey:question_id;references:id"`
}
