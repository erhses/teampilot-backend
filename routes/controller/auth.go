package controller

import (
	"teampilot/constants"
	"teampilot/integrations"
	"teampilot/integrations/dba"
	"teampilot/routes/controller/common"
	"teampilot/routes/middleware"
	"teampilot/structs/dts"
	"teampilot/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	common.Controller
}
const authKey = "key"
func (co AuthController) Register(router fiber.Router) {
	router.Post("/login", co.Login)
	router.Post("/register", co.SignUp)
	router.Post("/signout", middleware.AuthMiddleware, co.SignOut)
}

type (
	LoginInput struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		UserType  string `json:"user_type"`
		NotiToken string `json:"noti_token"`
	}
	LoginResponse struct {
		Token string   `json:"token"`
		User  dts.User `json:"data"`
	}
	SignupInput struct {
		Name      string             `json:"name"`
		Email     string             `json:"email"`
		Password  string             `json:"password"`
		UserType  constants.UserType `json:"user_type"`
		IsCompany bool               `json:"is_company"`
		Register  string             `json:"register"`
	}
	SignupResponse struct {
		Token string   `json:"token"`
		User  dts.User `json:"user"`
	}
)

// Login
// @Summary      мэдээлэлээр нэвтрэх
// @Description  имэйл хаяг, нууц үг, хэрэглэгчийн төрөлөөр нэвтэрж, JWT болон хэрэглэгчийн мэдээлэл буцаах
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body    body      LoginInput  true   "login credentials"
// @Success      200      {object}   common.BaseResponse{body=LoginResponse}
// @Failure      400      {object}   common.BaseResponse
// @Failure      401      {object}   common.BaseResponse
// @Failure      403      {object}   common.BaseResponse
// @Failure      404      {object}   common.BaseResponse
// @Failure      500      {object}   common.BaseResponse
// @Router       /auth/login [post]
func (co AuthController) Login(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	var params LoginInput
	if err := c.BodyParser(&params); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}

	var user dts.User
	if err := dba.DB.First(&user, "email", params.Email).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, errors.New("хэрэглэгчийн имэйл, нууц үг буруу байна"))
	}

	if err := utils.ComparePassword(user.Password, params.Password); err != nil {
		return co.SetError(c, http.StatusInternalServerError, errors.New("хэрэглэгчийн имэйл, нууц үг буруу байна"))
	}
	if err := dba.DB.Model(&dts.User{
		Base: dts.Base{
			ID: user.ID,
		},
	}).Update("notification_token", params.NotiToken).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)

	}
	integrations.SubscribeToTopic(params.NotiToken, string(user.UserType))
	return co.SetBody(LoginResponse{
		Token: utils.GenerateToken(user.ID, string(user.UserType)),
		User:  user,
	})
}

// Signup
// @Summary      шинэ хэрэглэгч үүсгэх
// @Description  хэрэгтэй мэдээлэл оруулан шинэ хэрэглэгч бүртгүүлэх, JWT болон хэрэглэгчийн мэдээлэл буцаах
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body    body      SignupInput  true   "user registration details"
// @Success      200   {object}   common.BaseResponse{body=SignupResponse}
// @Failure      400   {object}  common.BaseResponse "Invalid request body"
// @Failure      409   {object}  common.BaseResponse "User already exists"
// @Failure      500   {object}  common.BaseResponse "Internal server error"
// @Router       /auth/register [post]
func (co AuthController) SignUp(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	var params SignupInput
	if err := c.BodyParser(&params); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}
	if params.Email == "" || params.Password == "" {
		return co.SetError(c, http.StatusBadRequest, errors.New("email and password are required"))
	}
	var existingUser dts.User
	result := dba.DB.Where("email = ?", params.Email).First(&existingUser)
	if result.Error == nil {
		return co.SetError(c, http.StatusConflict, errors.New("хэрэглэгч бүртгэлтэй байна"))
	}
	hashedPassword, err := utils.GenerateHash(params.Password)
	if err != nil {
		return co.SetError(c, http.StatusInternalServerError, errors.New("нууц үг үүсгэхэд алдаа гарлаа"))
	}
	newUser := dts.User{
		Name:      params.Name,
		Email:     params.Email,
		Password:  hashedPassword,
		UserType:  params.UserType,
		IsActive:  true,
	}
	if err := dba.DB.Create(&newUser).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, errors.New("хэрэглэгч бүртгэх үед алдаа гарлаа"))
	}
	token := utils.GenerateToken(newUser.ID, string(newUser.UserType))
	return co.SetBody(SignupResponse{
		Token: token,
		User:  newUser,
	})
}

// SignOut
// @Summary      хэрэглэгч системээс гарах
// @Description  хэрэглэгчийн notification token-ийг хоослох, topic-оос устгах, token-ийг blacklist-д нэмэх
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}   common.BaseResponse
// @Failure      401      {object}   common.BaseResponse
// @Failure      500      {object}   common.BaseResponse
// @Router       /auth/signout [post]
func (co AuthController) SignOut(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()

	curr := middleware.RetrieveAuth(c)
	token := c.Get("Authorization")[7:]
	var user dts.User
	if err := dba.DB.First(&user, curr.ID).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}

	if user.NotificationToken != "" {
		integrations.UnsubscribeFromTopic(user.NotificationToken, string(user.UserType))
	}

	if err := dba.DB.Model(&user).Update("notification_token", "").Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, errors.New("гарах үед алдаа гарлаа"))
	}
	// c.Locals(authKey, nil)
	tokenManager := utils.GetTokenManager()
    if err := tokenManager.BlacklistToken(c.Context(), token); err != nil {
        return co.SetError(c, http.StatusInternalServerError, errors.New("token invalidation failed"))
    }

	return co.SetBody(common.SuccessResponse{Success: true})
}
