package entity

import (
	"time"
)

type LecturerCourseAssignment struct {
	ID                   uint   `gorm:"primaryKey"`
	LecturerTeachingCode string `gorm:"not null;unique"`
	CourseId             uint   `gorm:"not null;unique"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
