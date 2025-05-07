package teacher

import (
	"net/http"
	"teampilot/integrations/dba"
	"teampilot/routes/controller/common"
	"teampilot/routes/middleware"
	"teampilot/structs/dts"

	"github.com/gofiber/fiber/v2"
)

type (
	ClassController struct {
		common.Controller
	}
	ClassInput struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Semester    string `json:"semester"`
		TeacherID   uint   `json:"teacher_id"`
		// Students    []User `json:"students"`
	}
)

func (co ClassController) Register(router fiber.Router) {
	router.Get("/classes", co.GetClasses)
	router.Get("/class/:id", co.GetClass)
	router.Post("/create", co.CreateClass)
	router.Delete("/delete/:id", co.DeleteClass)
	router.Put("/update/:id", co.UpdateClass)
}
func (co ClassController) GetClass(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	id := c.Params("id")
	var class dts.Class
	if err := dba.DB.First(&class, id).Error; err != nil {
		return co.SetError(c, http.StatusNotFound, err)
	}
	return co.SetBody(class)
}

func (co ClassController) GetClasses(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	user := middleware.RetrieveAuth(c)
	var classes []dts.Class
	if err := dba.DB.Where("teacher_id = ?", user.ID).Find(&classes).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(classes)
}

func (co ClassController) CreateClass(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	user := middleware.RetrieveAuth(c)
	var class ClassInput
	if err := c.BodyParser(&class); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}
	class.TeacherID = user.ID
	if err := dba.DB.Create(&class).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(common.SuccessResponse{Success: true})
}

func (co ClassController) UpdateClass(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	id := c.Params("id")
	user := middleware.RetrieveAuth(c)
	var class dts.Class
	if err := dba.DB.Where("id = ? AND teacher_id = ?", id, user.ID).First(&class).Error; err != nil {
		return co.SetError(c, http.StatusNotFound, err)
	}

	if err := c.BodyParser(&class); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}
	if err := dba.DB.Save(&class).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(common.SuccessResponse{Success: true})
}

func (co ClassController) DeleteClass(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	id := c.Params("id")
	user := middleware.RetrieveAuth(c)
	if err := dba.DB.Where("id = ? AND teacher_id = ?", id, user.ID).Delete(&dts.Class{}).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(common.SuccessResponse{Success: true})
}
