package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClassUsecase interface {
	Create(ctx context.Context, request *model.ClassRequest) (*model.ClassResponse, error)
	Update(ctx context.Context, request *model.ClassUpdate) (*model.ClassResponse, error)
	Delete(ctx context.Context, classId uint) error
	FindAll(ctx context.Context) (*[]model.ClassResponse, error)
	FindByName(ctx context.Context, className string) (*model.ClassResponse, error)
}

type ClassUsecaseImpl struct {
	ClassRepo repository.ClassRepository
	DB        *gorm.DB
	Validate  *validator.Validate
}

func NewClassUsecase(classRepo repository.ClassRepository, DB *gorm.DB, validate *validator.Validate) ClassUsecase {
	return &ClassUsecaseImpl{
		ClassRepo: classRepo,
		DB:        DB,
		Validate:  validate,
	}
}

// Create implements ClassUsecase.
func (classUsecase *ClassUsecaseImpl) Create(ctx context.Context, request *model.ClassRequest) (*model.ClassResponse, error) {
	tx := classUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := classUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create class : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	class := &entity.Class{
		Name: request.Name,
	}

	err = classUsecase.ClassRepo.Create(tx, class)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when create repo class : ", err.Error())

		var errorField string
		parts := strings.Split(mysqlErr.Message, "'")
		if len(parts) > 2 {
			errorField = parts[1]
		}

		if mysqlErr.Number == 1062 {
			errorResponse.Message = "Duplicate entry"
			errorResponse.Details = []string{errorField + " already exists in the database."}

			jsonString, _ := json.Marshal(errorResponse)

			return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase class")

	return converter.ClassToResponse(class), nil
}

// Delete implements ClassUsecase.
func (classUsecase *ClassUsecaseImpl) Delete(ctx context.Context, classId uint) error {
	tx := classUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	class := &entity.Class{}
	class.ID = classId

	err := classUsecase.ClassRepo.FindById(tx, class)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "Class data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error delete class : ", err)

		return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	err = classUsecase.ClassRepo.Delete(tx, class)
	if err != nil {
		log.Println("failed when delete repo class : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase class")

	return nil
}

// FindAll implements ClassUsecase.
func (classUsecase *ClassUsecaseImpl) FindAll(ctx context.Context) (*[]model.ClassResponse, error) {
	tx := classUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var classes = &[]entity.Class{}
	err := classUsecase.ClassRepo.FindAll(tx, classes)
	if err != nil {
		log.Println("failed when find all repo class : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase class")

	return converter.ClassToResponses(classes), nil
}

// FindByName implements ClassUsecase.
func (classUsecase *ClassUsecaseImpl) FindByName(ctx context.Context, className string) (*model.ClassResponse, error) {
	tx := classUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	class := new(entity.Class)
	class.Name = className

	if err := classUsecase.ClassRepo.FindByName(tx, class); err != nil {
		errorResponse := model.ErrorResponse{
			Message: "Class data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error find by name class : ", err)

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by class name from usecase class")

	return converter.ClassToResponse(class), nil
}

// Update implements ClassUsecase.
func (classUsecase *ClassUsecaseImpl) Update(ctx context.Context, request *model.ClassUpdate) (*model.ClassResponse, error) {
	tx := classUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := classUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update class : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	class := &entity.Class{
		ID: request.ID,
	}

	err = classUsecase.ClassRepo.FindById(tx, class)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "Class data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update class : ", err)

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	class.Name = request.Name

	err = classUsecase.ClassRepo.Update(tx, class)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when update repo class : ", err.Error())

		if mysqlErr.Number == 1062 {
			return nil, fiber.ErrConflict
		}

		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase class")

	return converter.ClassToResponse(class), nil
}
