package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(tx *gorm.DB, question *entity.Question) error
	Update(tx *gorm.DB, question *entity.Question) error
	Delete(tx *gorm.DB, question *entity.Question) error
	FindByName(tx *gorm.DB, name string) error
	FindById(tx *gorm.DB, question *entity.Question) error
	FindByCourseAndId(tx *gorm.DB, question *entity.Question) error
	FindByCourseCode(tx *gorm.DB, courseCode string, questions *[]entity.Question) error
	FindAll(tx *gorm.DB, questions *[]entity.Question) error
}

type QuestionRepositoryImpl struct {
	Repository[entity.Question]
}

func NewQuestionRepository() QuestionRepository {
	return &QuestionRepositoryImpl{}
}

// FindAll implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindAll(tx *gorm.DB, questions *[]entity.Question) error {
	return tx.Preload("User").Preload("Course").Find(questions).Error
}

// FindById implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindById(tx *gorm.DB, question *entity.Question) error {
	return tx.First(question).Error
}

// FindByCourseAndId implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindByCourseAndId(tx *gorm.DB, question *entity.Question) error {
	return tx.Where("course_code = ? AND id = ?", question.CourseCode, question.ID).
		First(&question).Error
}

// FindByCourseCode implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindByCourseCode(tx *gorm.DB, courseCode string, questions *[]entity.Question) error {
	return tx.Find(questions, "course_code = ?", courseCode).Error
}

// FindByName implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindByName(tx *gorm.DB, name string) error {
	return tx.First(&entity.Question{}, "name=?", name).Error
}
