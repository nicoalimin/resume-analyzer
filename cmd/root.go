/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/nicoalimin/resume-analyzer/textract"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var inputDir string
var outputDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "resume-analyzer",
	Short: "A CLI tool to OCR PDFs using AWS Textract",
	Long:  `Processes all PDFs in a folder using AWS Textract and saves the extracted text to another folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		if inputDir == "" || outputDir == "" {
			fmt.Fprintln(os.Stderr, "Both --input and --output folders must be specified.")
			os.Exit(1)
		}

		files, err := os.ReadDir(inputDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read input directory: %v\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			if file.IsDir() || len(file.Name()) < 4 || file.Name()[len(file.Name())-4:] != ".pdf" {
				continue
			}
			pdfPath := inputDir + string(os.PathSeparator) + file.Name()
			outputPath := outputDir + string(os.PathSeparator) + file.Name()[:len(file.Name())-4] + ".txt"

			fmt.Printf("Processing %s...\n", file.Name())
			extractedText, err := textract.ExtractTextFromPDF(pdfPath)
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.resume-analyzer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&inputDir, "input", "i", "", "Input folder containing PDFs")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output folder for extracted text")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".resume-analyzer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".resume-analyzer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
