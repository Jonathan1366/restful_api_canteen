package repository

import (
	"context"
	"errors"
	entity "ubm-canteen/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func FindOrCreateGoogleUser(ctx context.Context, db *pgx.Conn, acc *entity.GoogleUser, role string) (string, error) {
	switch role {
		case "seller":
			var id string
			err := db.QueryRow(ctx, `SELECT id_seller from seller where google_uid = $1`, acc.Sub).Scan(&id)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows){
					id = uuid.NewString()
					_, err = db.Exec(ctx, `INSERT INTO seller (id_seller, nama_seller, email, google_uid) VALUES ($1, $2, $3, $4)`, id, acc.Name, acc.Email, acc.Sub)
					if err!=nil{
						return "", err
					}
				} else{
					return "", err
				}
			}
			return id, nil

		case "user":
			var id string
			err := db.QueryRow(ctx, `SELECT id_users from users where google_uid = $1`, acc.Sub).Scan(&id)
			if err!=nil{
				if errors.Is(err, pgx.ErrNoRows){
					id = uuid.NewString()
					_, err = db.Exec(ctx, `INSERT INTO users (id_users, nama_users, email, google_uid) VALUES ($1, $2, $3, $4)`, id, acc.Name, acc.Email, acc.Sub)
					if err!=nil{
						return "", err
					}
			} else {
				return "", err
			}
	}
	return id, nil
		default:
			return "", errors.New("invalid role")
		}
}
