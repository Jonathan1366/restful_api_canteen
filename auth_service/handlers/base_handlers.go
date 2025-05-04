package handlers

import (
	"ubm-canteen/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BaseHandler struct {
	DB *pgxpool.Pool
	Presigner utils.Presigner
	DefaultQueryExecMode pgx.QueryExecMode
}





	