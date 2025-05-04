package handlers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (h *BaseHandler) TokenRevocationLogic(IdSeller uuid.UUID, token string) error {
	ctx := context.Background()
	conn, err:=h.DB.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	query:= `INSERT INTO revoked_tokens (id_seller, token, revoked_at) values ($1, $2, $3)`
	_,err = conn.Exec(ctx, query, IdSeller, token, time.Now())
	if err != nil {
		return err
	}
	return nil
}