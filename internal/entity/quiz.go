package entity

import (
	"time"
)

type Quiz struct {
	ID           uint      `gorm:"primaryKey"`
	CourseId     uint      `gorm:"not null"`
	ClassId      uint      `gorm:"not null"`
	QuestionId   uint      `gorm:"not null"`
	Deadline     time.Time `gorm:"not null"`
	Status       string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Course       Course        `gorm:"foreignKey:course_id;references:id"`
	Question     Question      `gorm:"foreignKey:question_id;references:id"`
	Class        Class         `gorm:"foreignKey:class_id;references:id"`
	QuizzResults []QuizzResult `gorm:"foreignKey:quizz_id;references:id"`
	UserAnswers  []UserAnswer  `gorm:"foreignKey:quizz_id;references:id"`
}
