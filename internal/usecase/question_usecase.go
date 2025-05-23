package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuestionUsecase interface {
	Create(ctx context.Context, request *model.QuestionRequest) (*model.QuestionResponse, error)
	FindByCourseCode(ctx context.Context, userId uint, courseCode string) (*[]model.QuestionResponse, error)
}

type QuestionUsecaseImpl struct {
	QuestionRepo         repository.QuestionRepository
	LecturerTeachingRepo repository.LecturerTeachingRepository
	CourseRepo           repository.CourseRepository
	DB                   *gorm.DB
	Validate             *validator.Validate
}

func NewQuestionUsecase(CourseRepo repository.CourseRepository, questionRepo repository.QuestionRepository, lecturerTeachingRepo repository.LecturerTeachingRepository, DB *gorm.DB, validate *validator.Validate) QuestionUsecase {
	return &QuestionUsecaseImpl{
		QuestionRepo:         questionRepo,
		LecturerTeachingRepo: lecturerTeachingRepo,
		CourseRepo:           CourseRepo,
		DB:                   DB,
		Validate:             validate,
	}
}

// FindByCourseCode implements QuestionUsecase.
func (questionUsecase *QuestionUsecaseImpl) FindByCourseCode(ctx context.Context, userId uint, courseCode string) (*[]model.QuestionResponse, error) {
	tx := questionUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var questions = &[]entity.Question{}
	var errorResponse model.ErrorResponse

	err := questionUsecase.CourseRepo.FindByCourseCode(tx, &entity.Course{CourseCode: courseCode})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Course data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by course from question usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = questionUsecase.LecturerTeachingRepo.FindByCourseCode(tx, userId, courseCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Data not found"
			errorResponse.Details = []string{
				fmt.Sprintf("The lecturer is not assigned to teach the course code: %s.", courseCode),
			}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by course from question usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = questionUsecase.QuestionRepo.FindByCourseCode(tx, courseCode, questions)
	if err != nil {
		log.Println("failed when find all repo class : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase class")

	return converter.QuestionToResponses(questions), nil
}

// Create implements QuestionUsecase.
func (questionUsecase *QuestionUsecaseImpl) Create(ctx context.Context, request *model.QuestionRequest) (*model.QuestionResponse, error) {
	tx := questionUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var questionDetails []entity.QuestionDetail
	var questionDetail entity.QuestionDetail
	var answerEntity entity.Answer
	var errorResponse model.ErrorResponse

	err := questionUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create question : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	lecturerTeaching := &entity.LecturerTeaching{
		CourseCode: request.CourseCode,
		UserId:     request.UserId,
	}

	err = questionUsecase.LecturerTeachingRepo.FindLecturerTeaching(tx, lecturerTeaching)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Data not found"
			errorResponse.Details = []string{
				fmt.Sprintf("The lecturer is not assigned to teach the course: %s.", request.CourseCode),
			}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error FindLecturerTeaching from question usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = questionUsecase.QuestionRepo.FindByName(tx, request.Name)
	if err == nil {
		errorResponse.Message = "Duplicate entry"
		errorResponse.Details = []string{
			"Name already exists in the database.",
		}

		jsonString, _ := json.Marshal(errorResponse)
		return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
	}

	if len(request.QuestionDetail) != int(request.QuestionCount) {
		errorResponse.Message = "Bad Request"
		errorResponse.Details = []string{
			fmt.Sprintf("You must submit exactly %d questions.", request.QuestionCount),
		}

		jsonString, _ := json.Marshal(errorResponse)
		return nil, fiber.NewError(fiber.StatusBadRequest, string(jsonString))
	}

	for _, value := range request.QuestionDetail {
		questionDetail.CorrectAnswer = value.CorrectAnswer
		questionDetail.QuestionText = value.QuestionText

		for _, answer := range value.Answer {
			answerEntity.Choice = answer.Choice
			answerEntity.Answer = answer.Answer

			questionDetail.Answers = append(questionDetail.Answers, answerEntity)
		}

		questionDetails = append(questionDetails, questionDetail)

		// reset data
		questionDetail.Answers = nil
	}

	question := &entity.Question{
		Name:            request.Name,
		QuestionCount:   request.QuestionCount,
		Duration:        request.Duration,
		CourseCode:      request.CourseCode,
		UserId:          request.UserId,
		QuestionDetails: questionDetails,
	}

	err = questionUsecase.QuestionRepo.Create(tx, question)
	if err != nil {
		log.Println("error create question :", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create question from usecase question")
	return converter.QuestionToResponse(question), nil
}
