package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserAnswerController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type UserAnswerControllerImpl struct {
	UserAnswerUsecase usecase.UserAnswerUsecase
}

func NewUserAnswerController(classUsecase usecase.UserAnswerUsecase) UserAnswerController {
	return &UserAnswerControllerImpl{
		UserAnswerUsecase: classUsecase,
	}
}

// Create implements UserAnswerController.
func (controller *UserAnswerControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.UserAnswerRequests)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	quizId, err := ctx.ParamsInt("quiz_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	response, err := controller.UserAnswerUsecase.Create(ctx.UserContext(), request, uint(userID), uint(quizId))
	if err != nil {
		log.Println("failed to create class")
		return err
	}

	log.Println(response)

	return nil
}

// Update implements UserAnswerController.
func (controller *UserAnswerControllerImpl) Update(ctx *fiber.Ctx) error {
	panic("unimplemented")
}
