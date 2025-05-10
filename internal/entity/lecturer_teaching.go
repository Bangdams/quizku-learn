package entity

import (
	"time"
)

type LecturerTeaching struct {
	ID         uint   `gorm:"primaryKey"`
	CourseCode string `gorm:"not null;unique"`
	UserId     uint   `gorm:"not null;unique"`
	ClassId    uint   `gorm:"not null;unique"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	User       User   `gorm:"foreignKey:user_id;references:id"`
	Course     Course `gorm:"foreignKey:course_code;references:course_code"`
	Class      Class  `gorm:"foreignKey:class_id;references:id"`
}
