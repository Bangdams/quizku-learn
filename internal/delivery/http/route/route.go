package route

import (
	"github.com/Bangdams/quizku-learn/internal/delivery/http"
	"github.com/Bangdams/quizku-learn/internal/util"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	UserController   http.UserController
	CourseController http.CourseController
	ClassController  http.ClassController
}

func (config *RouteConfig) Setup() {
	// API for user
	config.App.Get("/api/users", util.CheckLevel("admin"), config.UserController.FindAll)
	config.App.Get("/api/users/:email", util.CheckLevel("admin"), config.UserController.FindByEmail)
	config.App.Post("/api/users", util.CheckLevel("admin"), config.UserController.Create)
	config.App.Delete("/api/users/:id", util.CheckLevel("admin"), config.UserController.Delete)
	config.App.Put("/api/users", util.CheckLevel("admin"), config.UserController.Update)

	// API for course
	config.App.Get("/api/courses", util.CheckLevel("admin"), config.CourseController.FindAll)
	config.App.Get("/api/courses/:course_code", util.CheckLevel("admin"), config.CourseController.FindByCourseCode)
	config.App.Post("/api/courses", util.CheckLevel("admin"), config.CourseController.Create)
	config.App.Delete("/api/courses/:course_code", util.CheckLevel("admin"), config.CourseController.Delete)
	config.App.Put("/api/courses", util.CheckLevel("admin"), config.CourseController.Update)

	// API for class
	config.App.Get("/api/classes", util.CheckLevel("admin"), config.ClassController.FindAll)
	config.App.Get("/api/classes/:class_name", util.CheckLevel("admin"), config.ClassController.FindByName)
	config.App.Post("/api/classes", util.CheckLevel("admin"), config.ClassController.Create)
	config.App.Delete("/api/classes/:class_id", util.CheckLevel("admin"), config.ClassController.Delete)
	config.App.Put("/api/classes", util.CheckLevel("admin"), config.ClassController.Update)

	// Api for login
	config.App.Post("/login", config.UserController.Login)
	config.App.Post("/logout", config.UserController.Logout)
	config.App.Post("/refresh", config.UserController.Refresh)
	config.App.Get("/api/status-login", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"message": "success"})
	})
}
