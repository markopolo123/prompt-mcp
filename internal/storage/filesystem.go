package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markopolo123/prompt-mcp/internal/prompt"
)

// FileSystemStorage provides filesystem-based storage for prompts
type FileSystemStorage struct {
	loader *prompt.Loader
	config Config
}

// Config holds storage configuration
type Config struct {
	PromptsDir   string
	WatchChanges bool
}

// NewFileSystemStorage creates a new filesystem storage instance
func NewFileSystemStorage(config Config) (*FileSystemStorage, error) {
	// Ensure prompts directory exists
	if err := ensureDirectoryExists(config.PromptsDir); err != nil {
		return nil, fmt.Errorf("failed to ensure prompts directory exists: %w", err)
	}

	loader := prompt.NewLoader(config.PromptsDir)

	return &FileSystemStorage{
		loader: loader,
		config: config,
	}, nil
}

// LoadLibrary loads all prompts from the filesystem
func (fs *FileSystemStorage) LoadLibrary() (*prompt.PromptLibrary, error) {
	return fs.loader.LoadAllPrompts()
}

// LoadPrompt loads a specific prompt by file path
func (fs *FileSystemStorage) LoadPrompt(filePath string) (*prompt.Prompt, error) {
	return fs.loader.LoadPrompt(filePath)
}

// SavePrompt saves a prompt to the filesystem
func (fs *FileSystemStorage) SavePrompt(p *prompt.Prompt, filePath string) error {
	return fs.loader.SavePrompt(p, filePath)
}

// GetPromptsDir returns the configured prompts directory
func (fs *FileSystemStorage) GetPromptsDir() string {
	return fs.config.PromptsDir
}

// ListPromptFiles returns all prompt files in the prompts directory
func (fs *FileSystemStorage) ListPromptFiles() ([]string, error) {
	var files []string

	err := filepath.Walk(fs.config.PromptsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list prompt files: %w", err)
	}

	return files, nil
}

// GetPromptURI generates a URI for a prompt based on its file path
func (fs *FileSystemStorage) GetPromptURI(filePath string) string {
	return fs.loader.GeneratePromptURI(filePath)
}

// GetCategoryFromPath extracts category from file path
func (fs *FileSystemStorage) GetCategoryFromPath(filePath string) string {
	return fs.loader.GetCategoryFromPath(filePath)
}

// ensureDirectoryExists creates directory if it doesn't exist
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}