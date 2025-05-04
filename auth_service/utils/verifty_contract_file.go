package utils

import (
	"context"
	"fmt"
	"strings"
	"ubm-canteen/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

  func ProcessDocumentWithTextractFromS3(ctx context.Context, Bucket, filePath string, keywords []string) ([]string, error) {

    // Textract client
    client := usecase.GetTextractClient()

    if !strings.HasSuffix(strings.ToLower(filePath), ".pdf") && !strings.HasSuffix(strings.ToLower(filePath), ".jpg") &&
    !strings.HasSuffix(strings.ToLower(filePath), ".jpeg") &&
    !strings.HasSuffix(strings.ToLower(filePath), ".png") {
    return nil, fmt.Errorf("invalid file type")
    }
    //send request to textract
    input := &textract.AnalyzeDocumentInput{
      Document:&types.Document{
        S3Object: &types.S3Object{
          Bucket: aws.String(Bucket),
          Name: aws.String(filePath),
        },
      },
      FeatureTypes: []types.FeatureType{types.FeatureTypeForms},
    }
    
    // Call AWS Textract to analyze the document.
    resp, err := client.AnalyzeDocument(context.TODO(), input)
    if err != nil {
      return nil, fmt.Errorf("failed to detect document text: %v", err)
    }

    //process the result & return it
    var lines[]string
    for _, block := range resp.Blocks{
      if block.BlockType == types.BlockTypeLine &&block.Text != nil {
        lines = append(lines, *block.Text)
      }
    }

    //check if the keywords are present in the document
    var foundKeywords [] string
    for _, keyword := range keywords{
      for _, line := range lines{
        if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
          foundKeywords = append(foundKeywords, keyword)
          break
        }
      }
    }

    //if no keywords are found, return an error
    if len(foundKeywords) == 0 {
      return nil, fmt.Errorf("no matching keywords found in the document")
    }
    return foundKeywords, nil
  }