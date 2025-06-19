package bedrock

import (
	"testing"
	"time"

	"github.com/nicoalimin/resume-analyzer/interfaces"
)

// MockBedrockService provides a mock implementation for testing
type MockBedrockService struct {
	responses map[string]string
	errors    map[string]error
	delay     time.Duration
}

// NewMockBedrockService creates a new mock service
func NewMockBedrockService() *MockBedrockService {
	return &MockBedrockService{
		responses: make(map[string]string),
		errors:    make(map[string]error),
		delay:     0,
	}
}

// SetResponse sets a mock response for a specific prompt
func (m *MockBedrockService) SetResponse(prompt, response string) {
	m.responses[prompt] = response
}

// SetError sets a mock error for a specific prompt
func (m *MockBedrockService) SetError(prompt string, err error) {
	m.errors[prompt] = err
}

// SetDelay sets a delay to simulate network latency
func (m *MockBedrockService) SetDelay(delay time.Duration) {
	m.delay = delay
}

// GenerateText implements the LLMService interface
func (m *MockBedrockService) GenerateText(prompt string) (string, error) {
	// Simulate network delay
	if m.delay > 0 {
		time.Sleep(m.delay)
	}

	// Check if we have a specific error for this prompt
	if err, exists := m.errors[prompt]; exists {
		return "", err
	}

	// Check if we have a specific response for this prompt
	if response, exists := m.responses[prompt]; exists {
		return response, nil
	}

	// Default mock response
	return "Mock response for: " + prompt, nil
}

// Test with mock service
func TestBedrockService_WithMock(t *testing.T) {
	mock := NewMockBedrockService()

	// Test basic functionality
	response, err := mock.GenerateText("test prompt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected non-empty response")
	}

	// Test custom response
	mock.SetResponse("custom prompt", "custom response")
	response, err = mock.GenerateText("custom prompt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response != "custom response" {
		t.Errorf("Expected 'custom response', got '%s'", response)
	}

	// Test custom error
	mock.SetError("error prompt", &MockError{message: "test error"})
	_, err = mock.GenerateText("error prompt")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "test error" {
		t.Errorf("Expected 'test error', got '%s'", err.Error())
	}
}

// Test interface compliance with mock
func TestBedrockService_MockInterfaceCompliance(t *testing.T) {
	mock := NewMockBedrockService()

	// Test that mock implements the interface
	var service interfaces.LLMService = mock
	if service == nil {
		t.Fatal("MockBedrockService does not implement LLMService interface")
	}

	// Test interface method
	response, err := service.GenerateText("test")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected non-empty response")
	}
}

// Test performance with mock
func BenchmarkBedrockService_Mock(b *testing.B) {
	mock := NewMockBedrockService()
	prompt := "Benchmark test prompt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mock.GenerateText(prompt)
		if err != nil {
			b.Errorf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// Test concurrent access with mock
func TestBedrockService_MockConcurrent(t *testing.T) {
	mock := NewMockBedrockService()

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			response, err := mock.GenerateText("concurrent test")
			if err != nil {
				t.Errorf("Concurrent test %d failed: %v", id, err)
			}
			if response == "" {
				t.Errorf("Concurrent test %d got empty response", id)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test edge cases with mock
func TestBedrockService_MockEdgeCases(t *testing.T) {
	mock := NewMockBedrockService()

	// Test empty prompt
	response, err := mock.GenerateText("")
	if err != nil {
		t.Errorf("Empty prompt test failed: %v", err)
	}
	if response == "" {
		t.Error("Empty prompt should return some response")
	}

	// Test very long prompt
	longPrompt := string(make([]byte, 100000)) // 100KB prompt
	response, err = mock.GenerateText(longPrompt)
	if err != nil {
		t.Errorf("Long prompt test failed: %v", err)
	}
	if response == "" {
		t.Error("Long prompt should return some response")
	}

	// Test special characters
	specialPrompt := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	response, err = mock.GenerateText(specialPrompt)
	if err != nil {
		t.Errorf("Special characters test failed: %v", err)
	}
	if response == "" {
		t.Error("Special characters should return some response")
	}
}
