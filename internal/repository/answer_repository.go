package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	Create(tx *gorm.DB, user *entity.Answer) error
	Update(tx *gorm.DB, user *entity.Answer) error
	Delete(tx *gorm.DB, user *entity.Answer) error
}

type AnswerRepositoryImpl struct {
	Repository[entity.Answer]
}

func NewAnswerRepository() AnswerRepository {
	return &AnswerRepositoryImpl{}
}
