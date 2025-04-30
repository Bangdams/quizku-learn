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
	Answer     Answer `gorm:"foreignKey:answer_id;references:id"`
	User       User   `gorm:"foreignKey:user_id;references:id"`
	Quiz       Quiz   `gorm:"foreignKey:quizz_id;references:id"`
}
