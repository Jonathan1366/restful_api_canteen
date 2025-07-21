// file: middleware/auth_middleware.go
package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// === LANGKAH 1: AMBIL TOKEN ===
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Request needs a token"})
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token format"})
		}

		// === LANGKAH 2: VALIDASI TOKEN & AMBIL ID ===
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid token claims"})
		}
		claimKey := "id_" + role
		id, idOk := claims[claimKey].(string)
		if !idOk {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Access forbidden for this role"})
		}

		// === LANGKAH 3: TERUSKAN REQUEST ===
		c.Locals(claimKey, id)
		return c.Next()
	}
}