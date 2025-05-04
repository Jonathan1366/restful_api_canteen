package handlers

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5"
)

type ProductService struct {
	DB                   *pgxpool.Pool
	DefaultQueryExecMode pgx.QueryExecMode
}

