package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrompt(t *testing.T) {
	// Create a temporary directory for test
	tempDir, err := os.MkdirTemp("", "prompt-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test prompt file
	testPromptContent := `metadata:
  id: "test-prompt"
  name: "Test Prompt"
  description: "A test prompt"
  author: "test"
  created: "2025-08-27T10:00:00Z"
  modified: "2025-08-27T10:00:00Z"
  version: "1.0.0"
  tags:
    - "test"

arguments:
  - name: "name"
    description: "Test name"
    type: "string"
    required: true

prompt: |
  Hello {{name}}!

usage_stats:
  usage_count: 0
  last_used: "2025-08-27T10:00:00Z"
`

	testFile := filepath.Join(tempDir, "test.yaml")
	err = os.WriteFile(testFile, []byte(testPromptContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading the prompt
	loader := NewLoader(tempDir)
	prompt, err := loader.LoadPrompt(testFile)
	if err != nil {
		t.Fatalf("Failed to load prompt: %v", err)
	}

	// Verify prompt content
	if prompt.Metadata.ID != "test-prompt" {
		t.Errorf("Expected ID 'test-prompt', got '%s'", prompt.Metadata.ID)
	}

	if prompt.Metadata.Name != "Test Prompt" {
		t.Errorf("Expected name 'Test Prompt', got '%s'", prompt.Metadata.Name)
	}

	if len(prompt.Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(prompt.Arguments))
	}

	if prompt.Arguments[0].Name != "name" {
		t.Errorf("Expected argument name 'name', got '%s'", prompt.Arguments[0].Name)
	}

	if prompt.Prompt != "Hello {{name}}!\n" {
		t.Errorf("Expected prompt 'Hello {{name}}!\\n', got '%s'", prompt.Prompt)
	}
}

func TestLoadAllPrompts(t *testing.T) {
	// Test loading from the actual prompts directory
	loader := NewLoader("../../prompts")
	library, err := loader.LoadAllPrompts()
	if err != nil {
		t.Fatalf("Failed to load prompts: %v", err)
	}

	if len(library.Prompts) == 0 {
		t.Error("Expected some prompts to be loaded")
	}

	// Verify we can get prompts by ID
	for id := range library.Prompts {
		prompt, exists := library.GetPrompt(id)
		if !exists {
			t.Errorf("Prompt %s should exist", id)
		}
		if prompt.Metadata.ID != id {
			t.Errorf("Prompt ID mismatch: expected %s, got %s", id, prompt.Metadata.ID)
		}
	}
}