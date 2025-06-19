package cmd

import (
	"fmt"
	"os"

	"github.com/nicoalimin/resume-analyzer/interfaces"
	"github.com/nicoalimin/resume-analyzer/textract"
	"github.com/spf13/cobra"
)

var convertInputDir string
var convertOutputDir string
var ocrService interfaces.OCRService

var convertPDFsCmd = &cobra.Command{
	Use:   "convert-pdfs",
	Short: "Convert PDFs in a folder to text using AWS Textract",
	Long:  `Processes all PDFs in a folder using AWS Textract and saves the extracted text to another folder.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Initialize OCR service if not already set
		if ocrService == nil {
			ocrService = textract.NewTextractService()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if convertInputDir == "" || convertOutputDir == "" {
			fmt.Fprintln(os.Stderr, "Both --input and --output folders must be specified.")
			os.Exit(1)
		}

		files, err := os.ReadDir(convertInputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input directory: %v\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			if file.IsDir() || len(file.Name()) < 4 || file.Name()[len(file.Name())-4:] != ".pdf" {
				continue
			}
			pdfPath := convertInputDir + string(os.PathSeparator) + file.Name()
			outputPath := convertOutputDir + string(os.PathSeparator) + file.Name()[:len(file.Name())-4] + ".txt"

			fmt.Printf("Processing %s...\n", file.Name())
			extractedText, err := ocrService.ExtractTextFromPDF(pdfPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Textract failed for %s: %v\n", file.Name(), err)
				continue
			}

			err = os.WriteFile(outputPath, []byte(extractedText), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write output for %s: %v\n", file.Name(), err)
			}
		}
		fmt.Println("Processing complete.")
	},
}

// SetOCRService allows dependency injection of OCR service (useful for testing)
func SetOCRService(service interfaces.OCRService) {
	ocrService = service
}

func init() {
	rootCmd.AddCommand(convertPDFsCmd)
	convertPDFsCmd.Flags().StringVarP(&convertInputDir, "input", "i", "", "Input folder containing PDFs")
	convertPDFsCmd.Flags().StringVarP(&convertOutputDir, "output", "o", "", "Output folder for extracted text")
}
