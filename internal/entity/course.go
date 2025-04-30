package entity

import (
	"time"
)

type Course struct {
	ID                        uint `gorm:"primaryKey"`
	Name                      uint `gorm:"not null"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	Questions                 []Question         `gorm:"foreignKey:course_id;references:id"`
	Quizzes                   []Quiz             `gorm:"foreignKey:course_id;references:id"`
	Classes                   []Class            `gorm:"many2many:class_subjects;foreignKey:id;joinForeignKey:course_id;references:id;joinReferences:class_id"`
	LecturerCourseAssignments []LecturerTeaching `gorm:"many2many:lecturer_course_assignments;foreignKey:id;joinForeignKey:course_id;references:id;joinReferences:lecturer_teaching_code"`
}
