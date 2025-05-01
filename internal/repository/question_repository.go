package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	Create(tx *gorm.DB, user *entity.Question) error
	Update(tx *gorm.DB, user *entity.Question) error
	Delete(tx *gorm.DB, user *entity.Question) error
}

type QuestionRepositoryImpl struct {
	Repository[entity.Question]
}

func NewQuestionRepository() QuestionRepository {
	return &QuestionRepositoryImpl{}
}
