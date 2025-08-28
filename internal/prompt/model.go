package prompt

import (
	"time"
)

// Prompt represents a prompt template with metadata
type Prompt struct {
	Metadata    Metadata     `yaml:"metadata"`
	Arguments   []Argument   `yaml:"arguments,omitempty"`
	Prompt      string       `yaml:"prompt"`
	UsageStats  UsageStats   `yaml:"usage_stats"`
	FilePath    string       `yaml:"-"` // Internal field, not serialized
}

// Metadata contains prompt metadata
type Metadata struct {
	ID          string    `yaml:"id"`
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Author      string    `yaml:"author"`
	Created     time.Time `yaml:"created"`
	Modified    time.Time `yaml:"modified"`
	Version     string    `yaml:"version"`
	Tags        []string  `yaml:"tags,omitempty"`
}

// Argument represents a prompt argument/parameter
type Argument struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Type        ArgumentType `yaml:"type"`
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default,omitempty"`
}

// ArgumentType defines the types of arguments supported
type ArgumentType string

const (
	ArgumentTypeString  ArgumentType = "string"
	ArgumentTypeNumber  ArgumentType = "number"
	ArgumentTypeBoolean ArgumentType = "boolean"
)

// UsageStats tracks usage statistics for a prompt
type UsageStats struct {
	UsageCount int       `yaml:"usage_count"`
	LastUsed   time.Time `yaml:"last_used"`
}

// PromptLibrary holds all loaded prompts
type PromptLibrary struct {
	Prompts map[string]*Prompt // keyed by prompt ID
}

// NewPromptLibrary creates a new prompt library
func NewPromptLibrary() *PromptLibrary {
	return &PromptLibrary{
		Prompts: make(map[string]*Prompt),
	}
}

// AddPrompt adds a prompt to the library
func (pl *PromptLibrary) AddPrompt(prompt *Prompt) {
	pl.Prompts[prompt.Metadata.ID] = prompt
}

// GetPrompt retrieves a prompt by ID
func (pl *PromptLibrary) GetPrompt(id string) (*Prompt, bool) {
	prompt, exists := pl.Prompts[id]
	return prompt, exists
}

// ListPrompts returns all prompts
func (pl *PromptLibrary) ListPrompts() []*Prompt {
	prompts := make([]*Prompt, 0, len(pl.Prompts))
	for _, prompt := range pl.Prompts {
		prompts = append(prompts, prompt)
	}
	return prompts
}

// GetPromptsByTag returns prompts that have the specified tag
func (pl *PromptLibrary) GetPromptsByTag(tag string) []*Prompt {
	var prompts []*Prompt
	for _, prompt := range pl.Prompts {
		for _, promptTag := range prompt.Metadata.Tags {
			if promptTag == tag {
				prompts = append(prompts, prompt)
				break
			}
		}
	}
	return prompts
}