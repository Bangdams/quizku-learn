package usecase

import (
	"context"
	"log"

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
		log.Println("Data not found : ")
		return nil, fiber.ErrNotFound
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

	err := courseUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadRequest
	}

	course := &entity.Course{
		CourseCode: request.CourseCode,
		Name:       request.Name,
	}

	err = courseUsecase.CourseRepo.Create(tx, course)
	if err != nil {
		log.Println("failed when create repo course : ", err)
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

	err := courseUsecase.CourseRepo.Delete(tx, course)
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

	err := courseUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadGateway
	}

	course := &entity.Course{
		CourseCode: request.OldCourseCode,
	}

	err = courseUsecase.CourseRepo.FindByCourseCode(tx, course)
	if err != nil {
		log.Println("Data not found")
		return nil, fiber.ErrNotFound
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
