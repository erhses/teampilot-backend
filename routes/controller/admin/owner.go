package admin

// import (
// 	"teampilot/constants"
// 	"teampilot/integrations/dba"
// 	"teampilot/routes/controller/common"
// 	"teampilot/structs/dts"
// 	"teampilot/utils"
// 	"errors"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// )

// type OwnerController struct {
// 	common.Controller
// }
// type NotificationRequest struct {
// 	Title   string `json:"title"`
// 	Content string `json:"content"`
// }

// func (co OwnerController) Register(router fiber.Router) {
// 	router.Post("/create-owner", co.Create).Name("create")
// 	router.Get("/:id", co.Get)
// 	router.Post("/page", co.Pagination)
// 	router.Delete("/:id", co.Delete)
// 	router.Put("/:id", co.Update)
// }

// type (
// 	SignupInput struct {
// 		Name      string             `json:"name"`
// 		Email     string             `json:"email"`
// 		Password  string             `json:"password"`
// 		UserType  constants.UserType `json:"user_type"`
// 		IsCompany bool               `json:"is_company"`
// 		Register  string             `json:"register"`
// 	}
// 	SignupResponse struct {
// 		Token string   `json:"token"`
// 		User  dts.User `json:"user"`
// 	}
// )

// // Register an owner
// // @Summary      шинээр эзэмшигч үүсгэх
// // @Description  хэрэгтэй мэдээлэл оруулан шинэ эзэмшигч бүртгүүлэх, JWT болон хэрэглэгчийн мэдээлэл буцаах
// // @Tags         Admin Owner
// // @Accept       json
// // @Produce      json
// // @Param        body    body      SignupInput  true   "user registration details"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security	 BearerAuth
// // @Router       /admin/owner/create-owner [post]
// func (co OwnerController) Create(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	var params SignupInput
// 	if err := c.BodyParser(&params); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	hashedPassword, err := utils.GenerateHash(params.Password)
// 	if err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, errors.New("нууц үг үүсгэхэд алдаа гарлаа"))
// 	}
// 	newUser := dts.User{
// 		Name:      params.Name,
// 		Email:     params.Email,
// 		Password:  hashedPassword,
// 		UserType:  params.UserType,
// 		IsActive:  true,
// 		IsCompany: params.IsCompany,
// 		Register:  params.Register,
// 	}
// 	if err := dba.DB.Create(&newUser).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	// fmt.Println(params)
// 	token := utils.GenerateToken(newUser.ID, string(newUser.UserType))
// 	return co.SetBody(SignupResponse{
// 		Token: token,
// 		User:  newUser,
// 	})
// }

// // Owner by ID
// // @Summary      эзэмшигчийн мэдээлэл харах
// // @Description  өгөгдсөн ID-тай эзэмшигчийн мэдээлэл
// // @Tags         Admin Owner
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "user id"
// // @Success      200      {object}   common.BaseResponse{body=dts.User}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security 	 BearerAuth
// // @Router       /admin/owner/{id} [get]
// func (co OwnerController) Get(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var owner dts.User
// 	if err := dba.DB.First(&owner, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(owner)
// }

// type OwnerFilter struct {
// 	dba.PaginationInput
// }

// // List of Owners
// // @Summary      эзэмшигчидийн жагсаалт
// // @Description  эзэмшигчидийн жагсаалт filter-тэй
// // @Tags         Admin Owner
// // @Accept       json
// // @Produce      json
// // @Param        input   body   OwnerFilter     true  "Owner filter"
// // @Success      200      {array}   common.BaseResponse{body=common.PaginationResult{items=[]dts.User}}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security 	 BearerAuth
// // @Router       /admin/owner/page [post]
// func (co OwnerController) Pagination(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()

// 	var params OwnerFilter
// 	if err := c.BodyParser(&params); err != nil {
// 		return co.SetError(c, fiber.StatusBadRequest, err)
// 	}

// 	orm := dba.QueryBuilder{
// 		DB: dba.DB.Model(&dts.User{}),
// 	}

// 	orm = *orm.Equal("user_type", "owner")

// 	var result common.PaginationResult

// 	result.Total = orm.Total()
// 	var owners []dts.User
// 	if err := orm.Scopes(dba.Paginate(&params.PaginationInput)).Find(&owners).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	result.Items = owners
// 	return co.SetBody(result)
// }

// // Delete an owner
// // @Summary      эзэмшигч устгах
// // @Description  өгөгдсөн ID-тай эзэмшигчийг устгах
// // @Tags         Admin Owner
// // @Accept       json
// // @Produce      json
// // @Param        id   path  int     true  "Owner ID"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security 	 BearerAuth
// // @Router       /admin/owner/{id} [delete]
// func (co OwnerController) Delete(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var owner dts.User
// 	if err := dba.DB.Delete(&owner, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Update an owner
// // @Summary      эзэмшигчийн мэдээлэл шинэчлэх
// // @Description  өгөгдсөн ID-тай эзэмшигчийн мэдээллийг шинэчлэх
// // @Tags         Admin Owner
// // @Accept       json
// // @Produce      json
// // @Param        id   path  int     true  "Owner ID"
// // @Success      200      {array}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security 	 BearerAuth
// // @Router       /admin/owner/{id} [put]
// func (co OwnerController) Update(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var owner dts.User
// 	dba.DB.First(&owner, id)
// 	if err := c.BodyParser(&owner); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Save(&owner).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }
