package handlers

import (
	entity "ubm-canteen/models"
	"ubm-canteen/repository"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
)

type GoogleHandler struct {
	*BaseHandler
	Repo *repository.UserRepo
}

func NewGoogleHandlers(base *BaseHandler) *GoogleHandler{
	return &GoogleHandler{BaseHandler: base, Repo: repository.NewUserRepo(base.DB)}
}

func (h *GoogleHandler) GoogleSignIn(c *fiber.Ctx) error {	
	//PARSE & VALIDATE PAYLOAD	
	req := new(entity.GoogleLoginReq)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request payload")
	}

	if req.IDToken == "" || req.Code == "" || req.Role == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing ID token or role")
	}

	ctx := c.Context()
	payload, err := utils.VerifyGoogleIDToken(c.Context(), req.IDToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid ID token: "+err.Error())
	}

	userID, err := h.Repo.FindOrCreateGoogleUser(
		ctx, payload.Subject, 
		payload.Claims["email"].(string),
		payload.Claims["name"].(string),
		req.Role,
	)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "db error: "+err.Error())
	}
	
	tok, err := utils.ExchangeAuthCode(ctx, req.Code)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to exchange auth code: "+err.Error())
	}

	if err := utils.StoreRefreshToken(ctx, h.RedisClient, userID, tok.RefreshToken); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to store refresh token: "+err.Error())
	}

	  // 6) Generate internal JWT (access token aplikasi)
  jwtStr, err := utils.GenerateToken(userID, payload.Claims["email"].(string), req.Role)
  if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate JWT: "+err.Error())
  }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "message": "Login sukses",
    "data": fiber.Map{
      "access_token":  jwtStr,
      "refresh_token": tok.RefreshToken,
    },
  })
}
