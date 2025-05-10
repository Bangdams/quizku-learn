package usecase

import (
	"context"
	"encoding/json"
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

type LecturerTeachingUsecase interface {
	Create(ctx context.Context, request *model.FlexibleLecturerTeachingRequest) (*[]model.LecturerTeachingResponse, error)
	Delete(ctx context.Context, lecturerId uint) error
	FindById(ctx context.Context, lecturerId uint) (*model.LecturerTeachingResponse, error)
	DisplayData(ctx context.Context) (*model.DisplayDataResponse, error)
	DisplayDataWithClassId(ctx context.Context, classId []uint) (*model.DisplayDataWitClassIdResponse, error)
}

type LecturerTeachingUsecaseImpl struct {
	LecturerTeachingRepo repository.LecturerTeachingRepository
	ClassRepo            repository.ClassRepository
	CourseRepo           repository.CourseRepository
	UserRepo             repository.UserRepository
	DB                   *gorm.DB
	Validate             *validator.Validate
}

func NewLecturerTeachingUsecase(lecturerTeaching repository.LecturerTeachingRepository, classRepo repository.ClassRepository, courseRepo repository.CourseRepository, userRepo repository.UserRepository, DB *gorm.DB, validate *validator.Validate) LecturerTeachingUsecase {
	return &LecturerTeachingUsecaseImpl{
		LecturerTeachingRepo: lecturerTeaching,
		ClassRepo:            classRepo,
		CourseRepo:           courseRepo,
		UserRepo:             userRepo,
		DB:                   DB,
		Validate:             validate,
	}
}

// DisplayData implements LecturerTeachingUsecase.
func (lecturerTeachingUsecase *LecturerTeachingUsecaseImpl) DisplayData(ctx context.Context) (*model.DisplayDataResponse, error) {
	tx := lecturerTeachingUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var classes []entity.Class
	var courses []entity.Course

	err := lecturerTeachingUsecase.ClassRepo.FindAll(tx, &classes)
	if err != nil {
		log.Println("error find all class in lecturer usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	err = lecturerTeachingUsecase.CourseRepo.FindAll(tx, &courses)
	if err != nil {
		log.Println("error find all course in lecturer usecase : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by id from usecase lecturerTeaching")

	return converter.DisplayDataToResponses(&classes, &courses), err
}

// DisplayDataWithClassId implements LecturerTeachingUsecase.
func (lecturerTeachingUsecase *LecturerTeachingUsecaseImpl) DisplayDataWithClassId(ctx context.Context, classId []uint) (*model.DisplayDataWitClassIdResponse, error) {
	tx := lecturerTeachingUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var courses []entity.Course

	if classId != nil {
		err := lecturerTeachingUsecase.CourseRepo.FindByIdWithClass(tx, &courses, classId)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success display data with class id from usecase lecturerTeaching")

	return converter.DisplayDataWitClassIdToResponses(&courses), nil
}

// FindById implements LecturerTeachingUsecase.
func (lecturerTeachingUsecase *LecturerTeachingUsecaseImpl) FindById(ctx context.Context, lecturerId uint) (*model.LecturerTeachingResponse, error) {
	tx := lecturerTeachingUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	lecturerTeaching := &entity.LecturerTeaching{}
	lecturerTeaching.ID = lecturerId

	err := lecturerTeachingUsecase.LecturerTeachingRepo.FindById(tx, lecturerTeaching)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "LecturerTeaching data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error find by id lecturerTeaching : ", err)

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by id from usecase lecturerTeaching")

	return converter.LecturerTeachingToResponse(lecturerTeaching), nil
}

// Create implements LecturerTeachingUsecase.
func (lecturerTeachingUsecase *LecturerTeachingUsecaseImpl) Create(ctx context.Context, request *model.FlexibleLecturerTeachingRequest) (*[]model.LecturerTeachingResponse, error) {
	tx := lecturerTeachingUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := lecturerTeachingUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create lecturerTeaching : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	err = lecturerTeachingUsecase.UserRepo.FindById(tx, &entity.User{ID: request.UserID})
	if err != nil {
		errorResponse.Message = "User data was not found"
		errorResponse.Details = []string{}

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error find by id user : ", err)

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	var lecturerTeachings []entity.LecturerTeaching

	for _, item := range request.Teachings {
		err = lecturerTeachingUsecase.ClassRepo.FindById(tx, &entity.Class{ID: item.ClassID})
		if err != nil {
			errorResponse.Message = "Class data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id class : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		for _, code := range item.CourseCodes {
			err = lecturerTeachingUsecase.CourseRepo.FindByCourseCode(tx, &entity.Course{CourseCode: code})
			if err != nil {
				errorResponse.Message = "Course data was not found"
				errorResponse.Details = []string{}

				jsonString, _ := json.Marshal(errorResponse)

				log.Println("error find by Code course : ", err)

				return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
			}

			lecturer := entity.LecturerTeaching{
				CourseCode: code,
				UserId:     request.UserID,
				ClassId:    item.ClassID,
			}

			if err := lecturerTeachingUsecase.LecturerTeachingRepo.OneDataCheck(tx, &lecturer); err == nil {
				errorResponse.Message = "Duplicate entry"
				errorResponse.Details = []string{"lecturer teaching already exists in the database."}

				jsonString, _ := json.Marshal(errorResponse)

				return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
			}

			lecturerTeachings = append(lecturerTeachings, lecturer)
		}
	}

	err = lecturerTeachingUsecase.LecturerTeachingRepo.CreateBacth(tx, &lecturerTeachings)
	if err != nil {
		log.Println("failed when create repo lecturerTeaching : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase lecturerTeaching")

	return converter.LecturerTeachingToResponses(&lecturerTeachings), nil
}

// Delete implements LecturerTeachingUsecase.
func (lecturerTeachingUsecase *LecturerTeachingUsecaseImpl) Delete(ctx context.Context, lecturerId uint) error {
	tx := lecturerTeachingUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	lecturerTeaching := &entity.LecturerTeaching{}
	lecturerTeaching.ID = lecturerId

	err := lecturerTeachingUsecase.LecturerTeachingRepo.FindById(tx, lecturerTeaching)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "LecturerTeaching data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error delete lecturerTeaching : ", err)

		return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	err = lecturerTeachingUsecase.LecturerTeachingRepo.Delete(tx, lecturerTeaching)
	if err != nil {
		log.Println("failed when delete repo lecturerTeaching : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase lecturerTeaching")

	return nil
}
