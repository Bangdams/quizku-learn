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
	UserAnswers      []UserAnswer   `gorm:"foreignKey:answer_id;references:id"`
	QuestionDetail   QuestionDetail `gorm:"foreignKey:question_detail_id;references:id"`
}
