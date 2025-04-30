package entity

import (
	"time"
)

type User struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string `gorm:"not null"`
	Email             string `gorm:"not null;unique"`
	Password          string `gorm:"not null"`
	Role              string `gorm:"not null"`
	Image             string `gorm:"not null"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	LecturerTeachings []LecturerTeaching `gorm:"foreignKey:user_id;references:id"`
	UserClass         UserClass          `gorm:"foreignKey:user_id;references:id"`
	Questions         []Question         `gorm:"foreignKey:user_id;references:id"`
	QuizzResults      []QuizzResult      `gorm:"foreignKey:user_id;references:id"`
	UserAnswers       []UserAnswer       `gorm:"foreignKey:user_id;references:id"`
}
