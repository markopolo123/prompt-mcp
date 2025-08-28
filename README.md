# Prompt MCP Server

A Golang MCP (Model Context Protocol) server that enables teams to share and manage reusable prompts for Claude Code through a GitOps workflow.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Prompt Management](#prompt-management)
- [Development](#development)
- [Testing](#testing)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Features

- **File-based Storage**: Store prompts in human-readable YAML format
- **GitOps Workflow**: Manage prompts through version control
- **Rich Metadata**: Support for versioning, tagging, and author attribution
- **Argument Substitution**: Dynamic prompts with configurable parameters
- **Category Organization**: Organize prompts by project, team, or use case
- **MCP Integration**: Seamless integration with Claude Code
- **Hot Reloading**: Automatic detection of prompt changes (planned)
- **Usage Statistics**: Track prompt usage and performance (planned)

## Installation

### Prerequisites

- Go 1.24.1 or later
- [just](https://github.com/casey/just) command runner (recommended)

### From Source

```bash
# Clone the repository
git clone https://github.com/markopolo123/prompt-mcp.git
cd prompt-mcp

# Install dependencies
go mod tidy

# Build the server
just build
# or
go build -o bin/prompt-mcp ./cmd/server
```

### Using Go Install

```bash
go install github.com/markopolo123/prompt-mcp/cmd/server@latest
```

## Usage

### Basic Usage

Start the server with default settings:

```bash
./bin/prompt-mcp
```

### Command Line Options

```bash
./bin/prompt-mcp [options]

Options:
  -prompts-dir string
        Directory containing prompt files (default "./prompts")
  -name string
        Server name (default "team-prompt-server")
  -ver string
        Server version (default "1.0.0")
  -version
        Print version and exit
```

### Integration with Claude Code

1. Start the prompt server:
   ```bash
   ./bin/prompt-mcp -prompts-dir /path/to/your/prompts
   ```

2. Configure Claude Code to connect to the MCP server following the MCP documentation.

3. Browse and use shared prompts directly within Claude Code.

## Configuration

### Server Configuration

Create a `config/server.yaml` file for advanced configuration:

```yaml
server:
  name: "team-prompt-server"
  version: "1.0.0"
  
storage:
  prompts_dir: "./prompts"
  watch_changes: false  # future feature
  
mcp:
  capabilities:
    prompts: true
    resources: false
    tools: false
```

### Prompt Structure

Prompts are stored as YAML files in the `prompts/` directory:

```yaml
# prompts/category/prompt-name.yaml
metadata:
  id: "unique-prompt-id"
  name: "Prompt Display Name"
  description: "Detailed description of prompt purpose"
  author: "author-name"
  created: "2025-08-27T10:00:00Z"
  modified: "2025-08-27T10:00:00Z"
  version: "1.0.0"
  tags:
    - "category"
    - "project-name"
    - "use-case"
  
arguments:
  - name: "variable_name"
    description: "What this variable represents"
    type: "string"  # string, number, boolean
    required: true
    default: "default_value"

prompt: |
  Your multi-line prompt content here.
  Use {{variable_name}} for argument substitution.
  
usage_stats:
  usage_count: 0
  last_used: "2025-08-27T10:00:00Z"
```

## API Documentation

### MCP Protocol Implementation

The server implements the Model Context Protocol (MCP) specification:

#### List Prompts
- **Operation**: List all available prompts
- **Returns**: Prompt metadata including ID, name, description, and tags

#### Get Prompt
- **Operation**: Retrieve a specific prompt with argument resolution
- **Input**: Prompt ID and argument values
- **Returns**: Rendered prompt content with substituted variables

#### Resource URIs
- **Format**: `prompt://category/prompt-name`
- **Usage**: Reference prompts by their file path structure

## Prompt Management

### Adding New Prompts

1. Create a YAML file in the appropriate category directory:
   ```bash
   mkdir -p prompts/your-category
   touch prompts/your-category/your-prompt.yaml
   ```

2. Follow the prompt schema structure (see Configuration section)

3. Commit to version control for team sharing

### Directory Structure

```
prompts/
├── development/
│   ├── code-review.yaml
│   └── bug-analysis.yaml
├── documentation/
│   ├── api-docs.yaml
│   └── readme-gen.yaml
└── testing/
    ├── test-gen.yaml
    └── test-review.yaml
```

### Best Practices

- Use descriptive prompt IDs and filenames
- Include comprehensive metadata
- Add relevant tags for discoverability
- Provide clear argument descriptions
- Test prompts before committing

## Development

### Development Setup

```bash
# Install development dependencies
just install-dev

# Format code
just fmt

# Run linter
just lint

# Build and run
just run
```

### Available Commands

- `just build` - Build the server binary
- `just run` - Build and run the server
- `just test` - Run all tests
- `just test-verbose` - Run tests with verbose output
- `just fmt` - Format Go code
- `just lint` - Run golangci-lint
- `just tidy` - Tidy Go modules
- `just clean` - Remove build artifacts
- `just dev` - Run development server with hot reload (requires air)

### Project Structure

```
/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── server/         # MCP server implementation
│   ├── prompt/         # Prompt models and validation
│   └── storage/        # File system storage layer
├── prompts/            # Example prompts
├── config/             # Configuration files
├── examples/           # Example configurations
└── bin/               # Built binaries
```

## Testing

### Running Tests

```bash
# Run all tests
just test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/prompt
```

### Test Coverage

The project includes:
- Unit tests for prompt loading and validation
- Integration tests for MCP server functionality
- Example prompts for testing various scenarios

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following the existing code style
4. Add tests for new functionality
5. Run tests and linting (`just test && just lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Write tests for new functionality
- Update documentation for user-facing changes
- Use semantic versioning for releases
- Keep prompts simple and focused

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

### Getting Help

- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/markopolo123/prompt-mcp/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/markopolo123/prompt-mcp/discussions)
- **Documentation**: Check the [ADR documentation](adr/) for architectural decisions

### Troubleshooting

#### Server Won't Start
- Verify the prompts directory exists and is readable
- Check for YAML syntax errors in prompt files
- Ensure Go version compatibility (1.24.1+)

#### Prompts Not Loading
- Validate YAML structure against the schema
- Check file permissions on prompt files
- Review server logs for specific error messages

#### MCP Connection Issues
- Verify Claude Code MCP configuration
- Check server is running on expected port
- Review MCP protocol compatibility

---