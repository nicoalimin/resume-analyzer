package textract

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// ExtractTextFromPDF calls AWS Textract DetectDocumentText on the given PDF file and returns the extracted text.
func ExtractTextFromPDF(pdfPath string) (string, error) {
	// Read the PDF file into memory
	fileBytes, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := textract.NewFromConfig(cfg)

	input := &textract.DetectDocumentTextInput{
		Document: &types.Document{
			Bytes: fileBytes,
		},
	}

	resp, err := client.DetectDocumentText(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("Textract DetectDocumentText failed: %w", err)
	}

	// Concatenate all detected lines into a single string
	var result string
	for _, block := range resp.Blocks {
		if block.BlockType == "LINE" && block.Text != nil {
			result += *block.Text + "\n"
		}
	}
	return result, nil
}
