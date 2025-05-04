package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Presigner struct{
	s3 *s3.Client
	PresignClient *s3.PresignClient
}

func NewS3Client() (*Presigner, error){
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if region == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing required AWS credentials in environment variables")
	}

	creds:= credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),	
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {      
		return nil, fmt.Errorf("failed to load configuration, %v", err)
	}
	
	//build s3 client
	client:= s3.NewFromConfig(cfg)
	//build presign client
	presignClient:= s3.NewPresignClient(client)
	return &Presigner{
		s3: client,
		PresignClient: presignClient,
	}, nil
}

func (presigner *Presigner) GetObject(
	ctx context.Context,
	bucketName string,
	objectKey string,
	lifetimeSecs int64,
) (*v4.PresignedHTTPRequest, error){
	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:aws.String(bucketName),
		Key: aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why:%v\n", bucketName, objectKey, err)
	}
	return request, err
} 

func (presigner *Presigner) PutObject(
	ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}
