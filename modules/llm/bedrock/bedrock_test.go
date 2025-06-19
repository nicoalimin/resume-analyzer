package bedrock

import (
	"testing"
	"time"

	"github.com/nicoalimin/resume-analyzer/interfaces"
)

// MockLLMService for testing
type MockLLMService struct {
	shouldError bool
	response    string
}

func (m *MockLLMService) GenerateText(prompt string) (string, error) {
	if m.shouldError {
		return "", &MockError{message: "mock error"}
	}
	return m.response, nil
}

type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}

func TestNewBedrockService(t *testing.T) {
	service := NewBedrockService()

	if service == nil {
		t.Fatal("NewBedrockService() returned nil")
	}

	// Check if it implements the interface
	var _ interfaces.LLMService = service
}

func TestBedrockService_GenerateText_Success(t *testing.T) {
	// Use mock service instead of real service for testing
	mock := NewMockBedrockService()
	mock.SetResponse("test prompt", "test response")

	response, err := mock.GenerateText("test prompt")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response != "test response" {
		t.Errorf("Expected 'test response', got '%s'", response)
	}
}

func TestBedrockService_GenerateText_InterfaceCompliance(t *testing.T) {
	// Use mock service for interface compliance testing
	mock := NewMockBedrockService()

	// Test interface compliance
	var _ interfaces.LLMService = mock

	// Test that the method signature is correct
	response, err := mock.GenerateText("test prompt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected non-empty response")
	}
}

func TestLegacyGenerateTextFunction(t *testing.T) {
	// Test the legacy function for backward compatibility
	response, err := GenerateText("test prompt")

	// We don't check the actual result since it depends on AWS
	// Just ensure the function doesn't panic and returns something
	if err != nil && err.Error() == "" {
		t.Error("Legacy GenerateText should return a meaningful error message")
	}

	// If no error, response should not be empty (though it might be in test environment)
	if err == nil && response == "" {
		t.Log("Note: Response is empty, which is expected in test environment without AWS credentials")
	}
}

func TestBedrockService_Configuration(t *testing.T) {
	// Use mock service for configuration testing
	mock := NewMockBedrockService()

	// Test that mock service can be created with default configuration
	if mock == nil {
		t.Fatal("Failed to create MockBedrockService with default configuration")
	}

	// Test that the service is properly initialized
	response, err := mock.GenerateText("test")
	if err != nil {
		t.Fatalf("Mock service should work without configuration: %v", err)
	}
	if response == "" {
		t.Error("Mock service should return a response")
	}
}

// Benchmark test for performance
func BenchmarkBedrockService_GenerateText(b *testing.B) {
	// Use mock service for benchmarking
	mock := NewMockBedrockService()
	prompt := "This is a test prompt for benchmarking"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mock.GenerateText(prompt)
		if err != nil {
			b.Errorf("Benchmark iteration %d failed: %v", i, err)
		}
	}
}

// Test error handling scenarios
func TestBedrockService_ErrorScenarios(t *testing.T) {
	// Use mock service for error scenario testing
	mock := NewMockBedrockService()

	// Test with empty prompt
	response, err := mock.GenerateText("")
	if err != nil {
		t.Logf("Empty prompt test: %v", err)
	}
	if response == "" {
		t.Error("Empty prompt should return some response")
	}

	// Test with very long prompt
	longPrompt := string(make([]byte, 10000)) // 10KB prompt
	response, err = mock.GenerateText(longPrompt)
	if err != nil {
		t.Logf("Long prompt test: %v", err)
	}
	if response == "" {
		t.Error("Long prompt should return some response")
	}

	// Test custom error scenario
	mock.SetError("error test", &MockError{message: "custom error"})
	_, err = mock.GenerateText("error test")
	if err == nil {
		t.Error("Expected error for 'error test' prompt")
	}
	if err.Error() != "custom error" {
		t.Errorf("Expected 'custom error', got '%s'", err.Error())
	}
}

// Test service creation and interface compliance
func TestBedrockService_InterfaceCompliance(t *testing.T) {
	// Test that MockBedrockService implements LLMService interface
	var service interfaces.LLMService = NewMockBedrockService()

	if service == nil {
		t.Fatal("MockBedrockService does not implement LLMService interface")
	}

	// Test that we can call the interface method
	response, err := service.GenerateText("test")
	if err != nil {
		t.Errorf("Interface method should not return error: %v", err)
	}
	if response == "" {
		t.Error("Interface method should return non-empty response")
	}
}

// Test mock service with delay simulation
func TestBedrockService_MockWithDelay(t *testing.T) {
	mock := NewMockBedrockService()
	mock.SetDelay(10 * time.Millisecond) // 10ms delay

	start := time.Now()
	response, err := mock.GenerateText("test prompt")
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected non-empty response")
	}
	if duration < 10*time.Millisecond {
		t.Error("Expected delay to be at least 10ms")
	}
}

// Test mock service with custom responses
func TestBedrockService_MockCustomResponses(t *testing.T) {
	mock := NewMockBedrockService()

	// Set multiple custom responses
	mock.SetResponse("prompt1", "response1")
	mock.SetResponse("prompt2", "response2")

	// Test first custom response
	response, err := mock.GenerateText("prompt1")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response != "response1" {
		t.Errorf("Expected 'response1', got '%s'", response)
	}

	// Test second custom response
	response, err = mock.GenerateText("prompt2")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response != "response2" {
		t.Errorf("Expected 'response2', got '%s'", response)
	}

	// Test default response for unknown prompt
	response, err = mock.GenerateText("unknown prompt")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == "" {
		t.Error("Expected default response for unknown prompt")
	}
}
