package utils

import (
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// const AbleJWTSecret = "UA1JJ8OxpVN7ZjjXk1R9"

// access secret key
func GetAccessSecretKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}


// GenerateToken ...
func GenerateToken(userId uint, userType string) string {
	accessExpTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"id":        userId,
		"user_type": userType,
		"exp":       accessExpTime.Unix(),
	}

	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(GetAccessSecretKey())
	return accessToken
}

// GenerateHash password hash generate
func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compare password and hash
func ComparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RandomWithCharset(length int, charset []rune) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomUpperCase(l int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return RandomWithCharset(l, letters)
}

func RandomString(l int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	return RandomWithCharset(l, letters)
}

func RandomNumber(l int) string {
	var letters = []rune("0123456789")
	return RandomWithCharset(l, letters)
}
