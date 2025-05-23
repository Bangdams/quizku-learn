package http

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByEmail(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	Refresh(ctx *fiber.Ctx) error
	AdminDashboardReport(ctx *fiber.Ctx) error
	LecturerDashboardReport(ctx *fiber.Ctx) error
}

type UserControllerImpl struct {
	UserUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
	return &UserControllerImpl{
		UserUsecase: userUsecase,
	}
}

// AdminDashboardReport implements UserController.
func (controller *UserControllerImpl) AdminDashboardReport(ctx *fiber.Ctx) error {
	response, err := controller.UserUsecase.AdminDashboardReport(ctx.UserContext())
	if err != nil {
		log.Println("failed to show dashborad")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AdminDashboardReportResponse]{Data: response})
}

// LecturerDashboardReport implements UserController.
func (controller *UserControllerImpl) LecturerDashboardReport(ctx *fiber.Ctx) error {
	// diambil dari jwt user id
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	response, err := controller.UserUsecase.LecturerDashboardReport(ctx.UserContext(), uint(userId))
	if err != nil {
		log.Println("failed to show dashborad")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LecturerDashboardReportResponse]{Data: response})
}

// Refresh implements UserController.
func (controller *UserControllerImpl) Refresh(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")
	if cookie == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	response, err := controller.UserUsecase.Refresh(ctx.UserContext(), cookie)
	if err != nil {
		log.Println("failed to create refresh token")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
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
		log.Println("failed to create user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Delete implements UserController.
func (controller UserControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.UserUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete user")
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
	roleParam := ctx.Query("role")
	if keyword != "" {
		responses, err = controller.UserUsecase.Search(ctx.UserContext(), keyword)
		if err != nil {
			log.Println("failed to search user")
			return err
		}
	} else {
		if roleParam != "" {
			responses, err = controller.UserUsecase.FindByRole(ctx.UserContext(), roleParam, uint(userID))
			if err != nil {
				log.Println("failed to find all user")
				return err
			}
		} else {
			responses, err = controller.UserUsecase.FindAll(ctx.UserContext(), uint(userID))
			if err != nil {
				log.Println("failed to find all user")
				return err
			}
		}
	}

	return ctx.JSON(model.WebResponses[model.UserResponse]{Data: responses})
}

// FindByEmail implements UserController.
func (controller *UserControllerImpl) FindByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	response, err := controller.UserUsecase.FindByEmail(ctx.UserContext(), email)
	if err != nil {
		log.Println("failed to find by email user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

// Login implements UserController.
func (controller *UserControllerImpl) Login(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")

	request := new(model.LoginRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, refreshToken, err := controller.UserUsecase.Login(ctx.UserContext(), request, cookie)
	if err != nil {
		log.Println("failed to login")
		return err
	}

	// durasi refreshToken
	duration := os.Getenv("DURATION_JWT_REFRESH_TOKEN")
	lifeTime, _ := strconv.Atoi(duration)

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Path:     "/",
		MaxAge:   60 * 60 * 24 * lifeTime,
	})

	return ctx.JSON(model.WebResponse[*model.LoginResponse]{Data: response})
}

// Logout implements UserController.
func (controller *UserControllerImpl) Logout(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("refresh_token")
	if cookie == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	err := controller.UserUsecase.Logout(ctx.UserContext(), cookie)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	})

	return ctx.JSON(model.WebResponse[string]{Data: "Logout successful"})
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
		log.Println("failed to update user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
