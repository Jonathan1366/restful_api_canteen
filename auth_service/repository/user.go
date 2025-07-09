package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
  DB *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
  return &UserRepo{DB: db}
}

func (r *UserRepo) FindOrCreateGoogleUser(ctx context.Context, sub, email, name, role string) (string, error) {
  var table, idCol, nameCol string
  if role == "seller" {
    table, idCol, nameCol = "seller", "id_seller", "nama_seller"
  } else {
    table, idCol, nameCol = "users", "id_users", "nama_users"
  }
  
  // 1) Cek existing
  var id string
  err := r.DB.QueryRow(ctx,
    `SELECT `+idCol+` FROM `+table+` WHERE google_uid = $1`,
    sub,
  ).Scan(&id)
  if err != nil && !errors.Is(err, pgx.ErrNoRows) {
    return "", err
  }

  // 2) Jika belum ada, insert baru
  if errors.Is(err, pgx.ErrNoRows) {
    id = uuid.NewString()
    _, err = r.DB.Exec(ctx,
      `INSERT INTO `+table+`
         (`+idCol+`, `+nameCol+`, email, google_uid)
       VALUES ($1,$2,$3,$4)`,
      id, name, email, sub,
    )
    if err != nil {
      return "", err
    }
  }

  return id, nil
}
