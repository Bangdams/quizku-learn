package route

import (
	"github.com/Bangdams/quizku-learn/internal/delivery/http"
	"github.com/Bangdams/quizku-learn/internal/util"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController http.UserController
}

func (config *RouteConfig) Setup() {
	// API for user
	config.App.Get("/api/users", util.CheckLevel("admin"), config.UserController.FindAll)
	config.App.Get("/api/users/:email", util.CheckLevel("admin"), config.UserController.FindByEmail)
	config.App.Post("/api/users", util.CheckLevel("admin"), config.UserController.Create)
	config.App.Delete("/api/users/:id", util.CheckLevel("admin"), config.UserController.Delete)
	config.App.Put("/api/users", util.CheckLevel("admin"), config.UserController.Update)
	config.App.Put("/api/users/change-password", util.CheckLevel("admin"), config.UserController.UpdatePassword)

	// Api for login
	config.App.Post("/login", config.UserController.Login)
	config.App.Get("/api/status-login", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"message": "success"})
	})
}
