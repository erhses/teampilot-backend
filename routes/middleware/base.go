package middleware

import (
	"errors"
	"log"

	"regexp"
	"strings"

	"teampilot/constants"
	"teampilot/integrations/dba"
	"teampilot/structs/dts"
	"teampilot/utils"

	"github.com/gofiber/fiber/v2"
)

const authKey = "key"

// const adminAuthKey = "admin_secret"

var publicUrls = []string{
	"admin/settings/rbac/group",
	"admin/auth/info",
}

func AuthMiddleware(c *fiber.Ctx) error {
	headers := c.GetReqHeaders()
	requiredToken, ok := headers["Authorization"]
	if !ok || len(requiredToken) == 0 || len(requiredToken[0]) < 8 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "та нэвтрэх хэрэгтэй",
			"body":    nil,
		})
	}

	tokenManager := utils.GetTokenManager()
	claims, err := tokenManager.ValidateToken(c.Context(), requiredToken[0][7:])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Таны нэвтрэх хугацаа дууссан байна",
			"body":    nil,
		})
	}

	db := dba.DB
	userID, ok := claims["id"].(float64)
	if ok {
		// if strings.HasPrefix(c.Path(), "/api/v1/admin") && claims["user_type"] != constants.USER_SYSTEM {
		// 	// log.Println("reaching here!!")
		// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		// 		"message": "Админ эрх шаардлагатай",
		// 		"body":    nil,
		// 	})
		// }
		switch claims["user_type"].(string) {
		case "owner":
			// agent
			var user *dts.User
			if err := db.Where("user_type = ?", constants.USER_TEACHER).First(&user, uint(userID)).Error; err != nil {
				return errors.New("хэрэглэгч олдсонгүй")
			}
			if !user.IsActive {
				return errors.New("та идэвхгүй хэрэглэгч байна")
			}
			c.Locals(authKey, user)
			return c.Next()
		case "tenant":
			var user *dts.User
			if err := db.Where("user_type = ?", constants.USER_STUDENT).First(&user, uint(userID)).Error; err != nil {
				return errors.New("хэрэглэгч олдсонгүй")
			}
			if !user.IsActive {
				return errors.New("та идэвхгүй хэрэглэгч байна")
			}
			c.Locals(authKey, user)
			return c.Next()
		case "system-user":
			// user
			var user *dts.User
			if err := db.Where("user_type = ?", constants.USER_SYSTEM).First(&user, uint(userID)).Error; err != nil {
				return errors.New("хэрэглэгч олдсонгүй")
			}
			if !user.IsActive {
				return errors.New("та идэвхгүй хэрэглэгч байна")
			}
			c.Locals(authKey, user)
			return c.Next()

		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Хандах эрхгүй байна.",
				"body":    nil,
			})
		}

	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"message": "Хандах эрхгүй байна.",
		"body":    nil,
	})
}

func RetrieveAuth(c *fiber.Ctx) *dts.User {
	user, ok := c.Locals(authKey).(*dts.User)
	if !ok {
		log.Println("RetrieveAuth: No user found in context")
		return nil
	}
	return user
}

// func GetCurrentPermission(c *fiber.Ctx) *databases.Permission {
// 	return c.Locals("permission").(*databases.Permission)
// }

func removeIdSegment(input string) string {
	// Check if ":id" is present at the end of the string
	if strings.HasSuffix(input, ":id") {
		// If yes, remove ":id"
		return input[:len(input)-3]
	}
	// If ":id" is not found, return the original string
	return input
}

func replaceID(path string) string {
	re := regexp.MustCompile(`/\d+`)

	result := re.ReplaceAllString(path, "/:id")
	return result
}

func matchPath(storedPath, actualPath string) bool {
	pattern := "^" + strings.ReplaceAll(storedPath, "{id}", "([^/]+)") + "$"
	re := regexp.MustCompile(pattern)
	return re.MatchString(actualPath)
}

func RBACMiddleware(c *fiber.Ctx) error {
	user := RetrieveAuth(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var role dts.Role
	if err := dba.DB.Preload("Groups.Permissions").First(&role, user.RoleID).Error; err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	path := strings.TrimPrefix(c.Path(), "/api/v1")
	method := strings.ToLower(c.Method())

	for _, group := range role.Groups {
		for _, permission := range group.Permissions {
			if matchPath(permission.Path, path) && permission.Method == method {
				return c.Next()
			}
		}
	}
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"message": "Access denied",
	})
}
