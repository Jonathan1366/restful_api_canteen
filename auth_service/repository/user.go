package repository

import (
	"context"
	"errors"
	entity "ubm-canteen/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)



func FindOrCreateGoogleUser(ctx context.Context, db *pgx.Conn, user *entity.GoogleUser, role string) (string, error) {
	
	var id string

	switch role {
	case "seller":
		// Cek apakah user seller sudah ada berdasarkan google_uid
		query := `SELECT id_seller FROM seller WHERE google_uid = $1`
		err := db.QueryRow(ctx, query, user.Sub).Scan(&id)
		if err == pgx.ErrNoRows {
			// Cek kalau email sudah pernah dipakai di tabel users
			emailCheck := `SELECT 1 FROM users WHERE email = $1 LIMIT 1`
			var dummy int
			if err := db.QueryRow(ctx, emailCheck, user.Email).Scan(&dummy); err == nil {
				return "", errors.New("email sudah digunakan oleh user biasa")
			}

			insertQuery := `INSERT INTO seller (id_seller, nama_seller, email, google_uid, profile_pic, created_at)
				VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id_seller`
			err = db.QueryRow(ctx, insertQuery, uuid.New(), user.Name, user.Email, user.Sub, user.Picture).Scan(&id)
			if err != nil {
				return "", err
			}
		} else if err != nil {
			return "", err
		}

	case "user":
		// Cek apakah user sudah ada
		query := `SELECT id_users FROM users WHERE email = $1`
		err := db.QueryRow(ctx, query, user.Email).Scan(&id)
		if err == pgx.ErrNoRows {
			// Pastikan tidak bentrok dengan seller
			sellerCheck := `SELECT 1 FROM seller WHERE google_uid = $1 LIMIT 1`
			var dummy int
			if err := db.QueryRow(ctx, sellerCheck, user.Sub).Scan(&dummy); err == nil {
				return "", errors.New("akun ini sudah digunakan sebagai seller")
			}

			insertQuery := `
				INSERT INTO users (id_users, nama_users, email, img_user, time_stamp)
				VALUES ($1, $2, $3, $4, NOW())
				RETURNING id_users`
			userID := uuid.New().String()
			err = db.QueryRow(ctx, insertQuery, userID, user.Name, user.Email, user.Picture).Scan(&id)
			if err != nil {
				return "", err
			}
		} else if err != nil {
			return "", err
		}

	default:
		return "", errors.New("invalid role")
	}

	return id, nil
}
