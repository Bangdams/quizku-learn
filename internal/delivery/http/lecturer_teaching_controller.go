package http

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Bangdams/quizku-learn/internal/model"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type LecturerTeachingController interface {
	Create(ctx *fiber.Ctx) error
	DisplayData(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
	FindById(ctx *fiber.Ctx) error
}

type LecturerTeachingrollerImpl struct {
	LecturerTeaching usecase.LecturerTeachingUsecase
}

func NewLecturerTeachingController(lecturerTeaching usecase.LecturerTeachingUsecase) LecturerTeachingController {
	return &LecturerTeachingrollerImpl{
		LecturerTeaching: lecturerTeaching,
	}
}

// DisplayData implements LecturerTeachingController.
func (controller *LecturerTeachingrollerImpl) DisplayData(ctx *fiber.Ctx) error {
	query := ctx.Query("class-ids")
	if query != "" {
		strIds := strings.Split(query, ",")

		var classIds []uint
		var validID = regexp.MustCompile(`^\d+$`)

		for _, strId := range strIds {
			if !validID.MatchString(strId) {
				return ctx.Status(400).JSON(fiber.Map{"error": "ID harus berupa angka"})
			}

			id, err := strconv.Atoi(strId)
			if err == nil {
				classIds = append(classIds, uint(id))
			}
		}

		responses, err := controller.LecturerTeaching.DisplayDataWithClassId(ctx.UserContext(), classIds)
		if err != nil {
			return err
		}

		return ctx.JSON(model.WebResponse[*model.DisplayDataWitClassIdResponse]{Data: responses})
	}

	responses, err := controller.LecturerTeaching.DisplayData(ctx.UserContext())
	if err != nil {
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DisplayDataResponse]{Data: responses})
}

// Create implements LecturerTeachingController.
func (controller *LecturerTeachingrollerImpl) Create(ctx *fiber.Ctx) error {
	var responses *[]model.LecturerTeachingResponse
	var err error

	request := new(model.FlexibleLecturerTeachingRequest)

	if err = ctx.BodyParser(request); err != nil {
		log.Println("failed to parse request : ", err)
		return fiber.ErrBadRequest
	}

	responses, err = controller.LecturerTeaching.Create(ctx.UserContext(), request)
	if err != nil {
		log.Println("failed to create course")
		return err
	}

	return ctx.JSON(model.WebResponses[model.LecturerTeachingResponse]{Data: responses})
}

// Delete implements LecturerTeachingController.
func (controller *LecturerTeachingrollerImpl) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("lecturer_teaching_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := controller.LecturerTeaching.Delete(ctx.UserContext(), uint(id)); err != nil {
		log.Println("failed to delete lecturerTeaching")
		return err
	}

	return nil
}

// FindById implements LecturerTeachingController.
func (controller *LecturerTeachingrollerImpl) FindById(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("lecturer_teaching_id")
	if err != nil {
		return fiber.ErrBadRequest
	}

	response, err := controller.LecturerTeaching.FindById(ctx.UserContext(), uint(id))
	if err != nil {
		log.Println("failed to find by id lecturerTeaching")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LecturerTeachingResponse]{Data: response})
}
