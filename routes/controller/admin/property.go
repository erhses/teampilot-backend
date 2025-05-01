package admin

// import (
// 	"teampilot/integrations/dba"
// 	"teampilot/routes/controller/common"
// 	"teampilot/structs/dts"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// )

// type PropertyController struct {
// 	common.Controller
// }

// func (co PropertyController) Register(router fiber.Router) {
// 	router.Post("/create", co.Create)
// 	router.Get("/:id", co.GetProperty)
// 	router.Post("/page", co.Pagination)
// 	router.Delete("/:id", co.Delete)
// 	router.Put("/:id", co.Update)
// }

// type PropertyInput struct {
// 	Name string
// }

// // Create Property
// // @Summary      үл хөдлөх үүсгэх
// // @Description  өгөгдсөн мэдээлэлээх үл хөдлөх үүсгэх
// // @Tags         Admin Property
// // @Accept       json
// // @Produce      json
// // @Param        property   body   dts.Property     true  "property details"
// // @Success      200  {object}  common.SuccessResponse
// // @Failure      400  {object}  common.BaseResponse
// // @Failure      404  {object}  common.BaseResponse
// // @Failure      500  {object}  common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/prop/create [post]
// func (co PropertyController) Create(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var params PropertyInput
// 	if err := c.BodyParser(&params); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Create(&dts.Property{
// 		Name: params.Name,
// 	}).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Property
// // @Summary      үл хөдлөх
// // @Description  өгөгдсөн ID-тай үл хөдлөх
// // @Tags         Admin Property
// // @Accept       json
// // @Produce      json
// // @Param        id   path   int    true  "property id"
// // @Success      200  {object}  common.BaseResponse{body=dts.Property}
// // @Failure      400  {object}  common.BaseResponse
// // @Failure      404  {object}  common.BaseResponse
// // @Failure      500  {object}  common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/prop/{id} [get]
// func (co PropertyController) GetProperty(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var property dts.Property
// 	if err := dba.DB.First(&property, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(property)
// }

// type PropertyFilter struct {
// 	dba.PaginationInput
// }

// // Propert List
// // @Summary      Үл хөдлөхийн жагсаалт
// // @Description  Үл хөдлөхийн жагсаалт filter-тэй
// // @Tags         Owner Property
// // @Accept       json
// // @Produce      json
// // @Param        input   body   PropertyFilter     true  "Propert filter"
// // @Success      200  {array}  common.BaseResponse{body=dts.Property}
// // @Failure      400  {object}  common.BaseResponse
// // @Failure      404  {object}  common.BaseResponse
// // @Failure      500  {object}  common.BaseResponse
// // @Router       /admin/prop/page [post]
// func (co PropertyController) Pagination(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
	
// 	var params PropertyFilter
// 	if err := c.BodyParser(&params); err != nil {
// 		return co.SetError(c, fiber.StatusBadRequest, err)
// 	}
// 	orm := dba.QueryBuilder{
// 		DB: dba.DB.Model(&dts.Property{}),
// 	}
// 	var result common.PaginationResult

// 	result.Total = orm.Total()

// 	var properties []dts.Property
// 	if err := orm.Scopes(dba.Paginate(&params.PaginationInput)).Find(&properties).Error; err != nil {
// 		return co.SetError(c, fiber.StatusNotFound, err)
// 	}
// 	result.Items = properties
// 	return co.SetBody(result)
// }

// // Delete Property
// // @Summary      үл хөдлөх устгах
// // @Description  өгөгдсөн ID-тай үл хөдлөхийг устгах
// // @Tags         Owner Property
// // @Accept       json
// // @Produce      json
// // @Param        id        path      int          true   "Property ID"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/prop/{id} [delete]
// func (co PropertyController) Delete(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var property dts.Property
// 	if err := dba.DB.Delete(&property, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Update Property
// // @Summary      үл хөдлөх шинэчлэх
// // @Description  өгөгдсөн ID-тай үл хөдлөхийн мэдээллийг шинэчлэх
// // @Tags         Owner Property
// // @Accept       json
// // @Produce      json
// // @Param        id        path      int          true   "Property ID"
// // @Param        property  body      dts.Property true   "Property details to update"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/prop/{id} [put]
// func (co PropertyController) Update(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var property dts.Property
// 	dba.DB.First(&property, id)
// 	if err := c.BodyParser(&property); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Save(&property).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{Success: true})
// }
