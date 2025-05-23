package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type QuizController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	QuizDashboard(ctx *fiber.Ctx) error
}

type QuizControllerImpl struct {
	QuizUsecase usecase.QuizUsecase
}

func NewQuizController(userUsecase usecase.QuizUsecase) QuizController {
	return &QuizControllerImpl{
		QuizUsecase: userUsecase,
	}
}

func (controller *QuizControllerImpl) QuizDashboard(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	responses, err := controller.QuizUsecase.QuizDashboard(ctx.UserContext(), uint(userID))
	if err != nil {
		log.Println("failed to show quiz dashborad")
		return err
	}

	return ctx.JSON(model.WebResponses[model.QuizDashboardResponse]{Data: responses})
}

// Create implements QuizController.
func (controller *QuizControllerImpl) Create(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// Delete implements QuizController.
func (controller *QuizControllerImpl) Delete(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// Update implements QuizController.
func (controller *QuizControllerImpl) Update(ctx *fiber.Ctx) error {
	panic("unimplemented")
}
