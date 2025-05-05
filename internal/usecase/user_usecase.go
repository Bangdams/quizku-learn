package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
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
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
	FindAll(ctx context.Context, userId uint) (*[]model.UserResponse, error)
	FindByEmail(ctx context.Context, emailRequest string) (*model.UserResponse, error)
	Search(ctx context.Context, keyword string) (*[]model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginRequest, requestRefreshToken string) (*model.LoginResponse, string, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*model.LoginResponse, error)
}

type UserUsecaseImpl struct {
	UserRepo         repository.UserRepository
	RefreshTokenRepo repository.RefreshTokenRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewUserUsecase(userRepo repository.UserRepository, RefreshTokenRepo repository.RefreshTokenRepository, DB *gorm.DB, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		UserRepo:         userRepo,
		RefreshTokenRepo: RefreshTokenRepo,
		DB:               DB,
		Validate:         validate,
	}
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
			return fiber.ErrBadRequest
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

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadRequest
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
		log.Println("duplicate email : ", user.Email)
		return nil, fiber.ErrConflict
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
func (userUsecase *UserUsecaseImpl) Delete(ctx context.Context, request *model.DeleteUserRequest) error {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return fiber.ErrBadRequest
	}

	user := &entity.User{}
	user.ID = request.ID

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
		log.Println("Data not found : ")
		return nil, fiber.ErrNotFound
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

	err := userUsecase.Validate.Struct(request)
	if err != nil {
		log.Println("Invalid request Body : ", err)
		return nil, fiber.ErrBadGateway
	}

	user := &entity.User{
		Email: request.Email,
	}

	err = userUsecase.UserRepo.FindByEmail(tx, user)
	if err != nil {
		log.Println("Data not found")
		return nil, fiber.ErrNotFound
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
		log.Println("failed when update repo user : ", err.Error())

		if mysqlErr.Number == 1062 {
			return nil, fiber.ErrConflict
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
