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
	"github.com/nicoalimin/resume-analyzer/prompts"
	"github.com/spf13/cobra"
)

var summarizeInputDir string
var summarizeOutputDir string
var llmService interfaces.LLMService

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Generate summaries of extracted text documents using AWS Bedrock",
	Long: `Reads all .txt files from a folder, generates summaries using AWS Bedrock,
and saves the summaries to an output folder.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Initialize LLM service if not already set
		if llmService == nil {
			llmService = bedrock.NewBedrockService()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if summarizeInputDir == "" || summarizeOutputDir == "" {
			fmt.Fprintln(os.Stderr, "Both --input and --output folders must be specified.")
			os.Exit(1)
		}

		files, err := os.ReadDir(summarizeInputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input directory: %v\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".txt") {
				continue
			}

			inputPath := filepath.Join(summarizeInputDir, file.Name())
			outputPath := filepath.Join(summarizeOutputDir, strings.TrimSuffix(file.Name(), ".txt")+"_summary.txt")

			fmt.Printf("Summarizing %s...\n", file.Name())

			// Read the text file
			content, err := os.ReadFile(inputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file.Name(), err)
				continue
			}

			// Generate summary using LLM service
			prompt := prompts.GetSummaryPrompt(string(content))
			summary, err := llmService.GenerateText(prompt)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Bedrock failed for %s: %v\n", file.Name(), err)
				continue
			}

			// Write the summary to output file
			err = os.WriteFile(outputPath, []byte(summary), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write summary for %s: %v\n", file.Name(), err)
			} else {
				fmt.Printf("Summary saved to %s\n", outputPath)
			}
		}
		fmt.Println("Summarization complete.")
	},
}

// SetLLMService allows dependency injection of LLM service (useful for testing)
func SetLLMService(service interfaces.LLMService) {
	llmService = service
}

func init() {
	rootCmd.AddCommand(summarizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// summarizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// summarizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	summarizeCmd.Flags().StringVarP(&summarizeInputDir, "input", "i", "", "Input folder containing .txt files")
	summarizeCmd.Flags().StringVarP(&summarizeOutputDir, "output", "o", "", "Output folder for summaries")
}
