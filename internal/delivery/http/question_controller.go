package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type QuestionController interface {
	Create(ctx *fiber.Ctx) error
	FindByCourseCode(ctx *fiber.Ctx) error
}

type QuestionControllerImpl struct {
	QuestionUsecase usecase.QuestionUsecase
}

func NewQuestionController(questionUsecase usecase.QuestionUsecase) QuestionController {
	return &QuestionControllerImpl{
		QuestionUsecase: questionUsecase,
	}
}

// FindByCourseCode implements QuestionController.
func (controller *QuestionControllerImpl) FindByCourseCode(ctx *fiber.Ctx) error {
	// diambil dari jwt user id
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	courseCode := ctx.Params("course_code")

	responses, err := controller.QuestionUsecase.FindByCourseCode(ctx.UserContext(), uint(userId), courseCode)
	if err != nil {
		log.Println("failed to find by course code from question")
		return err
	}

	return ctx.JSON(model.WebResponses[model.QuestionResponse]{Data: responses})
}

// Create implements QuestionController.
func (controller *QuestionControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.QuestionRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	// diambil dari jwt user id
	userToken := ctx.Locals("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userId := claims["user_id"].(float64)

	request.UserId = uint(userId)

	response, err := controller.QuestionUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create question")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.QuestionResponse]{Data: response})
}
