package http

import (
	"log"
	"strings"

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
	FindByUserAndCourse(ctx *fiber.Ctx) error
}

type QuizControllerImpl struct {
	QuizUsecase usecase.QuizUsecase
}

func NewQuizController(userUsecase usecase.QuizUsecase) QuizController {
	return &QuizControllerImpl{
		QuizUsecase: userUsecase,
	}
}

// FindByUserAndCourse implements QuizController.
func (controller *QuizControllerImpl) FindByUserAndCourse(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	courseCode := strings.ToUpper(ctx.Params("course_code"))

	response, err := controller.QuizUsecase.FindByUserAndCourse(ctx.UserContext(), uint(userID), courseCode)
	if err != nil {
		log.Println("failed to create quiz")
		return err
	}

	return ctx.JSON(model.WebResponses[model.QuizStudentResponse]{Data: response})
}

func (controller *QuizControllerImpl) QuizDashboard(ctx *fiber.Ctx) error {
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)

	status := ctx.Query("status")

	responses, err := controller.QuizUsecase.QuizDashboard(ctx.UserContext(), uint(userID), status)
	if err != nil {
		log.Println("failed to show quiz dashborad")
		return err
	}

	return ctx.JSON(model.WebResponses[model.QuizDashboardResponse]{Data: responses})
}

// Create implements QuizController.
func (controller *QuizControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.QuizRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	// diambil dari jwt user id
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	response, err := controller.QuizUsecase.Create(ctx.UserContext(), request, uint(userId))
	if err != nil {
		log.Println("failed to create quiz")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.QuizResponse]{Data: response})
}

// Delete implements QuizController.
func (controller *QuizControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("quiz_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.QuizUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete quiz")
		return err
	}

	return nil
}

// Update implements QuizController.
func (controller *QuizControllerImpl) Update(ctx *fiber.Ctx) error {
	panic("unimplemented")
}
