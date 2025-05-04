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
	input := new(entity.GoogleLogin)	
	ctx:= c.Context()

	//PARSE THE REQUEST BODY
	if err:= c.BodyParser(input); err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"cannot parse JSON",
		})
	}

	//VALIDATE INPUT
	if input.IdToken == "" || input.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"id token and valid role are required",
		})
	}
	
	//VALIDATE GOOGLE ID TOKEN

	audience := os.Getenv("WEB_CLIENT_ID")

	payload, err := idtoken.Validate(context.Background(), input.IdToken, audience)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": fmt.Sprintf("invalid google id token:%v", err),
		})
	}

	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	//acquire db connection
	conn, err:= g.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":fmt.Sprintf("failed to acquire db connection: %v", err),
		})
	}

	defer conn.Release()

	var roleLower = strings.ToLower(input.Role)

	switch roleLower {
	case "seller":
		squery := `SELECT id_seller FROM seller WHERE email=$1`
		db_seller:= new(entity.Seller)
		err = conn.QueryRow(ctx, squery, email).Scan(&db_seller.IdSeller, &db_seller.Email, &db_seller.Password)
		if err != nil {
			if err == pgx.ErrNoRows{
				//INSERT NEW SELLER
				iquery := `INSERT INTO seller (nama_seller, email, password, phone_num) VALUES ($1, $2, $3, $4) RETURNING id_seller`
				err = conn.QueryRow(ctx, iquery, name, email, "-").Scan(&db_seller.IdSeller)
				if err != nil {
					if strings.Contains(err.Error(), "duplicate key"){
						return c.Status(fiber.StatusConflict).JSON(fiber.Map{
							"status":"error",
							"message": "email already exist"})
					}
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "failed to create seller account",
					})
				}
		} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to query seller",
			},
		)
	}
}

case "user":
		query := `SELECT id_users FROM users WHERE email=$1`
		dbuser:= new(entity.User)
		err = conn.QueryRow(ctx, query, email).Scan(&dbuser.IdUsers, &dbuser.Email)
		if err != nil {
			if err == pgx.ErrNoRows{
				//insert user baru
				iquery := `INSERT INTO users (nama_users, email, password) VALUES ($1, $2, $3) RETURNING id_users`
				err = conn.QueryRow(ctx, iquery, name, email).Scan(&dbuser.IdUsers)
				if err != nil {
					if strings.Contains(err.Error(), "duplicate key") {
						return c.Status(fiber.StatusConflict).JSON(fiber.Map{
							"status":"error",
							"message": "email already exist",
						})
					}
				}
		}
	}

	//JWT TOKEN
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":dbuser.IdUsers.String(),
			"email": email,
			"role": roleLower,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)	
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":"failed to generate token",	
			})
	}

	//refresh token
	refreshToken:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id": dbuser.IdUsers.String(),
			"email":dbuser.Email,
			"role": roleLower,
			"ip": c.IP(),
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	)

	refreshTokenStr, err := refreshToken.SignedString(jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":"failed to generate refresh token",
			})
	}

	//save token and refresh token to redis
	tokenKey := fmt.Sprintf("%s:token:%s", input.Role, dbuser.IdUsers)
	refreshTokenKey := fmt.Sprintf("%s:refresh_token:%s", input.Role, dbuser.IdUsers)
	
	err = utils.RedisClient.Set(ctx, tokenKey, tokenString, 24*time.Hour).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to store token in redis",
		})
	}

	err = utils.RedisClient.Set(ctx, refreshTokenKey, refreshTokenStr, time.Hour*24*30).Err()
	if err!=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"failed to store token in redis",
		})	
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Success",
		"message": "Successfully logged in with Google",
		"data": fiber.Map{
			"email": email,
			"role": input.Role,
			"token": tokenString,
			"refresh_token": refreshTokenStr,
			}},
		)
	}

	// Optionally handle other roles or return an error fosw unsupported roles
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "unsupported role",
	})
}