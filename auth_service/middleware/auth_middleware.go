package middleware

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var jwksUrl = []byte(os.Getenv("JWKS_URL"))
var keySetCache *jwk.Cache

func init(){
	keySetCache = jwk.NewCache(context.Background())
	keySetCache.Register(string(jwksUrl), jwk.WithRefreshInterval(1*time.Hour))
}

func Supabase () fiber.Handler{
	return func (c*fiber.Ctx) error {
		authHeaders := c.Get("Authorization")
		if authHeaders == ""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}
		tokenString := strings.TrimPrefix(authHeaders, "Bearer")
		if tokenString==""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}
		keySet, err := jwk.Fetch(c.Context(), string(jwksUrl))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":"Failed to fetch authorization keys",
			})
		}
		token, err:= jwt.ParseString(tokenString, jwt.WithKeySet(keySet))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		supabaseUserID := token.Subject()
		if supabaseUserID == ""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User ID not found in token",
			})
		}

		c.Locals("id_seller", supabaseUserID)
		return c.Next()
	}
}
