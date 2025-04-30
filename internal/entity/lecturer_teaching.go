package entity

import (
	"time"
)

type LecturerTeaching struct {
	ID                        uint   `gorm:"primaryKey"`
	Code                      string `gorm:"not null;unique"`
	UserId                    uint   `gorm:"not null;unique"`
	ClassId                   uint   `gorm:"not null;unique"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	Class                     Class    `gorm:"foreignKey:class_id;references:id"`
	LecturerCourseAssignments []Course `gorm:"many2many:lecturer_course_assignments;foreignKey:id;joinForeignKey:lecturer_teaching_code;references:id;joinReferences:course_id"`
	User                      User     `gorm:"foreignKey:user_id;references:id"`
}
