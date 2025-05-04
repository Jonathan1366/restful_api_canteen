package usecase

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
)

var textractClient *textract.Client

func InitTextractClient(region string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("failed to load AWS SDK config: %v", err)
	}
	textractClient = textract.NewFromConfig(cfg)
	return nil
}

func GetTextractClient() *textract.Client  {
	if textractClient == nil{
		panic("Textract client is not initialized. Call initextractclient first")
	}
	return textractClient
}

