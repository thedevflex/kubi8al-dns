package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	logs "github.com/thedevflex/kubi8al-webhook/utils/logger"
)

func New() *fiber.App {
	return fiber.New(fiber.Config{
		AppName:               "kubi8al-dns",
		ErrorHandler:          customErrorHandler,
		DisableStartupMessage: true,
	})
}

func Start(app *fiber.App) error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logs.Infof("Starting server on port %s", port)
	return app.Listen(":" + port)
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	logs.Errorf("Error handling request: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": message,
	})
}
