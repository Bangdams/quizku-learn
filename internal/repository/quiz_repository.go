package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(tx *gorm.DB, user *entity.Quiz) error
	Update(tx *gorm.DB, user *entity.Quiz) error
	Delete(tx *gorm.DB, user *entity.Quiz) error
}

type QuizRepositoryImpl struct {
	Repository[entity.Quiz]
}

func NewQuizRepository() QuizRepository {
	return &QuizRepositoryImpl{}
}
