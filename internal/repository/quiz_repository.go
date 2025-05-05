package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(tx *gorm.DB, quiz *entity.Quiz) error
	Update(tx *gorm.DB, quiz *entity.Quiz) error
	Delete(tx *gorm.DB, quiz *entity.Quiz) error
	QuizDashboard(tx *gorm.DB, quizzes *[]entity.Quiz) error
}

type QuizRepositoryImpl struct {
	Repository[entity.Quiz]
}

func NewQuizRepository() QuizRepository {
	return &QuizRepositoryImpl{}
}

// QuizDashboard implements QuizRepository.
func (repository *QuizRepositoryImpl) QuizDashboard(tx *gorm.DB, quizzes *[]entity.Quiz) error {
	return tx.Preload("Course").Preload("Question").Find(quizzes).Error
}
