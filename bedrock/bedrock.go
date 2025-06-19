package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/spf13/viper"
)

var ctx = context.Background()

// SummaryRequest represents the request structure for Claude
type SummaryRequest struct {
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens"`
	Temperature      float64   `json:"temperature"`
	AnthropicVersion string    `json:"anthropic_version"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SummaryResponse represents the response structure from Claude
type SummaryResponse struct {
	Content []Content `json:"content"`
}

// Content represents the content in the response
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// GenerateText calls AWS Bedrock with any prompt and returns the response
func GenerateText(prompt string) (string, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	// Create the request
	request := SummaryRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:        1000,
		Temperature:      0.3,
		AnthropicVersion: viper.GetString("anthropic_version"),
	}
	if request.AnthropicVersion == "" {
		request.AnthropicVersion = "bedrock-2023-05-31"
	}

	// Marshal the request to JSON
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Bedrock (using Claude 3 Sonnet)
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(viper.GetString("bedrock_model_id")),
		ContentType: aws.String("application/json"),
		Body:        requestBytes,
	}
	if *input.ModelId == "" {
		*input.ModelId = "anthropic.claude-3-5-sonnet-20240620-v1:0"
	}

	resp, err := client.InvokeModel(ctx, input)
	if err != nil {
		return "", fmt.Errorf("bedrock invoke failed: %w", err)
	}

	// Parse the response
	var summaryResp SummaryResponse
	err = json.Unmarshal(resp.Body, &summaryResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(summaryResp.Content) == 0 {
		return "", fmt.Errorf("no content in response")
	}

	return summaryResp.Content[0].Text, nil
}
