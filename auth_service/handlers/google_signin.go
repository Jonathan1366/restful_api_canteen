package handlers

import (
	"log"
	"os"
	"time"
	entity "ubm-canteen/models"
	"ubm-canteen/repository"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

type GoogleHandler struct {
	*BaseHandler
}

func NewGoogleHandlers(base *BaseHandler) *GoogleHandler{
	return &GoogleHandler{base}
}

func (h *GoogleHandler) GoogleSignIn(c *fiber.Ctx) error {
	

	//PARSE & VALIDATE PAYLOAD	
	payload := new(entity.GoogleLogin)
	if err := c.BodyParser(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if payload.IdToken == "" || payload.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing ID token or role")
	}

	//VERIFY ID TOKEN WITH GOOGLE LIBRARY
	aud:=os.Getenv("WEB_CLIENT_ID")
	token, err := idtoken.Validate(c.Context(), payload.IdToken, aud)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid ID token: "+err.Error())
		
	} 

	subject := token.Subject
	claims := token.Claims
	email, _ := claims["email"].(string)
	name, _ := claims["name"].(string)

	googleuser := &entity.GoogleUser{
		Sub: subject,
		Email: email,
		Name: name,
	}

	conn, err := h.DB.Acquire(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to acquire DB connection")
	}
	
	defer conn.Release()

	userID, err := repository.FindOrCreateGoogleUser(c.Context(), conn.Conn(), googleuser, payload.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to find or create user: "+err.Error())
	}

	accessToken, err := utils.GenerateJWTSecret(userID, payload.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	//REFRESH TOKEN
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_token": payload.IdToken,
		"role":      payload.Role,
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate refresh token")	
	}

	// SAVE TO REDIS
	ctx := c.Context()

	if err := utils.RedisClient.Set(ctx, "token:"+payload.IdToken, accessToken, 24*time.Hour).Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to store access token in redis")
	}

	// Store refresh token in Redis
	if err := utils.RedisClient.Set(ctx, "refresh:"+payload.IdToken, refreshTokenString, time.Hour*24*30).Err(); err!=nil{
		return fiber.NewError(fiber.StatusInternalServerError, "failed to store refresh token in redis")
	}

	log.Println("Received UID:", subject) // Potong biar gak kepanjangan
	log.Println("Role:", payload.Role)
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login sukses",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshTokenString,
		},
	})
}
