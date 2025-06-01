package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserAnswerUsecase interface {
	Create(ctx context.Context, request *model.UserAnswerRequests, userId uint, quizId uint) (*model.QuizResponse, error)
	Update(ctx context.Context, request *model.UserAnswerRequest) (*model.QuizResponse, error)
}

type UserAnswerUsecaseImpl struct {
	QuizResultRepo repository.QuizResultRepository
	QuizRepo       repository.QuizRepository
	UserAnswerRepo repository.UserAnswerRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewUserAnswerUsecase(quizResultRepo repository.QuizResultRepository, quizRepo repository.QuizRepository, userAnswerRepo repository.UserAnswerRepository, DB *gorm.DB, validate *validator.Validate) UserAnswerUsecase {
	return &UserAnswerUsecaseImpl{
		QuizResultRepo: quizResultRepo,
		QuizRepo:       quizRepo,
		UserAnswerRepo: userAnswerRepo,
		DB:             DB,
		Validate:       validate,
	}
}

// Create implements UserAnswerUsecase.
func (userAnswerUsecase *UserAnswerUsecaseImpl) Create(ctx context.Context, request *model.UserAnswerRequests, userId uint, quizId uint) (*model.QuizResponse, error) {
	if len(request.UserAnswers) == 0 {
		return nil, fiber.NewError(fiber.ErrBadRequest.Code, "request body cannot be empty")
	}

	tx := userAnswerUsecase.DB.WithContext(ctx).Begin()

	errorResponse := &model.ErrorResponse{}
	for _, element := range request.UserAnswers {
		if err := userAnswerUsecase.Validate.Struct(element); err != nil {
			if validationErrs, ok := err.(validator.ValidationErrors); ok {
				var validationErrors []string
				for _, e := range validationErrs {
					msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
					validationErrors = append(validationErrors, msg)
				}

				errorResponse.Message = "invalid request parameter"
				errorResponse.Details = validationErrors

				jsonString, _ := json.Marshal(errorResponse)

				log.Println("error create user answer:", err)
				tx.Rollback()

				return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
			}
		}
	}

	quiz := &entity.Quiz{
		ID: quizId,
	}

	err := userAnswerUsecase.QuizRepo.FindById(tx, quiz)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "quiz data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id quiz : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = userAnswerUsecase.UserAnswerRepo.VerifyUserClassQuiz(tx, userId, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "You do not have access to this quiz because it is not assigned to your class."
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error verify user class : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error verify user class : ", err)
		return nil, fiber.ErrInternalServerError
	}

	var userAnswers []entity.UserAnswer
	for _, request := range request.UserAnswers {
		err = userAnswerUsecase.UserAnswerRepo.DataCheck(tx, request.AnswerId, quiz.QuestionId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorResponse.Message = "answer data was not found"
				errorResponse.Details = []string{}

				jsonString, _ := json.Marshal(errorResponse)

				log.Println("error data check : ", err)

				return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
			}

			log.Println("error data check : ", err)
			return nil, fiber.ErrInternalServerError
		}

		userAnswers = append(userAnswers, entity.UserAnswer{
			AnswerId: request.AnswerId,
			UserId:   userId,
			QuizzId:  quizId,
		})
	}

	err = userAnswerUsecase.UserAnswerRepo.CreateBatch(tx, &userAnswers)
	if err != nil {
		log.Println("failed when create repo user answer : ", err)
		return nil, fiber.ErrInternalServerError
	}

	var countCorrectAnswer int64
	var countIncorrectAnswer int64

	questions := &[]entity.QuestionDetail{}
	userAnswerUsecase.UserAnswerRepo.FindUnansweredQuestions(tx, questions, quiz.QuestionId, userId)

	for _, question := range *questions {
		countIncorrectAnswer += 1
		log.Println(question.QuestionText)
	}

	userAnswers = []entity.UserAnswer{}
	err = userAnswerUsecase.UserAnswerRepo.GetUserAnswersByQuestion(tx, &userAnswers, quiz.QuestionId, userId)
	if err != nil {
		log.Println("failed when GetUserAnswersByQuestion : ", err)
		return nil, fiber.ErrInternalServerError
	}

	for _, userAnswerCoba := range userAnswers {
		if userAnswerCoba.Answer.Choice == userAnswerCoba.Answer.QuestionDetail.CorrectAnswer {
			countCorrectAnswer += 1
			continue
		}
		countIncorrectAnswer += 1
	}

	var score int
	if quiz.Question.QuestionCount == 0 {
		score = 0
	} else {
		score = int((float64(countCorrectAnswer) / float64(quiz.Question.QuestionCount)) * 100)
	}

	quizResult := &entity.QuizzResult{
		UserId:               userId,
		QuizzId:              quizId,
		Score:                uint(score),
		CorrectAnswerCount:   uint(countCorrectAnswer),
		IncorrectAnswerCount: uint(countIncorrectAnswer),
	}

	// kkm
	thresholdScore, err := strconv.Atoi(os.Getenv("THRESHOLD_SCORE"))
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if score >= thresholdScore {
		quizResult.Status = "lulus"
	} else {
		quizResult.Status = "gagal"
	}

	err = userAnswerUsecase.QuizResultRepo.Create(tx, quizResult)
	if err != nil {
		log.Println("failed when Create quiz resut repo : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction:", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase user answer")
	return nil, nil
}

// Update implements UserAnswerUsecase.
func (userAnswerUsecase *UserAnswerUsecaseImpl) Update(ctx context.Context, request *model.UserAnswerRequest) (*model.QuizResponse, error) {
	panic("unimplemented")
}
