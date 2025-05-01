package config

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewFiber() *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      "Quizku API",
		ErrorHandler: NewErrorHandler(),
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://127.0.0.1:5500/",
		AllowCredentials: true,
	}))

	app.Use("/api", jwtware.New(jwtware.Config{
		TokenLookup: "cookie:token",
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("SECRET_KEY")),
		},
		ContextKey: "user",
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
