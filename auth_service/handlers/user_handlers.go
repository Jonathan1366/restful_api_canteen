package handlers

import (
	"fmt"
	"strings"
	"time"
	entity "ubm-canteen/models"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	*BaseHandler
}

func NewUserHandlers(base *BaseHandler) *UserHandler {
	return &UserHandler{base}
}

// REGISTER USER
func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	input := new(entity.User) // General registration struct
	ctx := c.Context()

	err := utils.DeallocateAllStatement(ctx, h.DB)
	if err != nil && err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to deallocate statement",
		})
	}

	// Parse and validate the request body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    400,
			"message": "Invalid input: Failed to parse request body",
		})
	}

	if input.Email == "" || input.Password == "" || input.NamaUsers == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Input: Email, Name, and Password are required.",
		})
	}

	// HASH PASSWORD
	hashedPass, err := utils.HashPass(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Failed to hash password",
		})
	}

	// Create user entity
	user := entity.User{
		IdUsers:   uuid.New(),
		NamaUsers: input.NamaUsers,
		Email:     input.Email,
		Password:  hashedPass,
	}

	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to acquire connection: %v", err),
		})
	}
	defer conn.Release()

	// Insert user into the database
	query := `INSERT INTO "user" (id_users, nama_users, email, password) VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(ctx, query, user.IdUsers, user.NamaUsers, user.Email, hashedPass)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to register user: %v", err),
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully registered",
		"data": fiber.Map{
			"id":    user.IdUsers.String(),
			"email": user.Email,
			"name":  user.NamaUsers,
		},
	})
}

// login user
func (h *UserHandler) LoginUser(c *fiber.Ctx) error {
	input := new(entity.User) // Use a general login struct
	ctx := c.Context()

	// Parse and validate the request body
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid input",
		})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Email and password are required.",
		})
	}
	
	// Query the database for user data based on the email
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to acquire connection: %v", err),
		})
	}

	defer conn.Release()

	dbUser := new(entity.User)
	query := `SELECT id_users, email, password FROM "user" WHERE email=$1`
	err = conn.QueryRow(ctx, query, input.Email).Scan(&dbUser.IdUsers, &dbUser.Email, &dbUser.Password)

	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Login failed",
		})
	}

	// Verify the password
	if !utils.CheckPassHash(input.Password, dbUser.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}
	
	//JWT TOKEN
	tokenString, err:= utils.GenerateToken(dbUser.IdUsers.String(), dbUser.Email, "user")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		})
	}

	//REFRESH TOKEN
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_users":  dbUser.IdUsers.String(),
		"email":      dbUser.Email,
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(), //valid for 30 days
	})
	
	refreshTokenString, err := refreshToken.SignedString(h.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}
	//SAVE TOKEN TO REDIS
	err = utils.RedisClient.Set(ctx, "token:"+dbUser.IdUsers.String(), tokenString, time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to store token in Redis",
		})
	}

	//SAVE REFRESH TOKEN TO REDIS
	err = utils.RedisClient.Set(ctx, "refresh token:"+dbUser.IdUsers.String(), refreshTokenString, time.Hour*24*30).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to store refresh token in Redis",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"email":         dbUser.Email,
			"token":         tokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

func (h *UserHandler) LogoutUser(c *fiber.Ctx) error {
	//invalid jwt token (for example, by storing it in a blacklist)
	token := c.Get("Authorization")

	if token == "" {
				token = c.Get("accessToken")
	}
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "no token provided",
		})
		// remove bearer from token string if present
	}
	
	token = strings.TrimPrefix(token, "Bearer ")

	// Parse token dulu untuk validasi sebelum delete dari Redis
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
		}
		return h.JWTSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token or error in parsing token",
		})
	}

	// Validate claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token",
		})
	} 
	// Extract user ID
	idUserStr, ok := claims["id_users"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user ID in token",
		})
	}

	idUser, err := uuid.Parse(idUserStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid user ID format",
		})
	}

	// DELETE TOKEN FROM REDIS (perbaikan: hilangkan space)
	ctx := c.Context()
	err = utils.RedisClient.Del(ctx, "token:"+idUser.String()).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete token from Redis",
		})
	}

		// DELETE REFRESH TOKEN FROM REDIS 
	err = utils.RedisClient.Del(ctx, "refresh token:"+idUser.String()).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete refresh token from Redis",
		})
	}

	// Input token to revocation list (perbaikan: tambah entityType parameter)
	err = h.TokenRevocationLogic(idUser, "user", token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to logout and revoke token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "logged out successfully",
	})
}