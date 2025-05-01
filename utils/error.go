package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	message := err.Error()
	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	// Send custom error page
	ctx.Status(code).JSON(fiber.Map{
		"message": message,
		"body":    nil,
	})
	// Return from handler
	return nil
}

func Info(name string, err zap.Field) {
	ZapLogger.Info(name, err)
}

func Error(name string, err zap.Field) {
	ZapLogger.Error(name, err)
}
