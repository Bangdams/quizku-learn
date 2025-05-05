package entity

import (
	"time"
)

type Course struct {
	CourseCode                string `gorm:"primaryKey"`
	Name                      string `gorm:"not null"`
	CreatedAt                 time.Time
	UpdatedAt                 time.Time
	Questions                 []Question         `gorm:"foreignKey:course_code;references:course_code"`
	Quizzes                   []Quiz             `gorm:"foreignKey:course_code;references:course_code"`
	Classes                   []Class            `gorm:"many2many:class_subjects;foreignKey:course_code;joinForeignKey:course_code;references:id;joinReferences:class_id"`
	LecturerCourseAssignments []LecturerTeaching `gorm:"many2many:lecturer_course_assignments;foreignKey:course_code;joinForeignKey:course_code;references:id;joinReferences:lecturer_teaching_code"`
}
