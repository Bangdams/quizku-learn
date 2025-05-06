package http

import (
	"log"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ClassController interface {
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
	FindByName(ctx *fiber.Ctx) error
}

type ClassControllerImpl struct {
	ClassUsecase usecase.ClassUsecase
}

func NewClassController(classUsecase usecase.ClassUsecase) ClassController {
	return &ClassControllerImpl{
		ClassUsecase: classUsecase,
	}
}

// Create implements ClassController.
func (controller *ClassControllerImpl) Create(ctx *fiber.Ctx) error {
	request := new(model.ClassRequest)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.ClassUsecase.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create class")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ClassResponse]{Data: response})
}

// Delete implements ClassController.
func (controller *ClassControllerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("class_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.ClassUsecase.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete class")
		return err
	}

	return nil
}

// FindAll implements ClassController.
func (controller *ClassControllerImpl) FindAll(ctx *fiber.Ctx) error {
	var responses *[]model.ClassResponse

	responses, err := controller.ClassUsecase.FindAll(ctx.UserContext())
	if err != nil {
		log.Println("failed to find all class")
		return err
	}

	return ctx.JSON(model.WebResponses[model.ClassResponse]{Data: responses})
}

// FindByName implements ClassController.
func (controller *ClassControllerImpl) FindByName(ctx *fiber.Ctx) error {
	className := ctx.Params("class_name")

	response, err := controller.ClassUsecase.FindByName(ctx.UserContext(), className)
	if err != nil {
		log.Println("failed to find by class name")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ClassResponse]{Data: response})
}

// Update implements ClassController.
func (controller *ClassControllerImpl) Update(ctx *fiber.Ctx) error {
	request := new(model.ClassUpdate)

	if err := ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	response, err := controller.ClassUsecase.Update(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to update class")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.ClassResponse]{Data: response})
}
