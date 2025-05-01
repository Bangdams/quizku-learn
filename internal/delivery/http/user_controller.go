package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByEmail(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
}

type UserControllerImpl struct {
	UserUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
	return &UserControllerImpl{
		UserUsecase: userUsecase,
	}
}

// Create implements UserController.
func (controller *UserControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.UserRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create user : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Delete implements UserController.
func (controller UserControllerImpl) Delete(ctx *fiber.Ctx) error {
	request := new(model.DeleteUserRequest)

	id, _ := ctx.ParamsInt("id")
	request.ID = uint(id)

	if err := controller.UserUsecase.Delete(ctx.UserContext(), request); err != nil {
		log.Println("failed to delete user : ", err)
		return err
	}

	return nil
}

// FindAll implements UserController.
func (controller *UserControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.UserResponse
	var err error

	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	// serach title
	keyword := ctx.Query("keyword")
	if keyword != "" {
		responses, err = controller.UserUsecase.Search(ctx.UserContext(), keyword)
		if err != nil {
			log.Println("failed to search user : ", err)
			return err
		}
	} else {
		responses, err = controller.UserUsecase.FindAll(ctx.UserContext(), uint(userID))
		if err != nil {
			log.Println("failed to find all user : ", err)
			return err
		}
	}

	return ctx.JSON(model.WebResponses[model.UserResponse]{Data: responses})
}

// FindByEmail implements UserController.
func (controller *UserControllerImpl) FindByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	response, err := controller.UserUsecase.FindByEmail(ctx.UserContext(), email)
	if err != nil {
		log.Println("failed to find by email user : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Login implements UserController.
func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.Login(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create user : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

// Update implements UserController.
func (controller *UserControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update user : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// UpdatePassword implements UserController.
func (controller *UserControllerImpl) UpdatePassword(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserPasswordRequest)

	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	request.ID = userId

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.UserUsecase.UpdatePassword(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update password user : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
