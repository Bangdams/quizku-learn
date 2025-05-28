package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuizUsecase interface {
	Create(ctx context.Context, request *model.QuizRequest, userId uint) (*model.QuizResponse, error)
	Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error)
	Delete(ctx context.Context, quizId uint) error
	QuizDashboard(ctx context.Context, userId uint, query string) (*[]model.QuizDashboardResponse, error)
}

type QuizUsecaseImpl struct {
	QuizRepo             repository.QuizRepository
	ClassRepo            repository.ClassRepository
	QuestionRepo         repository.QuestionRepository
	LecturerTeachingRepo repository.LecturerTeachingRepository
	DB                   *gorm.DB
	Validate             *validator.Validate
}

func NewQuizUsecase(quizRepo repository.QuizRepository, classRepo repository.ClassRepository, questionRepo repository.QuestionRepository, LecturerTeachingRepo repository.LecturerTeachingRepository, DB *gorm.DB, validate *validator.Validate) QuizUsecase {
	return &QuizUsecaseImpl{
		QuizRepo:             quizRepo,
		ClassRepo:            classRepo,
		QuestionRepo:         questionRepo,
		LecturerTeachingRepo: LecturerTeachingRepo,
		DB:                   DB,
		Validate:             validate,
	}
}

// QuizDashboard implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) QuizDashboard(ctx context.Context, userId uint, query string) (*[]model.QuizDashboardResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var quizzes = &[]entity.Quiz{}

	if query == "active" {
		err := quizUsecase.QuizRepo.QuizDashboardActive(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	} else if query == "archived" {
		err := quizUsecase.QuizRepo.QuizDashboardArchived(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		err := quizUsecase.QuizRepo.QuizDashboard(tx, quizzes, userId)
		if err != nil {
			log.Println("failed when find all repo quiz : ", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.QuizDashboardResponses(quizzes), nil
}

// Create implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Create(ctx context.Context, request *model.QuizRequest, userId uint) (*model.QuizResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := quizUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadRequest
	}

	lecturerTeaching := &entity.LecturerTeaching{
		CourseCode: request.CourseCode,
		UserId:     userId,
	}

	var errorResponse model.ErrorResponse

	err = quizUsecase.LecturerTeachingRepo.FindLecturerTeaching(tx, lecturerTeaching)
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

	class := &entity.Class{
		ID: request.ClassId,
	}

	err = quizUsecase.ClassRepo.FindById(tx, class)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "Class data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by class from quiz usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if lecturerTeaching.ClassId != request.ClassId {
		errorResponse.Message = "Data not found"
		errorResponse.Details = []string{
			fmt.Sprintf("The lecturer does not teach in this class : %s", class.Name),
		}

		jsonString, _ := json.Marshal(errorResponse)
		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	question := &entity.Question{
		ID:         request.QuestionId,
		CourseCode: request.CourseCode,
	}

	err = quizUsecase.QuestionRepo.FindByCourseAndId(tx, question)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse.Message = "question data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)
			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by question from quiz usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	// Format yang sesuai dengan input string
	layout := "2006-01-02"

	// Parsing string menjadi time.Time
	parsedDate, err := time.Parse(layout, request.Deadline)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return nil, fiber.NewError(fiber.ErrBadRequest.Code, "bad request")
	}

	quiz := &entity.Quiz{
		CourseCode: request.CourseCode,
		ClassId:    request.ClassId,
		QuestionId: request.QuestionId,
		Deadline:   parsedDate,
		Status:     "aktif",
	}

	err = quizUsecase.QuizRepo.Create(tx, quiz)
	if err != nil {
		log.Println("failed when create repo quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase quiz")

	return converter.QuizToResponse(quiz), nil
}

// Delete implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Delete(ctx context.Context, quizId uint) error {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	quiz := &entity.Quiz{}
	quiz.ID = quizId

	err := quizUsecase.QuizRepo.FindById(tx, quiz)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "Quiz data was not found",
				Details: []string{},
			}
			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error delete quiz : ", err)

			return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		log.Println("error find by id from quiz usecase : ", err)
		return fiber.ErrInternalServerError
	}

	err = quizUsecase.QuizRepo.Delete(tx, quiz)
	if err != nil {
		log.Println("failed when delete repo quiz : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase quiz")

	return nil
}

// Update implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error) {
	panic("unimplemented")
}
