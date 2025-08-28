package prompt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidatePrompt validates a prompt structure
func ValidatePrompt(prompt *Prompt) error {
	if err := validateMetadata(&prompt.Metadata); err != nil {
		return fmt.Errorf("metadata validation failed: %w", err)
	}

	if err := validateArguments(prompt.Arguments); err != nil {
		return fmt.Errorf("arguments validation failed: %w", err)
	}

	if err := validatePromptContent(prompt.Prompt, prompt.Arguments); err != nil {
		return fmt.Errorf("prompt content validation failed: %w", err)
	}

	return nil
}

// validateMetadata validates prompt metadata
func validateMetadata(metadata *Metadata) error {
	if strings.TrimSpace(metadata.ID) == "" {
		return errors.New("id is required")
	}

	if !isValidID(metadata.ID) {
		return errors.New("id must contain only alphanumeric characters, hyphens, and underscores")
	}

	if strings.TrimSpace(metadata.Name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(metadata.Description) == "" {
		return errors.New("description is required")
	}

	if strings.TrimSpace(metadata.Author) == "" {
		return errors.New("author is required")
	}

	if strings.TrimSpace(metadata.Version) == "" {
		return errors.New("version is required")
	}

	if metadata.Created.IsZero() {
		return errors.New("created timestamp is required")
	}

	if metadata.Modified.IsZero() {
		return errors.New("modified timestamp is required")
	}

	return nil
}

// validateArguments validates prompt arguments
func validateArguments(arguments []Argument) error {
	namesSeen := make(map[string]bool)

	for i, arg := range arguments {
		if strings.TrimSpace(arg.Name) == "" {
			return fmt.Errorf("argument %d: name is required", i)
		}

		if !isValidArgumentName(arg.Name) {
			return fmt.Errorf("argument %d: name must contain only alphanumeric characters and underscores", i)
		}

		if namesSeen[arg.Name] {
			return fmt.Errorf("argument %d: duplicate name '%s'", i, arg.Name)
		}
		namesSeen[arg.Name] = true

		if strings.TrimSpace(arg.Description) == "" {
			return fmt.Errorf("argument %d (%s): description is required", i, arg.Name)
		}

		if !isValidArgumentType(arg.Type) {
			return fmt.Errorf("argument %d (%s): invalid type '%s'", i, arg.Name, arg.Type)
		}

		if arg.Default != nil {
			if err := validateArgumentDefault(arg.Default, arg.Type); err != nil {
				return fmt.Errorf("argument %d (%s): %w", i, arg.Name, err)
			}
		}
	}

	return nil
}

// validatePromptContent validates the prompt content against its arguments
func validatePromptContent(content string, arguments []Argument) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("prompt content is required")
	}

	// Extract variable placeholders from the prompt
	variablePattern := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := variablePattern.FindAllStringSubmatch(content, -1)

	usedVariables := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			usedVariables[match[1]] = true
		}
	}

	// Create a map of defined arguments
	definedArgs := make(map[string]bool)
	for _, arg := range arguments {
		definedArgs[arg.Name] = true
	}

	// Check for undefined variables in prompt
	for variable := range usedVariables {
		if !definedArgs[variable] {
			return fmt.Errorf("undefined variable '%s' used in prompt", variable)
		}
	}

	// Check for unused required arguments
	for _, arg := range arguments {
		if arg.Required && !usedVariables[arg.Name] {
			return fmt.Errorf("required argument '%s' is not used in prompt", arg.Name)
		}
	}

	return nil
}

// isValidID checks if an ID is valid (alphanumeric, hyphens, underscores)
func isValidID(id string) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return pattern.MatchString(id)
}

// isValidArgumentName checks if an argument name is valid
func isValidArgumentName(name string) bool {
	pattern := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	return pattern.MatchString(name)
}

// isValidArgumentType checks if an argument type is valid
func isValidArgumentType(argType ArgumentType) bool {
	switch argType {
	case ArgumentTypeString, ArgumentTypeNumber, ArgumentTypeBoolean:
		return true
	default:
		return false
	}
}

// validateArgumentDefault validates that a default value matches the argument type
func validateArgumentDefault(defaultValue interface{}, argType ArgumentType) error {
	switch argType {
	case ArgumentTypeString:
		if _, ok := defaultValue.(string); !ok {
			return errors.New("default value must be a string")
		}
	case ArgumentTypeNumber:
		switch defaultValue.(type) {
		case int, int64, float32, float64:
			return nil
		default:
			return errors.New("default value must be a number")
		}
	case ArgumentTypeBoolean:
		if _, ok := defaultValue.(bool); !ok {
			return errors.New("default value must be a boolean")
		}
	}
	return nil
}