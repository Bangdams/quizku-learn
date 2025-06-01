package converter

import (
	"fmt"
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
)

func QuizStudentToResponse(quiz *entity.Quiz) *model.QuizStudentResponse {
	log.Println("log from Quiz Student To Response")

	return &model.QuizStudentResponse{
		ID:            quiz.ID,
		LecturerName:  quiz.Question.User.Name,
		CourseName:    quiz.Course.Name,
		QuestionName:  quiz.Question.Name,
		Deadline:      quiz.Deadline,
		QuestionCount: quiz.Question.QuestionCount,
		Duration:      quiz.Question.Duration,
	}
}

func QuizStudentResultToResponse(userAnswers *[]entity.UserAnswer, quizResult *entity.QuizzResult) *model.QuizStudentResultResponse {
	log.Println("log from Quiz Student Result To Response")

	response := model.QuizStudentResultResponse{}
	answerDetail := model.AnswerUserDetail{}

	response.Score = quizResult.Score
	response.QuestionCount = quizResult.CorrectAnswerCount + quizResult.IncorrectAnswerCount
	response.CountCorrectAnswer = quizResult.CorrectAnswerCount
	response.CountIncorrectAnswer = quizResult.IncorrectAnswerCount

	for _, userAnswer := range *userAnswers {
		if userAnswer.Answer.Choice != userAnswer.Answer.QuestionDetail.CorrectAnswer {
			answerDetail.QuestionText = userAnswer.Answer.QuestionDetail.QuestionText
			answerDetail.IncorrectAnswer = fmt.Sprintf("%s. %s", userAnswer.Answer.Choice, userAnswer.Answer.Answer)

			for _, element := range userAnswer.Answer.QuestionDetail.Answers {
				answerDetail.CorrectAnswer = fmt.Sprintf("%s. %s", element.Choice, element.Answer)
			}

			response.AnswerUserDetails = append(response.AnswerUserDetails, answerDetail)
		}
	}

	return &response
}

func QuizStudentToResponses(quizzes *[]entity.Quiz) *[]model.QuizStudentResponse {
	var quizResponses []model.QuizStudentResponse

	log.Println("log from Quiz Student To Responses")

	for _, quiz := range *quizzes {
		quizResponses = append(quizResponses, *QuizStudentToResponse(&quiz))
	}

	return &quizResponses
}

func QuizDashboardResponse(quiz *entity.Quiz) *model.QuizDashboardResponse {
	log.Println("log from quiz to response")

	return &model.QuizDashboardResponse{
		ID:            quiz.ID,
		QuestionName:  quiz.Question.Name,
		CourseName:    quiz.Course.Name,
		QuestionCount: quiz.Question.QuestionCount,
		StudentCount:  len(quiz.Class.UserClasses),
		CreatedAt:     quiz.CreatedAt,
	}
}

func QuizDashboardResponses(quizzes *[]entity.Quiz) *[]model.QuizDashboardResponse {
	var quizResponses []model.QuizDashboardResponse

	log.Println("log from quiz to responses")

	for _, quiz := range *quizzes {
		quizResponses = append(quizResponses, *QuizDashboardResponse(&quiz))
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
