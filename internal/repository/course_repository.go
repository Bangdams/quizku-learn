package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(tx *gorm.DB, user *entity.Course) error
	Update(tx *gorm.DB, user *entity.Course) error
	Delete(tx *gorm.DB, user *entity.Course) error
}

type CourseRepositoryImpl struct {
	Repository[entity.Course]
}

func NewCourseRepository() CourseRepository {
	return &CourseRepositoryImpl{}
}
