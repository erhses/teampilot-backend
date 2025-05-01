package routes

import (
	"teampilot/integrations/dba"
	"teampilot/integrations/rdb"
	"teampilot/routes/controller"
	"teampilot/routes/controller/common"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Register(app *fiber.App) *fiber.App {
	bc := common.Controller{
		Response: &common.BaseResponse{
			Message: "",
			Body:    nil,
			Code:    200,
		},
	}
	// Middlewares
	app.Use(cors.New())
	app.Use(recover.New())

	dba.LoadDatabase()
	rdb.InitRedis()

	AliveController{Controller: bc}.Register(app)
	controller.Register(bc, app.Group("/api/v1"))
	app.Hooks().OnShutdown(func() error {
		return rdb.CloseRedis()
	})
	return app
}
