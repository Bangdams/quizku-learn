package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func QuestionToResponse(question *entity.Question) *model.QuestionResponse {
	log.Println("log from question to response")

	return &model.QuestionResponse{
		ID:            question.ID,
		Name:          question.Name,
		QuestionCount: question.QuestionCount,
		Duration:      question.Duration,
		CourseCode:    question.Duration,
		UserId:        question.UserId,
	}
}

func QuestionToResponses(questions *[]entity.Question) *[]model.QuestionResponse {
	var questionResponses []model.QuestionResponse

	log.Println("log from question to responses")

	for _, question := range *questions {
		questionResponses = append(questionResponses, *QuestionToResponse(&question))
	}

	return &questionResponses
}
