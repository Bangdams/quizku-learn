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
	QuizDashboard(ctx context.Context) (*[]model.QuizResponse, error)
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
func (quizUsecase *QuizUsecaseImpl) QuizDashboard(ctx context.Context) (*[]model.QuizResponse, error) {
	tx := quizUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var quizzes = &[]entity.Quiz{}
	err := quizUsecase.QuizRepo.QuizDashboard(tx, quizzes)
	if err != nil {
		log.Println("failed when find all repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return nil, nil
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
		log.Println("failed when create repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase user")

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
