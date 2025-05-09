package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bangdams/quizku-learn/internal/entity"
	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/model/converter"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/Bangdams/quizku-learn/internal/util"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUsecase interface {
	Create(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, userId uint) error
	FindAll(ctx context.Context, userId uint) (*[]model.UserResponse, error)
	FindByRole(ctx context.Context, role string, userId uint) (*[]model.UserResponse, error)
	FindByEmail(ctx context.Context, emailRequest string) (*model.UserResponse, error)
	Search(ctx context.Context, keyword string) (*[]model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginRequest, requestRefreshToken string) (*model.LoginResponse, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
}

type UserUsecaseImpl struct {
	UserRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	ClassRepo        repository.ClassRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewUserUsecase(userRepo repository.UserRepository, RefreshTokenRepo repository.RefreshTokenRepository, classRepo repository.ClassRepository, DB *gorm.DB, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		UserRepo:         userRepo,
		RefreshTokenRepo: RefreshTokenRepo,
		ClassRepo:        classRepo,
		DB:               DB,
		Validate:         validate,
	}
}

// FindByRole implements UserUsecase.
func (userUsecase *UserUsecaseImpl) FindByRole(ctx context.Context, role string, userId uint) (*[]model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var users = &[]entity.User{}
	err := userUsecase.UserRepo.FindByRole(tx, role, userId, users)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by role from usecase user")

	return converter.UserToResponses(users), nil
}

// Logout implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Logout(ctx context.Context, refreshToken string) error {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := util.ParseToken(refreshToken, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return fiber.ErrUnauthorized
	}

	userId := claims["user_id"].(float64)

	if err := userUsecase.RefreshTokenRepo.FindById(tx, uint(userId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Logout failed because the user has not logged in.")
			return fiber.NewError(fiber.StatusBadRequest, "User has not logged in")
		} else {
			log.Println("Error RefreshToken findbyid:", err)
			return fiber.ErrInternalServerError
		}
	} else {
		log.Println("Logout successful.")
		userUsecase.RefreshTokenRepo.Update(tx, &entity.RefreshToken{UserId: uint(userId), StatusLogout: 1})
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

// Refresh implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	claims, err := util.ParseToken(refreshToken, []byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	userId := claims["user_id"].(float64)
	request := entity.User{
		ID:    uint(userId),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
		Role:  claims["role"].(string),
	}

	if err := userUsecase.RefreshTokenRepo.CheckStatusLogout(tx, uint(userId)); err != nil {
		return nil, fiber.ErrUnauthorized
	}

	newAccessToken, _ := util.GenerateAccessToken(&request)

	log.Println("success create access token")

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.LoginUserToResponse(newAccessToken), nil
}

// Create implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Create(ctx context.Context, request *model.UserRequest) (*model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error create user : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to generate password")
		return nil, fiber.ErrInternalServerError
	}

	user := &entity.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: string(password),
		Role:     request.Role,
		Image:    request.Image,
	}

	if err := userUsecase.UserRepo.FindByEmail(tx, user); err == nil {
		errorResponse.Message = "Duplicate entry"
		errorResponse.Details = []string{"email already exists in the database."}

		jsonString, _ := json.Marshal(errorResponse)

		return nil, fiber.NewError(fiber.ErrConflict.Code, string(jsonString))
	}

	if request.Role == "mahasiswa" {
		if request.ClassId == 0 {
			errorResponse.Message = "invalid request parameter"
			errorResponse.Details = []string{"Field ClassId failed on required rule"}

			jsonString, _ := json.Marshal(errorResponse)

			return nil, fiber.NewError(fiber.ErrBadGateway.Code, string(jsonString))
		}

		class := &entity.Class{
			ID: request.ClassId,
		}

		err := userUsecase.ClassRepo.FindById(tx, class)
		if err != nil {
			errorResponse.Message = "Class data was not found"
			errorResponse.Details = []string{}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by id class in create user repo : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		}

		user.UserClass.ClassId = request.ClassId
	}

	err = userUsecase.UserRepo.Create(tx, user)
	if err != nil {
		log.Println("failed when create repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success create from usecase user")

	return converter.UserToResponse(user), nil
}

// Delete implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Delete(ctx context.Context, userId uint) error {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{}
	user.ID = userId

	err := userUsecase.UserRepo.FindById(tx, user)
	if err != nil {
		errorResponse := model.ErrorResponse{
			Message: "User data was not found",
			Details: []string{},
		}
		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error delete user : ", err)

		return fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	err = userUsecase.UserRepo.Delete(tx, user)
	if err != nil {
		log.Println("failed when delete repo user : ", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return fiber.ErrInternalServerError
	}

	log.Println("success delete from usecase user")

	return nil
}

// FindAll implements UserUsecase.
func (userUsecase *UserUsecaseImpl) FindAll(ctx context.Context, userId uint) (*[]model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var users = &[]entity.User{}
	err := userUsecase.UserRepo.FindAll(tx, userId, users)
	if err != nil {
		log.Println("failed when find all repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find all from usecase user")

	return converter.UserToResponses(users), nil
}

// FindByEmail implements UserUsecase.
func (userUsecase *UserUsecaseImpl) FindByEmail(ctx context.Context, emailRequest string) (*model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	user.Email = emailRequest

	if err := userUsecase.UserRepo.FindByEmail(tx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := model.ErrorResponse{
				Message: "User data was not found",
				Details: []string{},
			}

			jsonString, _ := json.Marshal(errorResponse)

			log.Println("error find by email user usecase : ", err)

			return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
		} else {
			log.Println("Error find by email user usecase:", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success find by email from usecase user")

	return converter.UserToResponse(user), nil
}

// Login implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Login(ctx context.Context, request *model.LoginRequest, requestRefreshTokenUser string) (*model.LoginResponse, string, error) {
	now := time.Now()

	_, err := util.ParseToken(requestRefreshTokenUser, []byte(os.Getenv("SECRET_KEY")))
	if err == nil {
		return nil, "", fiber.NewError(fiber.StatusBadRequest, "Refresh token still valid")
	}

	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{}

	if err := userUsecase.UserRepo.Login(tx, user, request.Email); err != nil {
		log.Println("invalid Email : ", err)
		return nil, "", fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println("invalid password : ", err)
		return nil, "", fiber.ErrUnauthorized
	}

	accessToken, err := util.GenerateAccessToken(user)
	if err != nil {
		log.Println("Failed to generate token jwt")
		return nil, "", fiber.ErrInternalServerError
	}

	refreshToken, err := util.GenerateRefreshToken(user)
	if err != nil {
		log.Println("Failed to generate token jwt")
		return nil, "", fiber.ErrInternalServerError
	}

	duration := os.Getenv("DURATION_JWT_REFRESH_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	requestRefreshToken := &entity.RefreshToken{
		UserId:       user.ID,
		StatusLogout: 0,
		Token:        refreshToken,
		ExpiresAt:    now.Add(time.Minute * time.Duration(lifeTime)),
	}

	if err := userUsecase.RefreshTokenRepo.FindById(tx, user.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Creating new refresh token in database")
			userUsecase.RefreshTokenRepo.Create(tx, requestRefreshToken)
		} else {
			log.Println("Error fetching refresh token:", err)
			return nil, "", fiber.ErrInternalServerError
		}
	} else {
		log.Println("Updating existing refresh token in database")
		userUsecase.RefreshTokenRepo.Update(tx, requestRefreshToken)
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, "", fiber.ErrInternalServerError
	}

	log.Println("success login")

	return converter.LoginUserToResponse(accessToken), refreshToken, nil
}

// Search implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Search(ctx context.Context, keyword string) (*[]model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	var users = &[]entity.User{}
	err := userUsecase.UserRepo.Search(tx, users, keyword)
	if err != nil {
		log.Println("failed when search repo user : ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success search from usecase user")

	return converter.UserToResponses(users), nil
}

// Update implements UserUsecase.
func (userUsecase *UserUsecaseImpl) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	errorResponse := &model.ErrorResponse{}

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on '%s' rule", e.Field(), e.Tag())
			validationErrors = append(validationErrors, msg)
		}

		errorResponse.Message = "invalid request parameter"
		errorResponse.Details = validationErrors

		jsonString, _ := json.Marshal(errorResponse)

		log.Println("error update user : ", err)

		return nil, fiber.NewError(fiber.ErrBadRequest.Code, string(jsonString))
	}

	user := &entity.User{
		Email: request.Email,
	}

	err = userUsecase.UserRepo.FindByEmail(tx, user)
	if err != nil {
		errorResponse.Message = "User data was not found"
		errorResponse.Details = []string{}

		jsonString, _ := json.Marshal(errorResponse)
		log.Println("Data not found")

		return nil, fiber.NewError(fiber.ErrNotFound.Code, string(jsonString))
	}

	if request.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("failed to generate password")
			return nil, fiber.ErrInternalServerError
		}

		user.Password = string(password)
	}

	user.Name = request.Name
	user.Image = request.Image

	err = userUsecase.UserRepo.Update(tx, user)
	if err != nil {
		mysqlErr := err.(*mysql.MySQLError)
		log.Println("failed when update repo user : ", err)

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

	log.Println("success update from usecase user")

	return converter.UserToResponse(user), nil
}
