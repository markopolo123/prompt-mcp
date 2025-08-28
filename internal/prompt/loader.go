package prompt

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Loader handles loading prompts from the filesystem
type Loader struct {
	promptsDir string
}

// NewLoader creates a new prompt loader
func NewLoader(promptsDir string) *Loader {
	return &Loader{
		promptsDir: promptsDir,
	}
}

// LoadAllPrompts loads all prompts from the prompts directory
func (l *Loader) LoadAllPrompts() (*PromptLibrary, error) {
	library := NewPromptLibrary()

	err := filepath.Walk(l.promptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-YAML files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".yaml") {
			return nil
		}

		prompt, err := l.LoadPrompt(path)
		if err != nil {
			return fmt.Errorf("failed to load prompt from %s: %w", path, err)
		}

		// Store the file path for reference
		prompt.FilePath = path

		// Check for duplicate IDs
		if _, exists := library.GetPrompt(prompt.Metadata.ID); exists {
			return fmt.Errorf("duplicate prompt ID '%s' found in file %s", prompt.Metadata.ID, path)
		}

		library.AddPrompt(prompt)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load prompts from directory %s: %w", l.promptsDir, err)
	}

	return library, nil
}

// LoadPrompt loads a single prompt from a file
func (l *Loader) LoadPrompt(filePath string) (*Prompt, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var prompt Prompt
	if err := yaml.Unmarshal(data, &prompt); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate the loaded prompt
	if err := ValidatePrompt(&prompt); err != nil {
		return nil, fmt.Errorf("prompt validation failed: %w", err)
	}

	return &prompt, nil
}

// SavePrompt saves a prompt to a file
func (l *Loader) SavePrompt(prompt *Prompt, filePath string) error {
	// Validate before saving
	if err := ValidatePrompt(prompt); err != nil {
		return fmt.Errorf("prompt validation failed: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(prompt)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetCategoryFromPath extracts the category from a file path
func (l *Loader) GetCategoryFromPath(filePath string) string {
	// Get relative path from prompts directory
	relPath, err := filepath.Rel(l.promptsDir, filePath)
	if err != nil {
		return "uncategorized"
	}

	// Extract the first directory component as category
	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) > 1 {
		return parts[0]
	}

	return "uncategorized"
}

// GeneratePromptURI generates a URI for a prompt based on its file path
func (l *Loader) GeneratePromptURI(filePath string) string {
	category := l.GetCategoryFromPath(filePath)
	fileName := filepath.Base(filePath)
	promptName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	
	return fmt.Sprintf("prompt://%s/%s", category, promptName)
}