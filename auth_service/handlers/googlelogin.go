package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	entity "ubm-canteen/models"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/api/idtoken"
)

type GoogleHandlers struct {
	*BaseHandler
}

func NewGoogleHandlers(base *BaseHandler) *GoogleHandlers {
	return &GoogleHandlers{BaseHandler: base}
}

func (g *GoogleHandlers) GoogleLogin(c *fiber.Ctx) error {
	var in entity.GoogleLogin
	if err:= c.BodyParser(&in); err!=nil{
		return c.Status(400).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}
	
	if in.IdToken == "" || (in.Role != "seller" && in.Role!= "user"){
		return c.Status(400).JSON(fiber.Map{
			"error": "id_token & role are required",
		})
	}

	pl, err:= idtoken.Validate(context.Background(), in.IdToken, os.Getenv("WEB_CLIENT_ID"))
	if err != nil{
		return c.Status(401).JSON(fiber.Map{
			"error": "invalid google id token",
		})
	}

	email   := pl.Claims["email"].(string)
	name    := pl.Claims["name"].(string)
	googleU := pl.Subject

	// DB

	ctx:= context.Background()
	conn, err:= g.DB.Acquire(ctx)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "fail to connect the database",
		})
	}
	var id uuid.UUID
	switch strings.ToLower(in.Role) {
	case "seller":
		err = conn.QueryRow(ctx, `SELECT id_seller FROM seller WHERE email=$1`, email).Scan(&id)
		if err == pgx.ErrNoRows {
			err = conn.QueryRow(ctx,
				`INSERT INTO seller (google_uid, email, nama_seller)
				 VALUES ($1,$2,$3) RETURNING id_seller`,
				googleU, email, name).Scan(&id)
		}
	case "user":
		err = conn.QueryRow(ctx, `SELECT id_user FROM users WHERE email=$1`, email).Scan(&id)
		if err == pgx.ErrNoRows {
			err = conn.QueryRow(ctx,
				`INSERT INTO users (google_uid, email, fullname)
				 VALUES ($1,$2,$3) RETURNING id_user`,
				googleU, email, name).Scan(&id)
		}
	}
	if err != nil { return serverErr(c, err, "upsert") }

	// ---------- JWT ----------
	acc, ref, err := makeTokens(id, email, in.Role, c.IP())
	if err != nil { return serverErr(c, err, "jwt") }

	// ---------- Redis ----------
	if err := utils.RedisClient.Set(ctx,
		fmt.Sprintf("%s:token:%s", in.Role, id), acc, 24*time.Hour).Err(); err != nil {
		return serverErr(c, err, "redis token")
	}
	if err := utils.RedisClient.Set(ctx,
		fmt.Sprintf("%s:refresh:%s", in.Role, id), ref, 30*24*time.Hour).Err(); err != nil {
		return serverErr(c, err, "redis refresh")
	}

	// ---------- Done ----------
	return c.JSON(fiber.Map{
		"status":        "success",
		"email":         email,
		"role":          in.Role,
		"token":         acc,
		"refresh_token": ref,
	})

	
}


func makeTokens(id uuid.UUID, email, role, ip string) (string, string, error) {
	claims := jwt.MapClaims{
		"id": id.String(), "email": email, "role": role,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	if err != nil { return "", "", err }

	rClaims := jwt.MapClaims{
		"id": id.String(), "email": email, "role": role, "ip": ip,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims).SignedString(jwtSecret)
	return at, rt, err
}

func serverErr(c *fiber.Ctx, e error, tag string) error {
	return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("%s: %v", tag, e)})
}
