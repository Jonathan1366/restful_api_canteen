package test

// import (
// 	"context"
// 	"testing"
// 	"ubm-canteen/utils"
// )

// func TestBucketExist(t *testing.T)  {
// 	ctx:= context.Background()
// 	bucketName := "seller-textract-bucket"
// 	basics := utils.S3Bucket{}

// 	exist, err:= basics.BucketExists(ctx, bucketName)
// 	if err != nil {
// 		t.Fatalf("Failed to check bucket existance: %v",err)
// 	}
// 	if !exist {
// 		t.Errorf("Bucket %s should exist", bucketName)
// 	}
// }

// func TestUploadFile(t *testing.T)  {
// 	ctx:=context.Background()
// 	bucketName := "seller-textract-bucket"
// 	objectkey := "file"
// 	fileName := ""
// 	basics := utils.S3Bucket{}

// 	err := basics.UploadFile(ctx, bucketName, objectkey,)
// 	if err != nil {
// 		t.Fatalf("Failed to upload file: %v", err)
// 	}
// }