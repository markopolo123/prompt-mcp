package server

import (
	"context"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/markopolo123/prompt-mcp/internal/prompt"
	"github.com/markopolo123/prompt-mcp/internal/storage"
)

// Server represents the MCP server
type Server struct {
	mcpServer *server.MCPServer
	storage   *storage.FileSystemStorage
	library   *prompt.PromptLibrary
	config    Config
}

// Config holds server configuration
type Config struct {
	Name         string
	Version      string
	PromptsDir   string
	WatchChanges bool
}

// NewServer creates a new MCP server
func NewServer(config Config) (*Server, error) {
	// Initialize storage
	storageConfig := storage.Config{
		PromptsDir:   config.PromptsDir,
		WatchChanges: config.WatchChanges,
	}
	
	storage, err := storage.NewFileSystemStorage(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Create server instance
	srv := &Server{
		storage: storage,
		config:  config,
	}

	// Create MCP server with prompt capabilities
	mcpServer := server.NewMCPServer(config.Name, config.Version, 
		server.WithPromptCapabilities(true),
	)

	srv.mcpServer = mcpServer

	return srv, nil
}

// LoadPrompts loads all prompts from storage
func (s *Server) LoadPrompts() error {
	library, err := s.storage.LoadLibrary()
	if err != nil {
		return fmt.Errorf("failed to load prompt library: %w", err)
	}

	s.library = library
	s.registerPrompts()
	
	log.Printf("Loaded %d prompts", len(library.Prompts))
	return nil
}

// registerPrompts registers all loaded prompts with the MCP server
func (s *Server) registerPrompts() {
	for _, p := range s.library.ListPrompts() {
		// Create prompt options
		options := []mcp.PromptOption{
			mcp.WithPromptDescription(p.Metadata.Description),
		}
		
		// Add arguments
		for _, arg := range p.Arguments {
			argOptions := []mcp.ArgumentOption{
				mcp.ArgumentDescription(arg.Description),
			}
			if arg.Required {
				argOptions = append(argOptions, mcp.RequiredArgument())
			}
			options = append(options, mcp.WithArgument(arg.Name, argOptions...))
		}

		// Create prompt
		promptDef := mcp.NewPrompt(p.Metadata.ID, options...)

		// Register the prompt with its handler
		s.mcpServer.AddPrompt(promptDef, s.createPromptHandler(p))
	}
}

// createPromptHandler creates a handler function for a specific prompt
func (s *Server) createPromptHandler(p *prompt.Prompt) func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		// Debug: Log incoming arguments
		log.Printf("Debug: Prompt '%s' received arguments: %+v", p.Metadata.ID, request.Params.Arguments)
		
		// Convert string arguments to interface{} for processing
		args := make(map[string]interface{})
		for key, value := range request.Params.Arguments {
			log.Printf("Debug: Argument '%s' = '%v' (type: %T)", key, value, value)
			args[key] = value
		}
		
		// Resolve arguments and substitute in prompt content
		resolvedContent, err := s.resolvePromptContent(p, args)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve prompt content: %w", err)
		}

		// Update usage statistics
		s.updateUsageStats(p)

		// Return the resolved prompt
		return mcp.NewGetPromptResult(
			p.Metadata.Name,
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent(resolvedContent),
				),
			},
		), nil
	}
}

// Start starts the MCP server
func (s *Server) Start(ctx context.Context) error {
	// Load prompts before starting
	if err := s.LoadPrompts(); err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	log.Printf("Starting %s v%s", s.config.Name, s.config.Version)
	log.Printf("Loaded prompts from: %s", s.config.PromptsDir)
	
	// Run the MCP server using stdio transport
	return server.ServeStdio(s.mcpServer)
}

// GetLibrary returns the current prompt library
func (s *Server) GetLibrary() *prompt.PromptLibrary {
	return s.library
}

// Reload reloads prompts from storage
func (s *Server) Reload() error {
	return s.LoadPrompts()
}