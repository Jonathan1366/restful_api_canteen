package utils

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//validate mime type

func IsValidMimeType(mimeType string, allowedTypes []string) bool {
	for _, allowed := range allowedTypes {
		if mimeType == allowed {
			return true
		}
	}
	return false
}

//check seller in db

func IsSellerExist(ctx context.Context, db *pgxpool.Pool, sellerID string) bool {
	var exist bool
	query := "SELECT EXISTS (SELECT 1 FROM SELLER WHERE id_seller=$1)"
	err := db.QueryRow(ctx, query, sellerID).Scan( &exist)
	return err == nil && exist
}

//update database

func UpdateSellerContract(ctx context.Context, db *pgxpool.Pool, sellerID, fileName, mimeType string) error  {
	query := `UPDATE SELLER SET file_contract=$1, mime_type=$2, is_valid_contract=FALSE WHERE id_seller= $3`
	_, err := db.Exec(ctx, query, fileName, mimeType, sellerID)
	return err
}

func ProcessTextractAndValidate(ctx context.Context, db *pgxpool.Pool, bucket, key, sellerID string)  {
	keywords := []string{"nama", "nik", "tanggal"}
	foundKeywords, err:= ProcessDocumentWithTextractFromS3(
		ctx,
		bucket,
		key,
		keywords, 
	)

	isValid := err == nil && len(foundKeywords) > 0
	query := "UPDATE SELLER SET is_valid_contract=$1 WHERE id_seller=$2"
	_, err = db.Exec(ctx, query, isValid, sellerID)
	if err != nil {
		fmt.Printf("Failed to update validation status: %v\n",err)
	}
}