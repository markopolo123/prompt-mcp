package server

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/markopolo123/prompt-mcp/internal/prompt"
)

// resolvePromptContent resolves arguments in prompt content
func (s *Server) resolvePromptContent(promptObj *prompt.Prompt, args map[string]interface{}) (string, error) {
	content := promptObj.Prompt
	
	// Create a map of argument values, including defaults
	argValues := make(map[string]interface{})
	
	// First, set defaults
	for _, arg := range promptObj.Arguments {
		if arg.Default != nil {
			argValues[arg.Name] = arg.Default
		}
	}
	
	// Then, override with provided values
	if args != nil {
		for key, value := range args {
			argValues[key] = value
		}
	}
	
	// Validate required arguments
	for _, arg := range promptObj.Arguments {
		if arg.Required {
			if _, exists := argValues[arg.Name]; !exists {
				return "", fmt.Errorf("required argument '%s' not provided", arg.Name)
			}
		}
	}
	
	// Validate argument types and convert values
	for _, arg := range promptObj.Arguments {
		if value, exists := argValues[arg.Name]; exists {
			log.Printf("Debug: Converting argument '%s' (type: %s) with value '%v' (%T)", 
				arg.Name, arg.Type, value, value)
			convertedValue, err := s.convertArgumentValue(value, arg.Type)
			if err != nil {
				log.Printf("Debug: Failed to convert argument '%s': %v", arg.Name, err)
				return "", fmt.Errorf("argument '%s': %w", arg.Name, err)
			}
			argValues[arg.Name] = convertedValue
			log.Printf("Debug: Converted argument '%s' to '%v' (%T)", 
				arg.Name, convertedValue, convertedValue)
		}
	}
	
	// Replace placeholders in content
	variablePattern := regexp.MustCompile(`\{\{(\w+)\}\}`)
	result := variablePattern.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name
		varName := strings.Trim(match, "{}")
		
		if value, exists := argValues[varName]; exists {
			return fmt.Sprintf("%v", value)
		}
		
		// Return original if not found (shouldn't happen due to validation)
		return match
	})
	
	return result, nil
}

// convertArgumentValue converts and validates argument values based on their type
func (s *Server) convertArgumentValue(value interface{}, argType prompt.ArgumentType) (interface{}, error) {
	switch argType {
	case prompt.ArgumentTypeString:
		return fmt.Sprintf("%v", value), nil
		
	case prompt.ArgumentTypeNumber:
		switch v := value.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		case string:
			if parsed, err := strconv.ParseFloat(v, 64); err == nil {
				return parsed, nil
			}
			return nil, fmt.Errorf("cannot convert '%s' to number", v)
		default:
			return nil, fmt.Errorf("invalid number value: %v", value)
		}
		
	case prompt.ArgumentTypeBoolean:
		switch v := value.(type) {
		case bool:
			return v, nil
		case string:
			// Be more flexible with boolean parsing
			lowerV := strings.ToLower(strings.TrimSpace(v))
			switch lowerV {
			case "true", "yes", "1", "on":
				return true, nil
			case "false", "no", "0", "off", "":
				return false, nil
			default:
				// If it's not a recognizable boolean, treat as false and warn
				log.Printf("Warning: Unable to parse '%s' as boolean, defaulting to false", v)
				return false, nil
			}
		default:
			return nil, fmt.Errorf("invalid boolean value: %v", value)
		}
		
	default:
		return nil, fmt.Errorf("unsupported argument type: %s", argType)
	}
}

// updateUsageStats updates the usage statistics for a prompt
func (s *Server) updateUsageStats(promptObj *prompt.Prompt) {
	// Note: In a production system, you might want to persist these stats
	// For now, we just update the in-memory stats
	promptObj.UsageStats.UsageCount++
	// promptObj.UsageStats.LastUsed = time.Now() // Uncomment when needed
}