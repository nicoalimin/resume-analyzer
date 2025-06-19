/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicoalimin/resume-analyzer/interfaces"
	"github.com/nicoalimin/resume-analyzer/modules/llm/bedrock"
	"github.com/spf13/cobra"
)

var queryPrompt string
var queryInputDir string
var queryOutputFile string
var queryLLMService interfaces.LLMService

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query all resume texts with a custom prompt using AWS Bedrock",
	Long: `Reads all .txt files from a folder, combines them into a single prompt,
and sends it to AWS Bedrock with your custom question.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Initialize LLM service if not already set
		if queryLLMService == nil {
			queryLLMService = bedrock.NewBedrockService()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if queryPrompt == "" {
			fmt.Fprintln(os.Stderr, "A prompt must be specified with --prompt.")
			os.Exit(1)
		}

		if queryInputDir == "" {
			fmt.Fprintln(os.Stderr, "Input directory must be specified with --input.")
			os.Exit(1)
		}

		// Read all text files from the input directory
		files, err := os.ReadDir(queryInputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input directory: %v\n", err)
			os.Exit(1)
		}

		var allTexts []string
		var fileNames []string

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
				continue
			}

			inputPath := filepath.Join(queryInputDir, file.Name())
			fmt.Printf("Reading %s...\n", file.Name())

			// Read the text file
			content, err := os.ReadFile(inputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file.Name(), err)
				continue
			}

			allTexts = append(allTexts, string(content))
			fileNames = append(fileNames, file.Name())
		}

		if len(allTexts) == 0 {
			fmt.Fprintln(os.Stderr, "No .txt files found in the input directory.")
			os.Exit(1)
		}

		// Combine all texts into a single prompt
		combinedPrompt := buildCombinedPrompt(queryPrompt, allTexts, fileNames)

		fmt.Printf("Sending query to Bedrock with %d resume files...\n", len(allTexts))

		// Send to Bedrock
		response, err := queryLLMService.GenerateText(combinedPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bedrock query failed: %v\n", err)
			os.Exit(1)
		}

		// Output the response
		if queryOutputFile != "" {
			// Write to file
			err = os.WriteFile(queryOutputFile, []byte(response), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write response to file: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Response saved to %s\n", queryOutputFile)
		} else {
			// Print to stdout with better formatting
			fmt.Println("\n" + strings.Repeat("=", 80))
			fmt.Println("BEDROCK RESPONSE")
			fmt.Println(strings.Repeat("=", 80))
			fmt.Println(response)
			fmt.Println(strings.Repeat("=", 80))
		}
	},
}

// buildCombinedPrompt creates a comprehensive prompt combining the user's question with all resume texts
func buildCombinedPrompt(userPrompt string, texts []string, fileNames []string) string {
	var prompt strings.Builder

	prompt.WriteString("You are analyzing multiple resumes. Below are the extracted texts from ")
	prompt.WriteString(fmt.Sprintf("%d resume files.\n\n", len(texts)))

	prompt.WriteString("User Question: ")
	prompt.WriteString(userPrompt)
	prompt.WriteString("\n\n")

	prompt.WriteString("Resume Texts:\n")
	prompt.WriteString(strings.Repeat("=", 50))
	prompt.WriteString("\n\n")

	for i, text := range texts {
		prompt.WriteString(fmt.Sprintf("--- Resume %d: %s ---\n", i+1, fileNames[i]))
		prompt.WriteString(text)
		prompt.WriteString("\n\n")
	}

	prompt.WriteString(strings.Repeat("=", 50))
	prompt.WriteString("\n\n")
	prompt.WriteString("Please provide a comprehensive answer to the user's question based on the resume texts above. ")
	prompt.WriteString("If the question requires comparing candidates, please provide detailed analysis and comparisons. ")
	prompt.WriteString("If the question asks for specific information, please extract and present it clearly.\n\n")
	prompt.WriteString("Answer:")

	return prompt.String()
}

// SetQueryLLMService allows dependency injection of LLM service (useful for testing)
func SetQueryLLMService(service interfaces.LLMService) {
	queryLLMService = service
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().StringVarP(&queryPrompt, "prompt", "p", "", "The question or prompt to ask about the resumes")
	queryCmd.Flags().StringVarP(&queryInputDir, "input", "i", "", "Input folder containing .txt files")
	queryCmd.Flags().StringVarP(&queryOutputFile, "output", "o", "", "Output file for the response (optional, prints to stdout if not specified)")

	// Mark required flags
	queryCmd.MarkFlagRequired("prompt")
	queryCmd.MarkFlagRequired("input")
}
