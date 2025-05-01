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

	// usecase
	userUsecase := usecase.NewUserUsecase(userRepo, config.DB, config.Validate)

	// controller
	userController := http.NewUserController(userUsecase)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
	}

	routeConfig.Setup()
}
