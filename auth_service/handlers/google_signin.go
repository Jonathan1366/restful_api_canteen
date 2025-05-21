package handlers

import (
	entity "ubm-canteen/models"
	"ubm-canteen/repository"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
)



func (h *BaseHandler) GoogleSignIn(c *fiber.Ctx) error {
	payload:= new(entity.GoogleLogin)

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload"})
	}

	googleUser, err:= utils.VerifyGoogleIDToken(payload.IdToken)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Google ID token"})
	}

	conn, err:= h.DB.Acquire(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to acquire db connection",
		})
	}
	defer conn.Release()

	userId, err:= repository.FindOrCreateGoogleUser(c.Context(), conn.Conn(), googleUser, payload.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find or create user"})
	}

	accessToken, _:= utils.GenerateJWTSecret(userId, payload.Role)
	// Generate refresh token
	refreshToken, _:= utils.GenerateRefreshToken(userId)
	return c.JSON(fiber.Map{
		"message": "Login sukses",
		"data": fiber.Map{
			"access_token": accessToken,
			"refresh_token": refreshToken,
		},
	})
	}
