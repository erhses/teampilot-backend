package admin

// import (
// 	"teampilot/integrations/dba"
// 	"teampilot/routes/controller/common"
// 	"teampilot/structs/dts"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// )

// type TenantController struct {
// 	common.Controller
// }

// func (co TenantController) Register(router fiber.Router) {
// 	router.Delete("/:id", co.DeleteTenant)
// }

// // Delete User
// // @Summary      хэрэглэгч устгах
// // @Description  өгөгдсөг ID-тай хэрэглэгчийг устгах
// // @Tags         Admin Tenant
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "user id"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/tenant/{id} [delete]
// func (co TenantController) DeleteTenant(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var user dts.User
// 	id := c.Params("id")
// 	if err := dba.DB.Where("id = ?", id).Delete(&user).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{Success: true})
// }
