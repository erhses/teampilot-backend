package middleware

import (
	// "erhses/shorten-url-fiber-redis/routes/controller/common"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var SecretKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		log.Fatalf("SECRET_KEY not set in environment")
	}
	SecretKey = []byte(secret)
}

func GenerateJWT(userID uint, role string) (string, error) {

	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

func Authenticate(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
		// return common.Controller.SetError(c, http.StatusUnauthorized, "Missing or invalid token")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		// log.Println(err)
		// log.Println(token.Valid)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	// log.Println(tokenString)

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("userID", claims["userID"])
	c.Locals("role", claims["role"])
	// log.Println(c.Locals("userID"))
	// log.Println(c.Locals("role"))

	return c.Next()
}

func Authorize(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access forbidden"})
		}
		return c.Next()
	}
}
