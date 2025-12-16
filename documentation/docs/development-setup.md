---
sidebar_position: 11
title: Development Setup - Saturn CLI Environment Configuration
description: Complete guide to setting up your development environment for Saturn CLI. Learn how to build, test, and contribute to the Go job execution toolkit.
keywords: [saturn cli development setup, go cli development, saturn cli build, development environment setup, go job execution development]
---

# Development Setup

This guide covers setting up your development environment for working with Saturn CLI, whether you're using it in your projects or contributing to the project itself for Go job execution.

## Prerequisites

### System Requirements

- **Operating System**: macOS, Linux, or Windows
- **Go Version**: 1.19 or higher
- **Git**: For version control and dependency management
- **Make** (optional): For build automation (available on Unix-like systems, can be installed on Windows)

### Verify Prerequisites

Before starting, verify your system meets the requirements:

```bash
# Check Go installation
go version
# Should show: go version go1.19.x or higher

# Check Git installation
git --version
# Should show Git version

# Check environment variables
echo $GOPATH
echo $GOROOT

# Verify Go environment
go env GOOS GOARCH
```

## Setting Up Your Development Environment

### 1. Fork and Clone the Repository

If you plan to contribute to Saturn CLI:

```bash
# Fork the repository on GitHub (click Fork button on https://github.com/Kingson4Wu/saturncli)

# Clone your fork
git clone https://github.com/YOUR_USERNAME/saturncli.git
cd saturncli

# Add upstream remote for sync
git remote add upstream https://github.com/Kingson4Wu/saturncli.git
```

### 2. Initialize Go Modules

Saturn CLI uses Go modules for dependency management:

```bash
# Ensure you're in the project root
ls go.mod go.sum  # Should list these files

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```

### 3. Verify the Setup

Test that you can build the project:

```bash
# Build the examples (this creates the binaries)
make

# Or build manually
go build -o saturn_svr ./examples/server/server.go
go build -o saturn_cli ./examples/client/client.go

# Run basic tests
go test ./...
```

## Development Environment Configuration

### IDE Setup

Configure your IDE for optimal Go development:

#### VS Code
1. Install the Go extension
2. Configure settings in `.vscode/settings.json`:

```json
{
  "go.useLanguageServer": true,
  "go.formatTool": "gofumpt",
  "go.lintTool": "golangci-lint",
  "go.buildFlags": [],
  "go.testFlags": ["-v"]
}
```

#### GoLand/IntelliJ
1. Install Go plugin
2. Open project and configure Go SDK
3. Enable gofmt/goimports on save

### Linting and Formatting

Install development tools:

```bash
# Install Go linters
go install golang.org/x/tools/cmd/goimports@latest
go install mvdan.cc/gofumpt@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linters
golangci-lint run

# Format code
gofumpt -w .
goimports -w .
```

## Project Structure

Understanding the Saturn CLI codebase:

```
saturncli/
├── base/                 # Shared constants and types
│   └── result.go         # Result type definitions
├── client/               # Client-side implementation
│   ├── client.go         # Core client logic
│   ├── client_windows.go # Windows-specific client code
│   ├── cmd.go           # Command-line interface
│   └── client_manager.go # Client management utilities
├── server/               # Server-side implementation
│   ├── server.go        # Core server logic
│   ├── server_windows.go # Windows-specific server code
│   ├── job_manager.go   # Job registry and management
│   └── job_manager_test.go # Job manager tests
├── utils/               # Utility functions
│   └── logger.go        # Logging utilities
├── examples/            # Example implementations
│   ├── client/
│   └── server/
├── resource/            # Static resources
├── documentation/       # Documentation files
├── .githooks/          # Git hooks
├── go.mod, go.sum      # Go module files
└── Makefile            # Build automation
```

### Key Files and Their Purposes

| File | Purpose |
|------|---------|
| `client/cmd.go` | Command-line interface parser and runner |
| `client/client.go` | Core client communication logic |
| `server/server.go` | Core server implementation (Unix) |
| `server/server_windows.go` | Server implementation for Windows |
| `server/job_manager.go` | Job registration and execution management |
| `base/result.go` | Success/failure result types |
| `utils/logger.go` | Logging interface and implementations |

## Building and Testing

### Build Process

```bash
# Full build using Make
make

# Build server binary
go build -o saturn_svr ./examples/server/server.go

# Build client binary
go build -o saturn_cli ./examples/client/client.go

# Build with specific flags
go build -ldflags="-s -w" -o saturn_svr ./examples/server/server.go
```

### Testing Strategy

Saturn CLI uses Go's built-in testing framework:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...

# Run specific test package
go test ./server/...

# Run specific test
go test -run TestServerStartup ./server/...

# Get test coverage
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

Tests are organized alongside the code they test:

```go
// server/job_manager_test.go
package server

import (
    "testing"
    "time"
)

func TestAddJob(t *testing.T) {
    registry := NewRegistry()
    
    err := registry.AddJob("test", func(args map[string]string, signature string) bool {
        return true
    })
    
    if err != nil {
        t.Errorf("Expected no error, got: %v", err)
    }
    
    // Additional test assertions
}

func TestStoppableJob(t *testing.T) {
    // Test stoppable job functionality
    quit := make(chan struct{})
    done := make(chan bool, 1)
    
    handler := func(args map[string]string, signature string, quit chan struct{}) bool {
        // Implementation to test
        select {
        case <-quit:
            return true
        case <-time.After(100 * time.Millisecond):
            return true
        }
    }
    
    // Test the handler...
}
```

## Development Workflow

### 1. Creating a Development Branch

```bash
# Update your main branch
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/descriptive-name

# Or for bug fixes
git checkout -b fix/issue-description
```

### 2. Making Changes

```bash
# Make your changes to the code
# Add new functionality or fix bugs

# Format your code
gofumpt -w .
goimports -w .

# Run linters
golangci-lint run

# Run tests
go test ./...
```

### 3. Before Committing

```bash
# Ensure all tests pass
go test -race ./...

# Verify build
make

# Update documentation if needed
# (Check documentation guidelines)

# Commit with conventional commits format
git add .
git commit -m "feat: add new job type registration"

# Or for fixes
git commit -m "fix: resolve connection timeout issue"
```

### 4. Pushing Changes

```bash
# Push your branch
git push origin feature/descriptive-name

# Create a pull request on GitHub
```

## Debugging Saturn CLI

### Debugging Server Issues

```go
// Add debug logging to your server
func main() {
    logger := &utils.DebugLogger{} // Use a debug logger
    
    if err := server.AddJob("debug-job", func(args map[string]string, signature string) bool {
        logger.Debugf("Processing job %s with args: %+v", signature, args)
        
        // Your job logic here
        result := doWork(args)
        
        logger.Debugf("Job %s completed with result: %t", signature, result)
        return result
    }); err != nil {
        logger.Errorf("Failed to register job: %v", err)
    }
    
    server.NewServer(logger, "/tmp/debug.sock").Serve()
}
```

### Using Delve for Debugging

```bash
# Install Delve debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug server
dlv debug ./examples/server/server.go -- /tmp/debug.sock

# Debug with arguments
dlv exec ./saturn_svr -- /tmp/debug.sock

# Set breakpoints
(dlv) break server.go:45
(dlv) continue
```

### Debugging Client Issues

```go
// Create a debug version of your client
func main() {
    logger := &utils.DebugLogger{}
    
    cli := client.NewClient(logger, "/tmp/debug.sock")
    
    result := cli.Run(&client.Task{
        Name:   "debug-job",
        Params: map[string]string{"debug": "true"},
    })
    
    logger.Infof("Job result: %v", result)
}
```

## Running Examples and Demos

### Server Examples

```bash
# Start the example server
go run ./examples/server/server.go

# In another terminal, run client commands
./saturn_cli --name hello --param id=33
./saturn_cli --name hello_stoppable
./saturn_cli --name hello_stoppable --stop
```

### Custom Examples

Create your own test examples:

```go
// custom_example.go
package main

import (
    "log"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Register various test jobs
    server.AddJob("fast-job", func(args map[string]string, signature string) bool {
        log.Printf("[%s] Fast job executed with args: %+v", signature, args)
        return true
    })
    
    server.AddStoppableJob("slow-job", func(args map[string]string, signature string, quit chan struct{}) bool {
        log.Printf("[%s] Slow job started", signature)
        
        for i := 0; i < 100; i++ {
            select {
            case <-quit:
                log.Printf("[%s] Slow job stopped at iteration %d", signature, i)
                return true
            default:
                log.Printf("[%s] Slow job working... step %d", signature, i)
                time.Sleep(500 * time.Millisecond)
            }
        }
        
        log.Printf("[%s] Slow job completed", signature)
        return true
    })
    
    log.Println("Server starting on /tmp/custom.sock")
    server.NewServer(&utils.DefaultLogger{}, "/tmp/custom.sock").Serve()
}
```

## Performance Testing

### Benchmarking

Create benchmarks for performance-critical code:

```go
// benchmark_test.go
package main

import (
    "testing"
    "time"
)

func BenchmarkJobExecution(b *testing.B) {
    // Setup code that doesn't count toward benchmark time
    job := func(args map[string]string, signature string) bool {
        time.Sleep(1 * time.Millisecond) // Simulate work
        return true
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = job(map[string]string{}, "benchmark")
    }
}

func BenchmarkConcurrentJobs(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Test concurrent execution
        }
    })
}
```

### Memory Profiling

```bash
# Run with memory profiling
go test -memprofile=mem.prof -memprofilerate=1 ./...
go tool pprof mem.prof

# Run server with profiling enabled
go build -o saturn_svr_with_profiling ./examples/server/server.go
./saturn_svr_with_profiling
```

## Continuous Integration Setup

If contributing, understand the CI setup:

```yaml
# .github/workflows/test.yml (example)
name: Tests
on: [push, pull_request]
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.19.x, 1.20.x, 1.21.x]
        os: [ubuntu-latest, macos-latest, windows-latest]
    runs-on: ${{ matrix.os }}
    steps:
    - uses: actions/checkout@v3
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
    - name: Test
      run: go test -v ./...
    - name: Lint
      run: |
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        golangci-lint run
```

## Troubleshooting Development Setup

### Common Issues

#### Module Issues
```bash
# If getting "package not found" errors
go mod tidy

# If dependencies are out of sync
go clean -modcache
go mod download
```

#### Build Issues
```bash
# Clean build cache
go clean -cache

# Force rebuild
go build -a -o saturn_svr ./examples/server/server.go
```

#### Test Issues
```bash
# If tests are failing due to timeouts
go test -timeout=30s ./...

# Run tests in verbose mode to see what's failing
go test -v -timeout=60s ./...
```

## Development Best Practices

### Code Standards

1. **Format Consistently**: Use `gofumpt` and `goimports`
2. **Write Tests**: Add unit tests for new functionality
3. **Document Public APIs**: Use proper Go doc comments
4. **Handle Errors Properly**: Don't ignore errors
5. **Use Context**: For timeout and cancellation propagation

### Performance Considerations

```go
// Good: Efficient resource management
func efficientJob(args map[string]string, signature string) bool {
    // Use buffered I/O
    reader := bufio.NewReaderSize(inputFile, 64*1024)
    
    // Pre-allocate slices when size is known
    results := make([]string, 0, expectedSize)
    
    // Minimize allocations in loops
    var buffer strings.Builder
    return true
}

// Good: Proper error handling
func safeJob(args map[string]string, signature string) bool {
    resource, err := acquireResource(args["id"])
    if err != nil {
        return false // Proper error handling
    }
    defer resource.Close()
    
    return doWork(resource)
}
```

## Next Steps

- Review the [Contributing](./contributing.md) guidelines for contribution processes
- Check the [Architecture](./architecture.md) documentation to understand the system design
- Look at the [Examples](./examples.md) for practical implementations
- Follow the [Best Practices](./best-practices.md) guide for quality implementation