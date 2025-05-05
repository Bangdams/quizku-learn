package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type CourseController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByCourseCode(ctx *fiber.Ctx) error
}

type CourseControllerImpl struct {
	CourseUsecase usecase.CourseUsecase
}

func NewCourseController(courseUsecase usecase.CourseUsecase) CourseController {
	return &CourseControllerImpl{
		CourseUsecase: courseUsecase,
	}
}

// FindByCourseCode implements CourseController.
func (controller *CourseControllerImpl) FindByCourseCode(ctx *fiber.Ctx) error {
	courseCode := ctx.Params("course_code")

	response, err := controller.CourseUsecase.FindByCourseCode(ctx.UserContext(), courseCode)
	if err != nil {
		log.Println("failed to find by course code : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: response})
}

// FindAll implements CourseController.
func (controller *CourseControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.CourseResponse

	responses, err := controller.CourseUsecase.FindAll(ctx.UserContext())
	if err != nil {
		log.Println("failed to find all course : ", err)
		return err
	}

	return ctx.JSON(model.WebResponses[model.CourseResponse]{Data: responses})
}

// Create implements CourseController.
func (controller *CourseControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.CourseRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.CourseUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create course : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: response})
}

// Delete implements CourseController.
func (controller *CourseControllerImpl) Delete(ctx *fiber.Ctx) error {
	if err := controller.CourseUsecase.Delete(ctx.UserContext(), ctx.Params("course_code")); err != nil {
		log.Println("failed to delete course : ", err)
		return err
	}

	return nil
}

// Update implements CourseController.
func (controller *CourseControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.CourseRequestUpdate)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.CourseUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update course : ", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CourseResponse]{Data: response})
}
