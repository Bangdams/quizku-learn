package usecase

import (
	"context"
	"log"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type QuizUsecase interface {
	Create(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error)
	Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error)
	Delete(ctx context.Context, request *model.QuizRequest) error
	QuizDashboard(ctx context.Context, userId uint) (*[]model.QuizDashboardResponse, error)
}

type QuizUsecaseImpl struct {
	QuizRepo repository.QuizRepository
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewQuizUsecase(quizRepo repository.QuizRepository, DB *gorm.DB, validate *validator.Validate) QuizUsecase {
	return &QuizUsecaseImpl{
		QuizRepo: quizRepo,
		DB:       DB,
		Validate: validate,
	}
}

// QuizDashboard implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) QuizDashboard(ctx context.Context, userId uint) (*[]model.QuizDashboardResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var quizzes = &[]entity.Quiz{}
	err := quizUsecase.QuizRepo.QuizDashboard(tx, quizzes, userId)
	if err != nil {
		log.Println("failed when find all repo quiz : ", err)
		return nil, fiber.ErrInternalServerError
	}

	// for _, quiz := range *quizzes {
	// 	log.Println(quiz.CourseCode)
	// 	log.Println("--------------")
	// 	log.Println(quiz.Course.Name)
	// 	log.Println("--------------")
	// 	log.Println(quiz.Question.Name)
	// 	log.Println("--------------")
	// 	log.Println("question_count : ", quiz.Question.QuestionCount)
	// 	log.Println("==============")
	// 	log.Println("user_count : ", len(quiz.Class.UserClasses))
	// 	log.Println("created at : ", quiz.CreatedAt)
	// }

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.QuizDashboardResponses(quizzes), nil
}

// Create implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Create(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := quizUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadRequest
	}

	quiz := &entity.Quiz{}

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
func (quizUsecase *QuizUsecaseImpl) Delete(ctx context.Context, request *model.QuizRequest) error {
	panic("unimplemented")
}

// Update implements QuizUsecase.
func (quizUsecase *QuizUsecaseImpl) Update(ctx context.Context, request *model.QuizRequest) (*model.QuizResponse, error) {
	panic("unimplemented")
}
