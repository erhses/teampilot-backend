package admin

// import (
// 	"teampilot/integrations/dba"
// 	"teampilot/routes/controller/common"
// 	"teampilot/routes/middleware"
// 	"teampilot/structs/dts"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// )

// type SpaceController struct {
// 	common.Controller
// }

// func (co SpaceController) Register(router fiber.Router) {
// 	router.Post("/create-space", co.CreateSpace)
// 	router.Get("/:id", co.GetSpace)
// 	router.Post("/page", middleware.AuthMiddleware, co.Pagination)
// 	router.Delete("/:id", co.Delete)
// 	router.Put("/:id", co.Update)
// 	router.Post("/apps", co.Apps)
// }

// // List of Applications
// // @Summary      бүх талбайн дээр ирсэн хүсэлтийн жагсаалт
// // @Description  бүх талбайн дээрх хүсэлтийн жагсаалт
// // @Tags         Owner Space
// // @Accept       json
// // @Produce      json
// // @Success      200      {array}   common.BaseResponse{body=dts.Application}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/space/apps [post]
// func (co SpaceController) Apps(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var apps []dts.Application
// 	if err := dba.DB.Find(&apps).Error; err != nil {
// 		return co.SetError(c, http.StatusNotFound, err)
// 	}
// 	return co.SetBody(apps)
// }

// // Create Space
// // @Summary      талбай үүсгэх
// // @Description  шинээр талбай үүсгэх
// // @Tags         Admin Space
// // @Accept       json
// // @Produce      json
// // @Param        Space    body      dts.Space true   "Space details"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/space/create-space [post]
// func (co SpaceController) CreateSpace(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var space dts.Space
// 	if err := c.BodyParser(&space); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Create(&space).Error; err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Space by ID
// // @Summary      талбайн мэдээлэл харах
// // @Description  өгөгдсөн ID-тай талбайн мэдээлэл
// // @Tags         Admin Space
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "space id"
// // @Success      200      {object}   common.BaseResponse{body=dts.Space}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security  	 BearerAuth
// // @Router       /admin/space/{id} [get]
// func (co SpaceController) GetSpace(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var space dts.Space
// 	if err := dba.DB.Preload("Property").Preload("User").First(&space, id).Error; err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	return co.SetBody(space)
// }

// type SpaceFilter struct {
// 	dba.PaginationInput
// }

// // List of spaces
// // @Summary      талбайн жагсаалт
// // @Description  талбайн жагсаалт filter-тэй
// // @Tags         Admin Space
// // @Accept       json
// // @Produce      json
// // @Param        input   body   PropertyFilter     true  "Propert filter"
// // @Success      200      {array}   common.BaseResponse{body=dts.Space}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security  	 BearerAuth
// // @Router       /admin/space/page [post]
// func (co SpaceController) Pagination(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var params SpaceFilter
// 	if err := c.BodyParser(&params); err != nil {
// 		return co.SetError(c, fiber.StatusBadRequest, err)
// 	}

// 	orm := dba.QueryBuilder{
// 		DB: dba.DB.Model(&dts.Space{}),
// 	}
	
// 	var result common.PaginationResult
// 	result.Total = orm.Total()

// 	var space []dts.Space
// 	if err := orm.Scopes(dba.Paginate(&params.PaginationInput)).Find(&space).Error; err != nil {
// 		return co.SetError(c, fiber.StatusNotFound, err)
// 	}
// 	result.Items = space
// 	return co.SetBody(result)
// }

// // Delete Space
// // @Summary      талбай устгах
// // @Description  өгөгдсөг ID-тай талбайг устгах
// // @Tags         Owner Space
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "space id"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/space/{id} [delete]
// func (co SpaceController) Delete(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var space dts.Space
// 	if err := dba.DB.Delete(&space, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Update Space
// // @Summary      талбай шинэчлэх
// // @Description  өгөгдсөг ID-тай талбайн мэлээллийг шинэчлэх
// // @Tags         Owner Space
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "space id"
// // @Param        Space    body      dts.Space true   "Space details to update"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/space/{id} [put]
// func (co SpaceController) Update(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var space dts.Space
// 	dba.DB.First(&space, id)
// 	if err := c.BodyParser(&space); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Save(&space).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }
