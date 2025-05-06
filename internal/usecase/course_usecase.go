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

type CourseUsecase interface {
	Create(ctx context.Context, request *model.CourseRequest) (*model.CourseResponse, error)
	Update(ctx context.Context, request *model.CourseRequestUpdate) (*model.CourseResponse, error)
	Delete(ctx context.Context, courseCode string) error
	FindAll(ctx context.Context) (*[]model.CourseResponse, error)
	FindByCourseCode(ctx context.Context, courseCode string) (*model.CourseResponse, error)
}

type CourseUsecaseImpl struct {
	CourseRepo repository.CourseRepository
	DB         *gorm.DB
	Validate   *validator.Validate
}

func NewCourseUsecase(courseRepo repository.CourseRepository, DB *gorm.DB, validate *validator.Validate) CourseUsecase {
	return &CourseUsecaseImpl{
		CourseRepo: courseRepo,
		DB:         DB,
		Validate:   validate,
	}
}

// FindByCourseCode implements CourseUsecase.
func (courseUsecase *CourseUsecaseImpl) FindByCourseCode(ctx context.Context, courseCode string) (*model.CourseResponse, error) {
	tx := courseUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	course := new(entity.Course)
	course.CourseCode = courseCode

	if err := courseUsecase.CourseRepo.FindByCourseCode(tx, course); err != nil {
		errorResponse := model.ErrorResponse{
			Message: "Course data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("Data not found : ")

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by course code from usecase course")

	return converter.CourseToResponse(course), nil
}

// FindAll implements CourseUsecase.
func (courseUsecase *CourseUsecaseImpl) FindAll(ctx context.Context) (*[]model.CourseResponse, error) {
	tx := courseUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var courses = &[]entity.Course{}
	err := courseUsecase.CourseRepo.FindAll(tx, courses)
	if err != nil {
		log.Println("failed when find all repo course : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase course")

	return converter.CourseToResponses(courses), nil
}

// Create implements CourseUsecase.
func (courseUsecase *CourseUsecaseImpl) Create(ctx context.Context, request *model.CourseRequest) (*model.CourseResponse, error) {
	tx := courseUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := courseUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create course : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	course := &entity.Course{
		CourseCode: request.CourseCode,
		Name:       request.Name,
	}

	err = courseUsecase.CourseRepo.Create(tx, course)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when create repo course : ", err)

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

		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase course")

	return converter.CourseToResponse(course), nil
}

// Delete implements CourseUsecase.
func (courseUsecase *CourseUsecaseImpl) Delete(ctx context.Context, courseCode string) error {
	tx := courseUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	course := &entity.Course{}
	course.CourseCode = courseCode

	_, err := courseUsecase.FindByCourseCode(ctx, courseCode)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "Course data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error delete course : ", err)

		return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	err = courseUsecase.CourseRepo.Delete(tx, course)
	if err != nil {
		log.Println("failed when delete repo course : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase course")

	return nil
}

// Update implements CourseUsecase.
func (courseUsecase *CourseUsecaseImpl) Update(ctx context.Context, request *model.CourseRequestUpdate) (*model.CourseResponse, error) {
	tx := courseUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := courseUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("Invalid request Body : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	course := &entity.Course{
		CourseCode: request.OldCourseCode,
	}

	err = courseUsecase.CourseRepo.FindByCourseCode(tx, course)
	if err != nil {
		errorResponse.Message = "Course data was not found"
		errorResponse.Details = []string{}

		jsonString, _ := json.Marshal(errorResponse)
		log.Println("Data not found")

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	if course.CourseCode != request.NewCourseCode {
		course.CourseCode = request.NewCourseCode
	}

	course.Name = request.Name

	err = courseUsecase.CourseRepo.Update(tx, course)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when update repo course : ", err.Error())

		if mysqlErr.Number == 1062 {
			return nil, fiber.ErrConflict
		}

		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success update from usecase course")

	return converter.CourseToResponse(course), nil
}
