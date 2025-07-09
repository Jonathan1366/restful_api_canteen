package middleware

import (
	"context"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
)

type AuthHandlers struct {
	DB *pgxpool.Pool
	DefaultQueryExecMode pgx.QueryExecMode
}

func (h*AuthHandlers) TokenValidationMiddleware(c*fiber.Ctx) error  {
	token:= c.Get("Authorization")
	if token=="" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":"error",
			"message":"no token provided",
		})
	}
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
	
	//check is there any token in revocation list
	var exists bool
		err:=h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE token = $1)", token).Scan(&exists)
		if err!=nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":"error",
				"message":"database error",
			})
		}
		if exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":"error",
				"message":"token has been revoked",
			})
		}
		return c.Next()
}