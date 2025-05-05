package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func QuizDashboardResponse(quiz *entity.Quiz, studentCount uint) *model.QuizDashboardResponse {
	log.Println("log from quiz to response")

	return &model.QuizDashboardResponse{
		ID:            quiz.ID,
		CourseName:    quiz.Course.Name,
		QuestionCount: quiz.Question.QuestionCount,
		StudentCount:  studentCount,
		CreatedAt:     quiz.CreatedAt,
	}
}

func QuizDashboardResponses(quizzes *[]entity.Quiz, studentCount uint) *[]model.QuizDashboardResponse {
	var quizResponses []model.QuizDashboardResponse

	log.Println("log from quiz to responses")

	for _, quiz := range *quizzes {
		quizResponses = append(quizResponses, *QuizDashboardResponse(&quiz, studentCount))
	}

	return &quizResponses
}

func QuizToResponse(quiz *entity.Quiz) *model.QuizResponse {
	log.Println("log from quiz to response")

	return &model.QuizResponse{
		ID:         quiz.ID,
		CourseCode: quiz.CourseCode,
		ClassId:    quiz.ClassId,
		QuestionId: quiz.QuestionId,
		Deadline:   quiz.Deadline,
		Status:     quiz.Status,
		CreatedAt:  quiz.CreatedAt,
	}
}

func QuizToResponses(quizzes *[]entity.Quiz) *[]model.QuizResponse {
	var quizResponses []model.QuizResponse

	log.Println("log from quiz to responses")

	for _, quiz := range *quizzes {
		quizResponses = append(quizResponses, *QuizToResponse(&quiz))
	}

	return &quizResponses
}
