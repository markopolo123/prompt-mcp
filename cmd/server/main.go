package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/markopolo123/prompt-mcp/internal/server"
)

func main() {
	// Command line flags
	var (
		promptsDir = flag.String("prompts-dir", "./prompts", "Directory containing prompt files")
		version    = flag.Bool("version", false, "Print version and exit")
		name       = flag.String("name", "team-prompt-server", "Server name")
		ver        = flag.String("ver", "1.0.0", "Server version")
	)
	flag.Parse()

	if *version {
		log.Printf("%s v%s", *name, *ver)
		os.Exit(0)
	}

	// Ensure prompts directory exists and is absolute
	absPromptsDir, err := filepath.Abs(*promptsDir)
	if err != nil {
		log.Fatalf("Failed to resolve prompts directory path: %v", err)
	}

	if _, err := os.Stat(absPromptsDir); os.IsNotExist(err) {
		log.Fatalf("Prompts directory does not exist: %s", absPromptsDir)
	}

	// Create server configuration
	config := server.Config{
		Name:         *name,
		Version:      *ver,
		PromptsDir:   absPromptsDir,
		WatchChanges: false, // TODO: Implement file watching
	}

	// Create and initialize server
	srv, err := server.NewServer(config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		cancel()
	}()

	// Start the server
	log.Printf("Starting server with prompts directory: %s", absPromptsDir)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	
	log.Println("Server stopped")
}