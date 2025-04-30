package entity

import (
	"time"
)

type QuizzResult struct {
	ID                   uint   `gorm:"primaryKey"`
	UserId               uint   `gorm:"not null"`
	QuizzId              uint   `gorm:"not null"`
	Score                uint   `gorm:"not null"`
	Status               string `gorm:"not null"`
	CorrectAnswerCount   uint   `gorm:"not null"`
	IncorrectAnswerCount uint   `gorm:"not null"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Quiz                 Quiz `gorm:"foreignKey:quizz_id;references:id"`
	User                 User `gorm:"foreignKey:user_id;references:id"`
}
