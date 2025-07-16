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
	input := new(entity.RegisterUser) // General registration struct
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
	input := new(entity.LoginUser) // Use a general login struct
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
	query := `SELECT id_users, email, password FROM user WHERE email=$1`
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
	if input.Password != dbUser.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid password",
		})
	}

	// JWT Token generation
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_users": dbUser.IdUsers.String(),
		"email":    dbUser.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	//Sign the token with secret key upload to frontend
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("JWT Token generation error: %v\n", err) // Log error jika gagal
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate jwt token",
		})
	}

	//Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_seller":  dbUser.IdUsers.String(),
		"email":      dbUser.Email,
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(), //valid for 30 days
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Login successful",
		"data": fiber.Map{
			"email":         dbUser.Email,
			"token":         accessTokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

func (h *UserHandler) LogoutUser(c *fiber.Ctx) error {
	//invalid jwt token (for example, by storing it in a blacklist)
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "no token provided",
		})
	}
	// remove bearer from token string if present
	token = strings.TrimPrefix(token, "Bearer ")

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token or error in parsing token",
		})
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		IdUsers := claims["id_users"].(string)
		//input token to revocation list
		err := h.TokenRevocationLogic(uuid.MustParse(IdUsers), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to logout and revoke token",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "Success",
			"message": "logged out successfully",
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token",
		})
	}
}

func (h *UserHandler) TokenRevocationLogic(d uuid.UUID, token string) any {
	panic("unimplemented")
}
