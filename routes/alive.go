package routes

import (
	"teampilot/routes/controller/common"

	"github.com/gofiber/fiber/v2"
)

type AliveController struct {
	common.Controller
}

func (co AliveController) Register(router *fiber.App) {

}
