package repository

import (
	"context"
	entity "ubm-canteen/models"
)

type SellerRepository interface {
	Create(ctx context.Context, seller *entity.Seller) error
	Update(ctx context.Context, seller *entity.Seller) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entity.Seller, error)
	GetByEmail(ctx context.Context, email string) (*entity.Seller, error)
}

