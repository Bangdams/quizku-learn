package converter

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func QuestionWithCourseUserToResponse(question *entity.Question) *model.QuestionWithCourseUserResponse {
	log.Println("log from question to response")

	return &model.QuestionWithCourseUserResponse{
		ID:           question.ID,
		QuestionName: question.Name,
		CourseName:   question.Course.Name,
		LecturerName: question.User.Name,
		CreatedAt:    question.CreatedAt,
	}
}

func QuestionWithCourseUserToResponses(questions *[]entity.Question) *[]model.QuestionWithCourseUserResponse {
	var questionResponses []model.QuestionWithCourseUserResponse

	log.Println("log from question to responses")

	for _, question := range *questions {
		questionResponses = append(questionResponses, *QuestionWithCourseUserToResponse(&question))
	}

	return &questionResponses
}

func QuestionToResponse(question *entity.Question) *model.QuestionResponse {
	log.Println("log from question to response")

	return &model.QuestionResponse{
		ID:            question.ID,
		Name:          question.Name,
		QuestionCount: question.QuestionCount,
		Duration:      question.Duration,
		CourseCode:    question.CourseCode,
		UserId:        question.UserId,
		CreatedAt:     question.CreatedAt,
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
