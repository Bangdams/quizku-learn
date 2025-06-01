package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type AnswerRepository interface {
	Create(tx *gorm.DB, answer *entity.Answer) error
	Update(tx *gorm.DB, answer *entity.Answer) error
	Delete(tx *gorm.DB, answer *entity.Answer) error
	GetAnswersByQuestionDetail(tx *gorm.DB, questionDetailsId uint, choice string, answer *entity.Answer) error
}

type AnswerRepositoryImpl struct {
	Repository[entity.Answer]
}

func NewAnswerRepository() AnswerRepository {
	return &AnswerRepositoryImpl{}
}

// GetAnswersByQuestionDetailAndChoice implements AnswerRepository.
func (repository *AnswerRepositoryImpl) GetAnswersByQuestionDetail(tx *gorm.DB, questionDetailsId uint, choice string, answer *entity.Answer) error {
	return tx.Where("question_detail_id = ? and choice = ?", questionDetailsId, choice).
		First(&answer).Error
}
