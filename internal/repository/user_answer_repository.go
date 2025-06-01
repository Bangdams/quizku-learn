package repository

import (
	"github.com/Bangdams/quizku-learn/internal/entity"
	"gorm.io/gorm"
)

type UserAnswerRepository interface {
	CreateBatch(tx *gorm.DB, userAnswers *[]entity.UserAnswer) error
	Update(tx *gorm.DB, userAnswer *entity.UserAnswer) error
	Delete(tx *gorm.DB, userAnswer *entity.UserAnswer) error
	FindById(tx *gorm.DB, userAnswer *entity.UserAnswer) error
	DataCheck(tx *gorm.DB, answerId uint, questionId uint) error
	FindUnansweredQuestions(tx *gorm.DB, questions *[]entity.QuestionDetail, questionId uint, userId uint) error
	VerifyUserClassQuiz(tx *gorm.DB, userId uint, quizId uint) error
	GetUserAnswersByQuestion(tx *gorm.DB, userAnswers *[]entity.UserAnswer, questionId uint, userId uint) error
}

type UserAnswerRepositoryImpl struct {
	Repository[entity.UserAnswer]
}

func NewUserAnswerRepository() UserAnswerRepository {
	return &UserAnswerRepositoryImpl{}
}

// GetUserAnswersByQuestion implements UserAnswerRepository.
func (repository *UserAnswerRepositoryImpl) GetUserAnswersByQuestion(tx *gorm.DB, userAnswers *[]entity.UserAnswer, questionId uint, userId uint) error {
	return tx.Preload("Answer.QuestionDetail", "question_details.question_id = ?", questionId).
		Where("user_answers.user_id = ?", userId).
		Find(&userAnswers).Error
}

// VerifyUserClassQuiz implements UserAnswerRepository.
func (repository *UserAnswerRepositoryImpl) VerifyUserClassQuiz(tx *gorm.DB, userId uint, quizId uint) error {
	return tx.Joins("JOIN class_subjects ON user_classes.class_id = class_subjects.class_id").
		Joins("JOIN quizzes ON class_subjects.course_code = quizzes.course_code AND class_subjects.class_id = quizzes.class_id").
		Where("user_classes.user_id = ?", userId).
		Where("quizzes.id = ?", quizId).
		First(&entity.UserClass{}).Error
}

// FindUnansweredQuestions implements UserAnswerRepository.
func (repository *UserAnswerRepositoryImpl) FindUnansweredQuestions(tx *gorm.DB, questions *[]entity.QuestionDetail, questionId uint, userId uint) error {
	return tx.Table("question_details").
		Select("question_details.question_text").
		Where("question_details.question_id = ?", questionId).
		Where("question_details.id NOT IN (?)",
			tx.Table("question_details").
				Select("DISTINCT question_details.id").
				Joins("JOIN answers ON question_details.id = answers.question_detail_id").
				Joins("JOIN user_answers ON answers.id = user_answers.answer_id").
				Where("user_answers.user_id = ?", userId),
		).
		Find(questions).Error
}

// DataCheck implements UserAnswerRepository.
func (repository *UserAnswerRepositoryImpl) DataCheck(tx *gorm.DB, answerId uint, questionId uint) error {
	return tx.Table("answers").
		Joins("JOIN question_details ON answers.question_detail_id = question_details.id").
		Where("answers.id = ? AND question_details.question_id = ?", answerId, questionId).
		First(&entity.Answer{}).Error
}

// FindById implements UserAnswerRepository.
func (repository *UserAnswerRepositoryImpl) FindById(tx *gorm.DB, userAnswer *entity.UserAnswer) error {
	panic("unimplemented")
}
