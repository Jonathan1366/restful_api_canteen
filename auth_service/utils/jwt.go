package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateToken(entityID, email, entityType string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"sub":   entityID, // 'sub' (subject) adalah klaim standar untuk ID entitas
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token berlaku selama 24 jam
	}

	// Tentukan nama klaim ID berdasarkan entityType
	switch entityType {
	case "user":
		claims["id_users"] = entityID
	case "seller":
		claims["id_seller"] = entityID
	default:
		// Jika tipe entitas tidak valid, kembalikan error
		return "", fmt.Errorf("tipe entitas tidak valid: %s", entityType)
	}

	// Buat token dengan klaim yang sudah ditentukan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token dengan secret key Anda
	return token.SignedString(jwtSecret)
}