package utils

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DeallocateStatement(ctx context.Context, db *pgxpool.Pool, statementName string)error {
	conn, err:= db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.Conn().Deallocate(ctx, statementName)
}

func DeallocateAllStatement(ctx context.Context, db *pgxpool.Pool) error  {
	conn, err:=db.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	return conn.Conn().DeallocateAll(ctx)
}