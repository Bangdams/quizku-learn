package repository

import (
	"time"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"gorm.io/gorm"
)

type QuizResultRepository interface {
	Create(tx *gorm.DB, quizResult *entity.QuizzResult) error
	FindByUserAndQuiz(tx *gorm.DB, userId uint, quizId uint, quizResult *entity.QuizzResult) error
	QuizResultAnalysis(tx *gorm.DB, quizId uint, quizResult *[]entity.QuizzResult) (model.QuizResultAnalysisResponse, error)
	GetQuizHistoryForStudent(tx *gorm.DB, userId uint, quizResults *[]entity.QuizzResult) error
}

type QuizResultRepositoryImpl struct {
	Repository[entity.QuizzResult]
}

func NewQuizResultRepository() QuizResultRepository {
	return &QuizResultRepositoryImpl{}
}

// GetQuizHistoryForStudent implements QuizResultRepository.
func (repository *QuizResultRepositoryImpl) GetQuizHistoryForStudent(tx *gorm.DB, userId uint, quizResults *[]entity.QuizzResult) error {
	return tx.Preload("Quiz.Question.Course").
		Where("quizz_results.user_id = ?", userId).
		Find(quizResults).Error
}

// QuizResultAnalysis implements QuizResultRepository.
func (repository *QuizResultRepositoryImpl) QuizResultAnalysis(tx *gorm.DB, quizId uint, quizResult *[]entity.QuizzResult) (model.QuizResultAnalysisResponse, error) {
	var userCount int64
	var scoreMinMax []model.ScoreResult
	var userCompleteCount int64
	var userIncompleteCount int64
	var completionRate int64
	var totalScore int64
	var avarageScore int64
	var courseName string
	var questionName string
	var questionCount uint
	var CreatedAt time.Time

	err := tx.Table("quizzes").
		Select("COUNT(*) AS user_count").
		Joins("JOIN user_classes ON quizzes.class_id = user_classes.class_id").
		Where("quizzes.id = ?", quizId).
		Scan(&userCount).Error
	if err != nil {
		return model.QuizResultAnalysisResponse{}, err
	}

	err = tx.Raw(`
    SELECT score, 'max_score' AS type FROM 
    (SELECT score FROM quizz_results WHERE quizz_id = ? ORDER BY score DESC LIMIT 1) AS max_score
    UNION
    SELECT score, 'min_score' AS type FROM 
    (SELECT score FROM quizz_results WHERE quizz_id = ? ORDER BY score ASC LIMIT 1) AS min_score;
	`, quizId, quizId).Scan(&scoreMinMax).Error
	if err != nil {
		return model.QuizResultAnalysisResponse{}, err
	}

	err = tx.Preload("Quiz.Question").
		Preload("Quiz.Course").
		Preload("User.UserClass.Class").
		Where("quizz_results.quizz_id = ?", quizId).
		Find(&quizResult).Error
	if err != nil {
		return model.QuizResultAnalysisResponse{}, err
	}

	students := []model.QuizResultAnalysisStudent{}
	student := model.QuizResultAnalysisStudent{}

	for _, element := range *quizResult {
		student.Name = element.User.Name
		student.Email = element.User.Email
		student.ClassName = element.User.UserClass.Class.Name
		student.Score = element.Score
		student.Status = element.Status
		students = append(students, student)

		courseName = element.Quiz.Course.Name
		questionName = element.Quiz.Question.Name
		questionCount = element.Quiz.Question.QuestionCount
		CreatedAt = element.Quiz.CreatedAt
		totalScore += int64(element.Score)
	}

	userCompleteCount = int64(len(*quizResult))
	userIncompleteCount = userCount - userCompleteCount
	completionRate = int64((float64(userCompleteCount) / float64(userCount)) * 100)
	avarageScore = totalScore / userCompleteCount

	response := model.QuizResultAnalysisResponse{
		AvarageScore:        uint(avarageScore),
		CompletionRate:      uint(completionRate),
		StudentCount:        uint(userCount),
		UserCompleteCount:   uint(userCompleteCount),
		UserIncompleteCount: uint(userIncompleteCount),
		CoursesName:         courseName,
		QuizName:            questionName,
		QuestionCount:       questionCount,
		CreatedAt:           CreatedAt,
		Students:            students,
	}

	for _, element := range scoreMinMax {
		if element.Type == "max_score" {
			response.HighestScore = uint(element.Score)
			continue
		}
		response.LowestScore = uint(element.Score)
	}

	return response, nil
}

// FindByUserAndQuiz implements QuizResultRepository.
func (repository *QuizResultRepositoryImpl) FindByUserAndQuiz(tx *gorm.DB, userId uint, quizId uint, quizResult *entity.QuizzResult) error {
	return tx.Where("user_id = ? AND quizz_id = ?", userId, quizId).
		First(&quizResult).Error
}
