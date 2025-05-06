package entity

import (
	"time"
)

type Class struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string `gorm:"not null;unique"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Courses           []Course           `gorm:"many2many:class_subjects;foreignKey:id;joinForeignKey:class_id;references:course_code;joinReferences:course_code"`
	Quizzes           []Quiz             `gorm:"foreignKey:class_id;references:id"`
	UserClasses       []UserClass        `gorm:"foreignKey:class_id;references:id"`
	LecturerTeachings []LecturerTeaching `gorm:"foreignKey:class_id;references:id"`
}
