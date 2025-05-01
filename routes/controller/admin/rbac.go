package admin

import (
	"teampilot/integrations/dba"
	"teampilot/routes/controller/common"
	"teampilot/routes/middleware"
	"teampilot/structs/dts"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type RBACController struct {
	common.Controller
}
func(co RBACController) Register (router fiber.Router){
	router.Post("/role", co.CreateRole)
	router.Post("/roles", co.ListRoles)
	router.Post("/group", co.CreateGroup)
	router.Post("/groups", co.ListGroups)
	router.Post("/permission", co.CreatePermission)
	router.Post("/permissions", co.ListPermissions)
	router.Post("/assign/:id", co.AssignPermissionsToGroup)
	router.Post("/assignrole", co.AssignRole)
	router.Post("/role/:id", co.AssignGroupToRole)
	router.Post("/init", co.Init)
	router.Post("/see", co.All)
}

func (co RBACController) Init(c *fiber.Ctx) error {
	return common.InitializeRBAC()
}

func (co RBACController) All(c *fiber.Ctx) error {
	var roles []dts.Group
	if err := dba.DB.Find(&roles).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(roles)
}

func (co RBACController) CreateRole(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()

	var role dts.Role
	if err := c.BodyParser(&role); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}

	if err := dba.DB.Create(&role).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}

	return co.SetBody(role)
}
func (co RBACController) ListRoles(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
    var roles []dts.Role
    if err := dba.DB.Preload("Groups").Find(&roles).Error; err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }
    return co.SetBody(roles)
}
func (co RBACController) CreateGroup(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	var group dts.Group
	if err := c.BodyParser(&group); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}
	if err := dba.DB.Create(&group).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(group)
}
func (co RBACController) ListGroups(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
    var groups []dts.Group
    if err := dba.DB.Preload("Permissions").Find(&groups).Error; err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }
    return co.SetBody(groups)
}
func (co RBACController) CreatePermission(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
	var permission dts.Permission
	if err := c.BodyParser(&permission); err != nil {
		return co.SetError(c, http.StatusBadRequest, err)
	}
	if err := dba.DB.Create(&permission).Error; err != nil {
		return co.SetError(c, http.StatusInternalServerError, err)
	}
	return co.SetBody(permission)
}
func (co RBACController) ListPermissions(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
    var permissions []dts.Permission
    if err := dba.DB.Find(&permissions).Error; err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }
    return co.SetBody(permissions)
}
func (co RBACController) AssignPermissionsToGroup(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
    groupID := c.Params("id")
    var permissionIDs []uint
    if err := c.BodyParser(&permissionIDs); err != nil {
        return co.SetError(c, http.StatusBadRequest, err)
    }

    var group dts.Group
    if err := dba.DB.First(&group, groupID).Error; err != nil {
        return co.SetError(c, http.StatusNotFound, err)
    }

    var permissions []dts.Permission
    if err := dba.DB.Find(&permissions, permissionIDs).Error; err != nil {
        return co.SetError(c, http.StatusNotFound, err)
    }

    if err := dba.DB.Model(&group).Association("Permissions").Replace(&permissions); err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }

    return co.SetBody(common.SuccessResponse{Success: true})
}

func (co RBACController) AssignGroupToRole(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			co.RespondPanic(c, r)
		} else {
			co.GetBody(c)
		}
	}()
    roleID := c.Params("id")
    var groupIDs []uint
    if err := c.BodyParser(&groupIDs); err != nil {
        return co.SetError(c, http.StatusBadRequest, err)
    }

    var role dts.Role
    if err := dba.DB.First(&role, roleID).Error; err != nil {
        return co.SetError(c, http.StatusNotFound, err)
    }

    var groups []dts.Group
    if err := dba.DB.Find(&groups, groupIDs).Error; err != nil {
        return co.SetError(c, http.StatusNotFound, err)
    }

    if err := dba.DB.Model(&role).Association("Groups").Replace(&groups); err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }

    return co.SetBody(common.SuccessResponse{Success: true})
}

func (co RBACController) AssignRole(c *fiber.Ctx) error {
    defer func() {
        if r := recover(); r != nil {
            co.RespondPanic(c, r)
        }
    }()
    
    current := middleware.RetrieveAuth(c)
    
    var request struct {
        RoleID uint `json:"role_id"`
    }
    
    // Parse the request body
    if err := c.BodyParser(&request); err != nil {
        return co.SetError(c, http.StatusBadRequest, err)
    }
    
    // Check if role exists
    var role dts.Role
    if err := dba.DB.First(&role, request.RoleID).Error; err != nil {
        return co.SetError(c, http.StatusNotFound, err)
    }
    
    // Update user's role
    current.RoleID = request.RoleID
    if err := dba.DB.Save(&current).Error; err != nil {
        return co.SetError(c, http.StatusInternalServerError, err)
    }
    
    return co.SetBody(common.SuccessResponse{Success: true})
}
// func (co RBACController) Init(c *fiber.Ctx) error {
// 	tx := dba.DB.Begin()

//     // Create admin role
//     adminRole := dts.Role{
//         Name: "system-admin",
//         Description: "System Administrator with full access",
//     }
//     if err := tx.Create(&adminRole).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     // Create admin group
//     adminGroup := dts.Group{
//         Name: "full-access",
//         Description: "Full system access group",
//     }
//     if err := tx.Create(&adminGroup).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     // Assign admin group to admin role
//     if err := tx.Model(&adminRole).Association("Groups").Append(&adminGroup); err != nil {
//         tx.Rollback()
//         return err
//     }

//     adminUser := dts.User{
//         Email: "admin@example.com",
//         Password: "hashed_password", // Make sure to hash this
//         UserType: "system-user",
//         RoleID: adminRole.ID,
//         IsActive: true,
//     }
//     if err := tx.Create(&adminUser).Error; err != nil {
//         tx.Rollback()
//         return err
//     }

//     return tx.Commit().Error
// }