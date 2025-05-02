package usecase

import (
	"context"
	"log"

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
	UpdatePassword(ctx context.Context, request *model.UpdateUserPasswordRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
	FindAll(ctx context.Context, userId uint) (*[]model.UserResponse, error)
	FindByEmail(ctx context.Context, emailRequest string) (*model.UserResponse, error)
	Search(ctx context.Context, keyword string) (*[]model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error)
}

type UserUsecaseImpl struct {
	UserRepo repository.UserRepository
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewUserUsecase(userRepo repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserUsecase {
	return &UserUsecaseImpl{
		UserRepo: userRepo,
		DB:       DB,
		Validate: validate,
	}
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
		log.Println("duplicate email : ", err)
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
func (userUsecase *UserUsecaseImpl) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	tx := userUsecase.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := &entity.User{}

	if err := userUsecase.UserRepo.Login(tx, user, request.Email); err != nil {
		log.Println("invalid Email : ", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		log.Println("invalid password : ", err)
		return nil, fiber.ErrUnauthorized
	}

	accessToken, err := util.GenerateAccessToken(user)
	// refreshToken, err := util.GenerateAccessToken(user)
	if err != nil {
		log.Println("Failed to generate token jwt")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		log.Println("Failed commit transaction : ", err)
		return nil, fiber.ErrInternalServerError
	}

	log.Println("success login")

	return converter.LoginUserToResponse(accessToken), nil
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

// UpdatePassword implements UserUsecase.
func (*UserUsecaseImpl) UpdatePassword(ctx context.Context, request *model.UpdateUserPasswordRequest) (*model.UserResponse, error) {
	panic("unimplemented")
}
