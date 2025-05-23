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
	FindByCourseCode(tx *gorm.DB, courseCode string, questions *[]entity.Question) error
}

type QuestionRepositoryImpl struct {
	Repository[entity.Question]
}

func NewQuestionRepository() QuestionRepository {
	return &QuestionRepositoryImpl{}
}

// FindByCourseCode implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindByCourseCode(tx *gorm.DB, courseCode string, questions *[]entity.Question) error {
	return tx.Debug().Find(questions, "course_code = ?", courseCode).Error
}

// FindByName implements QuestionRepository.
func (repository *QuestionRepositoryImpl) FindByName(tx *gorm.DB, name string) error {
	return tx.First(&entity.Question{}, "name=?", name).Error
}
