package interfaces

// OCRService defines the interface for Optical Character Recognition services
type OCRService interface {
	// ExtractTextFromPDF extracts text from a PDF file
	// Returns the extracted text as a string and any error
	ExtractTextFromPDF(pdfPath string) (string, error)
}

// LLMService defines the interface for Large Language Model services
type LLMService interface {
	// GenerateText generates text based on a prompt
	// Returns the generated text and any error
	GenerateText(prompt string) (string, error)
}
