package admin

// import (
// 	"teampilot/routes/controller/common"
// 	"teampilot/routes/middleware"

// 	"github.com/gofiber/fiber/v2"
// )

// func Register(bc common.Controller, router fiber.Router) {

// 	OwnerController{Controller: bc}.Register(router.Group("/owner").Use(middleware.AuthMiddleware))
// 	TenantController{Controller: bc}.Register(router.Group("/tenant").Use(middleware.AuthMiddleware))
// 	UserController{Controller: bc}.Register(router.Group("/user").Use(middleware.AuthMiddleware))
// 	RBACController{Controller: bc}.Register(router.Group("/rbac").Use(middleware.AuthMiddleware))
// 	protected := router.Group("", middleware.AuthMiddleware, middleware.RBACMiddleware)
// 	PropertyController{Controller: bc}.Register(protected.Group("/prop"))
// 	SpaceController{Controller: bc}.Register(protected.Group("/space"))
// }
