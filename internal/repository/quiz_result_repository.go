package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuizResultRepository interface {
	Create(tx *gorm.DB, quizResult *entity.QuizzResult) error
	FindByUserAndQuiz(tx *gorm.DB, userId uint, quizId uint, quizResult *entity.QuizzResult) error
}

type QuizResultRepositoryImpl struct {
	Repository[entity.QuizzResult]
}

func NewQuizResultRepository() QuizResultRepository {
	return &QuizResultRepositoryImpl{}
}

// FindByUserAndQuiz implements QuizResultRepository.
func (repository *QuizResultRepositoryImpl) FindByUserAndQuiz(tx *gorm.DB, userId uint, quizId uint, quizResult *entity.QuizzResult) error {
	return tx.Where("user_id = ? AND quizz_id = ?", userId, quizId).
		First(&quizResult).Error
}
