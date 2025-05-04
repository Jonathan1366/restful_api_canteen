package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func generateJWTSecret(userId, email string) (string, error){
	claims:= jwt.MapClaims{
		"sub": userId,
		"email": email,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //valid for 30 days
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}