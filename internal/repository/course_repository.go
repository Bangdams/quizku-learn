package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type CourseRepository interface {
	Create(tx *gorm.DB, course *entity.Course) error
	Update(tx *gorm.DB, course *entity.Course) error
	Delete(tx *gorm.DB, course *entity.Course) error
	FindByCourseCode(tx *gorm.DB, course *entity.Course) error
	FindAll(tx *gorm.DB, courses *[]entity.Course) error
}

type CourseRepositoryImpl struct {
	Repository[entity.Course]
}

func NewCourseRepository() CourseRepository {
	return &CourseRepositoryImpl{}
}

// FindByCourseCode implements CourseRepository.
func (repository *CourseRepositoryImpl) FindAll(tx *gorm.DB, courses *[]entity.Course) error {
	return tx.Find(courses).Error
}

// FindByCourseCode implements CourseRepository.
func (repository *CourseRepositoryImpl) FindByCourseCode(tx *gorm.DB, course *entity.Course) error {
	return tx.First(course, "course_code=?", course.CourseCode).Error
}
