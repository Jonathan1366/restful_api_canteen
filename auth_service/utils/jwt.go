package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWTSecret(userId, email string) (string, error){
	claims:= jwt.MapClaims{
		"id_seller": userId,
		"sub": userId,
		"email": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //valid for 30 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}