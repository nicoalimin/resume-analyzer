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

var consolidateInputDir string
var consolidateOutputFile string
var consolidateLLMService interfaces.LLMService

type ApplicantInfo struct {
	Name            string
	Role            string
	Seniority       string
	Status          string
	CurrentPosition string
	CurrentCompany  string
	YearsOfExp      string
	CVLink          string
	Skillset        string
	Remarks         string
}

// consolidateCmd represents the consolidate command
var consolidateCmd = &cobra.Command{
	Use:   "consolidate",
	Short: "Consolidate all summaries into a single summary table",
	Long:  `Reads all summary files and generates a consolidated table with applicant information.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Initialize LLM service if not already set
		if consolidateLLMService == nil {
			consolidateLLMService = bedrock.NewBedrockService()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if consolidateInputDir == "" || consolidateOutputFile == "" {
			fmt.Fprintln(os.Stderr, "Both --input and --output must be specified.")
			os.Exit(1)
		}

		files, err := os.ReadDir(consolidateInputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input directory: %v\n", err)
			os.Exit(1)
		}

		var applicants []ApplicantInfo

		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), "_summary.txt") {
				continue
			}

			inputPath := filepath.Join(consolidateInputDir, file.Name())
			fmt.Printf("Processing %s...\n", file.Name())

			// Read the summary file
			content, err := os.ReadFile(inputPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", file.Name(), err)
				continue
			}

			// Extract structured information using LLM service
			applicant, err := extractApplicantInfo(string(content), file.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to extract info from %s: %v\n", file.Name(), err)
				continue
			}

			applicants = append(applicants, applicant)
		}

		// Generate the consolidated table
		table := generateConsolidatedTable(applicants)

		// Write to output file
		err = os.WriteFile(consolidateOutputFile, []byte(table), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write consolidated table: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Consolidated table saved to %s\n", consolidateOutputFile)
	},
}

// SetConsolidateLLMService allows dependency injection of LLM service (useful for testing)
func SetConsolidateLLMService(service interfaces.LLMService) {
	consolidateLLMService = service
}

func extractApplicantInfo(summary string, filename string) (ApplicantInfo, error) {
	prompt := prompts.GetExtractionPrompt(summary)

	response, err := consolidateLLMService.GenerateText(prompt)
	if err != nil {
		return ApplicantInfo{}, err
	}

	// Parse the JSON response (simplified - in production you'd want proper JSON parsing)
	// For now, we'll create a basic structure
	applicant := ApplicantInfo{
		Name:            extractField(response, "name"),
		Role:            extractField(response, "role"),
		Seniority:       extractField(response, "seniority"),
		Status:          extractField(response, "status"),
		CurrentPosition: extractField(response, "current_position"),
		CurrentCompany:  extractField(response, "current_company"),
		YearsOfExp:      extractField(response, "years_of_exp"),
		CVLink:          extractField(response, "cv_link"),
		Skillset:        extractField(response, "skillset"),
		Remarks:         extractField(response, "remarks"),
	}

	// If name is not found, use filename
	if applicant.Name == "N/A" || applicant.Name == "" {
		applicant.Name = strings.TrimSuffix(filename, "_summary.txt")
	}

	return applicant, nil
}

func extractField(response, field string) string {
	// Simple field extraction - in production you'd want proper JSON parsing
	fieldLower := strings.ToLower(field)
	responseLower := strings.ToLower(response)

	if strings.Contains(responseLower, `"`+fieldLower+`"`) {
		// Basic extraction - find the field and get the value
		start := strings.Index(responseLower, `"`+fieldLower+`"`)
		if start != -1 {
			valueStart := strings.Index(response[start:], ":")
			if valueStart != -1 {
				valueStart += start + 1
				valueEnd := strings.Index(response[valueStart:], "\n")
				if valueEnd == -1 {
					valueEnd = len(response) - valueStart
				}
				value := strings.TrimSpace(response[valueStart : valueStart+valueEnd])
				value = strings.Trim(value, `",`)
				return value
			}
		}
	}
	return "N/A"
}

func generateConsolidatedTable(applicants []ApplicantInfo) string {
	var csv strings.Builder

	// CSV header
	csv.WriteString("Applicant,Role,Seniority,Status,Current Position,Current Company,Years of Exp,CV Link,Skillset,Remarks\n")

	// Data rows
	for _, applicant := range applicants {
		// Escape CSV fields that contain commas or quotes
		row := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			escapeCSVField(applicant.Name),
			escapeCSVField(applicant.Role),
			escapeCSVField(applicant.Seniority),
			escapeCSVField(applicant.Status),
			escapeCSVField(applicant.CurrentPosition),
			escapeCSVField(applicant.CurrentCompany),
			escapeCSVField(applicant.YearsOfExp),
			escapeCSVField(applicant.CVLink),
			escapeCSVField(applicant.Skillset),
			escapeCSVField(applicant.Remarks))
		csv.WriteString(row)
	}

	return csv.String()
}

// escapeCSVField properly escapes CSV fields that contain commas, quotes, or newlines
func escapeCSVField(field string) string {
	// If field contains comma, quote, or newline, wrap in quotes and escape internal quotes
	if strings.ContainsAny(field, ",\"\n\r") {
		// Replace any existing quotes with double quotes
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return field
}

func init() {
	rootCmd.AddCommand(consolidateCmd)
	consolidateCmd.Flags().StringVarP(&consolidateInputDir, "input", "i", "", "Input folder containing summary files")
	consolidateCmd.Flags().StringVarP(&consolidateOutputFile, "output", "o", "", "Output file for consolidated table")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// consolidateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// consolidateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
