package config

import (
	"github.com/Bangdams/quizku-learn/internal/delivery/http"
	"github.com/Bangdams/quizku-learn/internal/delivery/http/route"
	"github.com/Bangdams/quizku-learn/internal/repository"
	"github.com/Bangdams/quizku-learn/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Validate *validator.Validate
}

func Bootstrap(config *BootstrapConfig) {
	// repo
	userRepo := repository.NewUserRepository()
	refreshTokenRepo := repository.NewRefreshTokenRepository()
	courseRepo := repository.NewCourseRepository()
	classRepo := repository.NewClassRepository()
	lecturerTeachingRepo := repository.NewLecturerTeachingRepository()

	// usecase
	userUsecase := usecase.NewUserUsecase(userRepo, refreshTokenRepo, classRepo, config.DB, config.Validate)
	courseUsecase := usecase.NewCourseUsecase(courseRepo, config.DB, config.Validate)
	classUsecase := usecase.NewClassUsecase(classRepo, courseRepo, config.DB, config.Validate)
	lecturerTeachingUsecase := usecase.NewLecturerTeachingUsecase(lecturerTeachingRepo, classRepo, courseRepo, userRepo, config.DB, config.Validate)

	// controller
	userController := http.NewUserController(userUsecase)
	courseController := http.NewCourseController(courseUsecase)
	classController := http.NewClassController(classUsecase)
	lecturerTeachingController := http.NewLecturerTeachingController(lecturerTeachingUsecase)

	routeConfig := route.RouteConfig{
		App:                        config.App,
		UserController:             userController,
		CourseController:           courseController,
		ClassController:            classController,
		LecturerTeachingController: lecturerTeachingController,
	}

	routeConfig.Setup()
}
