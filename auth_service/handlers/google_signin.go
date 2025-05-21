package handlers

import (
	"os"
	"time"
	entity "ubm-canteen/models"
	"ubm-canteen/repository"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)


type GoogleHandler struct {
	*BaseHandler
}

func NewGoogleHandlers(base *BaseHandler) *GoogleHandler{
	return &GoogleHandler{base}
}

func (h *GoogleHandler) GoogleSignIn(c *fiber.Ctx) error {
	
	ctx := c.Context()

	payload := new(entity.GoogleLogin)


	if err := c.BodyParser(payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate the payload
	if payload.IdToken == "" || payload.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing ID token or role")
	}

	
	aud:=os.Getenv("WEB_CLIENT_ID")
	if aud == ""{
		return fiber.NewError(fiber.StatusInternalServerError, "WEB_CLIENT_ID not set in env")
	}

	googleUser, err := utils.VerifyGoogleIDToken(ctx, payload.IdToken, aud)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Google ID token: "+err.Error())
	}

	conn, err := h.DB.Acquire(c.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to acquire DB connection")
	}
	defer conn.Release()

	userID, err := repository.FindOrCreateGoogleUser(c.Context(), conn.Conn(), googleUser, payload.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to find or create user: "+err.Error())
	}

	accessToken, err := utils.GenerateJWTSecret(userID, payload.Role)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate access token")
	}

//REFRESH TOKEN
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   userID,
		"role":      payload.Role,
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sign refresh token",
		})
	}

	// SAVE TO REDIS

	err = utils.RedisClient.Set(ctx, "token:"+userID, accessToken, time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store token in Redis",
		})
	}
	
	// Store refresh token in Redis

	err = utils.RedisClient.Set(ctx, "refresh:"+userID, refreshTokenString, time.Hour*24*30).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store refresh token in Redis",
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login sukses",
		"data": fiber.Map{
			"access_token":  accessToken,
			"refresh_token": refreshTokenString,
		},
	})
}
