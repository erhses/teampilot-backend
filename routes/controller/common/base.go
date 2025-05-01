package common

import (
	"fmt"
	"teampilot/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Controller struct {
	Response *BaseResponse
}

func (co Controller) GetBody(c *fiber.Ctx) {
	if co.Response.Message != "" {
		return
	}
	c.Status(200).JSON(co.Response)
}

func (co Controller) SetBody(body interface{}) error {
	co.Response.Code = 200
	co.Response.Message = ""
	co.Response.Body = body
	return nil
}

func (co Controller) SetError(c *fiber.Ctx, code int, err error) error {
	co.Response.Code = 500
	co.Response.Message = err.Error()
	if err != nil {
		utils.Error(c.BaseURL(), zap.String("Error", err.Error()))
	}
	if code == 0 {
		co.Response.Code = fiber.StatusInternalServerError
	} else {
		co.Response.Code = code
	}
	co.Response.Body = nil
	return c.Status(co.Response.Code).JSON(co.Response)
}

func (co Controller) RespondPanic(c *fiber.Ctx, err interface{}) error {
	utils.Error("PANIC", zap.String("", fmt.Sprint(err)))
	return c.Status(500).JSON(&fiber.Map{
		"body":    nil,
		"message": "panic error happend",
	})
}
