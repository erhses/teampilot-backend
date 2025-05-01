package admin

// import (
// 	"teampilot/integrations/dba"
// 	"teampilot/routes/controller/common"
// 	"teampilot/structs/dts"
// 	"teampilot/utils"
// 	"errors"
// 	"net/http"

// 	"github.com/gofiber/fiber/v2"
// )

// type UserController struct {
// 	common.Controller
// }

// func (co UserController) Register(router fiber.Router) {
// 	router.Post("/register", co.SignUp).Name("create")
// 	router.Get("/:id", co.Get)
// 	router.Post("/page", co.Pagination)
// 	router.Delete("/:id", co.Delete)
// 	router.Put("/:id", co.Update)

// }

// // Signup
// // @Summary      шинэ хэрэглэгч үүсгэх
// // @Description  хэрэгтэй мэдээлэл оруулан шинэ хэрэглэгч бүртгүүлэх, JWT болон хэрэглэгчийн мэдээлэл буцаах
// // @Tags         Admin User
// // @Accept       json
// // @Produce      json
// // @Param        body    body      SignupInput  true   "user registration details"
// // @Success      200      {object}   common.BaseResponse{body=SignupResponse}
// // @Failure      400   {object}  common.BaseResponse "Invalid request body"
// // @Failure      409   {object}  common.BaseResponse "User already exists"
// // @Failure      500   {object}  common.BaseResponse "Internal server error"
// // @Router       /admin/user/register [post]
// func (co UserController) SignUp(c *fiber.Ctx) error {
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
// 	if params.Email == "" || params.Password == "" {
// 		return co.SetError(c, http.StatusBadRequest, errors.New("email and password are required"))
// 	}
// 	var existingUser dts.User
// 	result := dba.DB.Where("email = ?", params.Email).First(&existingUser)
// 	if result.Error == nil {
// 		return co.SetError(c, http.StatusConflict, errors.New("хэрэглэгч бүртгэлтэй байна"))
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
// 	}
// 	if err := dba.DB.Create(&newUser).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, errors.New("хэрэглэгч бүртгэх үед алдаа гарлаа"))
// 	}
// 	token := utils.GenerateToken(newUser.ID, string(newUser.UserType))
// 	return co.SetBody(SignupResponse{
// 		Token: token,
// 		User:  newUser,
// 	})
// }

// // Get users by ID
// // @Summary      хэрэглэгчийн мэдээлэл харах
// // @Description  өгөгдсөн ID-тай хэрэглэгчийн мэдээлэл харах
// // @Tags         Admin User
// // @Accept       json
// // @Produce      json
// // @Param        id    path      int true   "user id"
// // @Success      200   {object}   common.BaseResponse{body=dts.User}
// // @Failure      400   {object}  common.BaseResponse "Invalid request body"
// // @Failure      409   {object}  common.BaseResponse "User already exists"
// // @Failure      500   {object}  common.BaseResponse "Internal server error"
// // @Router       /admin/user/{id} [get]
// func (co UserController) Get(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var user dts.User
// 	if err := dba.DB.First(&user, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(user)
// }

// type UserFilter struct {
// 	Limit int `json:"limit"`
// 	Page  int `json:"page"`
// }

// // List of Owners
// // @Summary      хэрэглэгчдийн жагсаалт
// // @Description  хэрэглэгчидийн жагсаалт filter-тэй
// // @Tags         Admin User
// // @Accept       json
// // @Produce      json
// // @Param        input   body   UserFilter     true  "User filter"
// // @Success      200      {array}   common.BaseResponse{body=dts.User}
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/user/page [post]
// func (co UserController) Pagination(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()

// 	var users []dts.User
// 	if err := dba.DB.Find(&users).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(users)
// }

// // Delete user by ID
// // @Summary      хэрэглэгч устгах
// // @Description  өгөгдсөн ID-тай хэрэглэгчийн устгах
// // @Tags         Admin User
// // @Accept       json
// // @Produce      json
// // @Param        id 	   path  	int     true  "User ID"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/user/{id} [delete]
// func (co UserController) Delete(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var user dts.User
// 	if err := dba.DB.Delete(&user, id).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }

// // Delete user by ID
// // @Summary      хэрэглэгчийн мэдээлэл шинэчлэх
// // @Description  өгөгдсөн ID-тай хэрэглэгчийн мэдээлэл шинэчлэх
// // @Tags         Admin User
// // @Accept       json
// // @Produce      json
// // @Param        id 	   path  	int     true  "User ID"
// // @Success      200      {object}   common.SuccessResponse
// // @Failure      400      {object}   common.BaseResponse
// // @Failure      401      {object}   common.BaseResponse
// // @Failure      403      {object}   common.BaseResponse
// // @Failure      404      {object}   common.BaseResponse
// // @Failure      500      {object}   common.BaseResponse
// // @Security     BearerAuth
// // @Router       /admin/user/{id} [put]
// func (co UserController) Update(c *fiber.Ctx) error {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			co.RespondPanic(c, r)
// 		} else {
// 			co.GetBody(c)
// 		}
// 	}()
// 	id := c.Params("id")
// 	var user dts.User
// 	dba.DB.First(&user, id)
// 	if err := c.BodyParser(&user); err != nil {
// 		return co.SetError(c, http.StatusBadRequest, err)
// 	}
// 	if err := dba.DB.Save(&user).Error; err != nil {
// 		return co.SetError(c, http.StatusInternalServerError, err)
// 	}
// 	return co.SetBody(common.SuccessResponse{
// 		Success: true,
// 	})
// }
