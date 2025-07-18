package handlers

import (
	"ubm-canteen/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)
type BaseHandler struct {
	DB *pgxpool.Pool
	RedisClient *redis.Client
	Presigner utils.Presigner
	DefaultQueryExecMode pgx.QueryExecMode
	JWTSecret []byte // Secret key for JWT signing
}





	