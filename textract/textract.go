package textract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var ctx = context.Background()

// ExtractTextFromPDF calls AWS Textract DetectDocumentText on the given PDF file and returns the extracted text.
func ExtractTextFromPDF(pdfPath string) (string, error) {
	// Split PDF into single-page PDFs in a temp dir
	tempDir, err := os.MkdirTemp("", "pdfpages")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	err = api.SplitFile(pdfPath, tempDir, 1, nil)
	if err != nil {
		return "", fmt.Errorf("failed to split PDF: %w", err)
	}

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to read temp dir: %w", err)
	}

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := textract.NewFromConfig(cfg)

	var combinedText string
	for _, f := range files {
		pagePath := filepath.Join(tempDir, f.Name())
		pageBytes, err := os.ReadFile(pagePath)
		if err != nil {
			return "", fmt.Errorf("failed to read page PDF: %w", err)
		}
		if len(pageBytes) > 5*1024*1024 {
			return "", fmt.Errorf("PDF page is too large: %d bytes (max 5MB)", len(pageBytes))
		}
		input := &textract.DetectDocumentTextInput{
			Document: &types.Document{
				Bytes: pageBytes,
			},
		}
		resp, err := client.DetectDocumentText(ctx, input)
		if err != nil {
			return "", fmt.Errorf("textract failed for page %s: %w", f.Name(), err)
		}
		for _, block := range resp.Blocks {
			if block.BlockType == "LINE" && block.Text != nil {
				combinedText += *block.Text + "\n"
			}
		}
	}
	return combinedText, nil
}
