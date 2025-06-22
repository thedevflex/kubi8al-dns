package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Setup(app *fiber.App) {
	app.Use(cors.New(
		cors.Config{

			AllowMethods: "GET,POST,DELETE",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "ok",
			"message": "healthy",
		})
	})
	app.Use("*", func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"message": "you missed the server",
		})
	})
}
