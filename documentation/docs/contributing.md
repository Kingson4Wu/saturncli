---
sidebar_position: 9
title: Contributing to Saturn CLI - Development & Contribution Guide
description: Guidelines and information to help you contribute effectively to Saturn CLI. Learn about development setup, coding standards, and contribution process.
keywords: [saturn cli contributing, contribute to go project, saturn cli development, go cli contribution guide, open source contribution]
---

# Contributing

Thank you for your interest in contributing to Saturn CLI! This document provides guidelines and information to help you contribute effectively to the project and improve Go job execution capabilities.

## Code of Conduct

Please read and follow our [Code of Conduct](https://github.com/Kingson4Wu/saturncli/blob/main/CODE_OF_CONDUCT.md) to help create a positive and inclusive community.

## How to Contribute

### Reporting Issues

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Check the documentation** to see if your question is answered
3. **Verify it's a Saturn CLI issue** and not an environment problem

When reporting issues, please include:
- Go version (`go version`)
- Operating system and version
- Saturn CLI version
- Steps to reproduce the issue
- Expected vs actual behavior
- Any relevant error messages or logs

### Feature Requests

We welcome feature requests! Before submitting:
- Check if the feature already exists
- Consider if it fits the project scope
- Provide a clear use case and potential implementation approach

### Pull Requests

#### Preparing Your Pull Request

1. **Fork the repository** and clone your fork
2. **Create a feature branch** from the `main` branch
3. **Add tests** if your change affects functionality
4. **Follow the code style** used throughout the project
5. **Update documentation** if needed
6. **Run tests** to ensure all functionality still works

#### Setting Up Your Development Environment

1. Install Go 1.19 or later
2. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/saturncli.git
   cd saturncli
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run tests to verify your setup:
   ```bash
   go test ./...
   ```

#### Code Standards

- **Go Formatting**: Use `gofmt` or `go fmt` for consistent formatting
- **Naming**: Follow Go naming conventions and idioms
- **Comments**: Document exported functions, types, and packages
- **Error Handling**: Follow Go error handling patterns
- **Testing**: Include tests for new functionality and bug fixes

Example of proper function documentation:

```go
// AddJob registers a regular (non-stoppable) job with the default registry.
// It returns an error if the job name is already registered.
func AddJob(name string, handler JobHandler) error {
    // implementation
}
```

#### Testing Guidelines

All contributions should include appropriate tests:

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test the interaction between components
- **Edge Cases**: Test error conditions and boundary values

Run all tests before submitting:
```bash
go test ./...
```

For more comprehensive testing:
```bash
go test -v ./...  # Verbose output
go test -race ./...  # Race condition detection
```

## Development Workflow

### Making Changes

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/descriptive-name
   ```

2. **Make your changes** following the project's coding standards

3. **Write or update tests** as needed

4. **Run the test suite**:
   ```bash
   go test ./...
   ```

5. **Check code quality**:
   ```bash
   # You may need to install these tools first
   go vet ./...
   golangci-lint run  # if available
   ```

### Commit Guidelines

- **Use clear, descriptive commit messages**
- **Follow conventional commits** format when possible:
  - `feat:` for new features
  - `fix:` for bug fixes
  - `docs:` for documentation changes
  - `test:` for test additions
  - `refactor:` for code restructuring
  - `perf:` for performance improvements
  - `chore:` for maintenance tasks

Example commit message:
```
feat: add support for custom logger interface

- Add Logger interface definition
- Update server and client to use interface
- Provide default logger implementation
- Update tests to use mock logger
```

### Submitting Pull Requests

1. **Squash commits** if you have multiple small changes
2. **Update your branch** with the latest from main:
   ```bash
   git remote add upstream https://github.com/Kingson4Wu/saturncli.git
   git fetch upstream
   git rebase upstream/main
   ```
3. **Push your branch** to your fork
4. **Create a pull request** from your fork to the main repository

## Project Structure

Understanding the project layout will help you navigate and contribute:

```
saturncli/
├── base/                 # Basic types and constants
├── client/               # Client-side implementation
│   ├── client.go         # Main client implementation
│   ├── cmd.go            # Command-line interface
│   └── client_windows.go # Windows-specific client code
├── server/               # Server-side implementation
│   ├── server.go         # Main server implementation
│   ├── server_windows.go # Windows-specific server code
│   └── job_manager.go    # Job registration and management
├── utils/                # Utility functions
├── examples/             # Example implementations
│   ├── client/
│   └── server/
├── resource/             # Static resources (images, etc.)
├── documentation/        # Documentation files
└── Makefile              # Build automation
```

### Key Components

- **base**: Contains shared types and constants used by both client and server
- **client**: Implements client-side functionality for connecting to Saturn servers
- **server**: Implements server-side functionality for handling job requests  
- **utils**: Reusable utility functions like logging
- **examples**: Sample implementations demonstrating Saturn CLI usage

## Areas to Contribute

### Documentation Improvements

- API documentation
- Usage examples
- Architecture documentation
- Tutorials
- Migration guides

### Feature Development

- New transport mechanisms
- Enhanced security features
- Performance improvements
- Additional job types
- Monitoring and metrics

### Bug Fixes

- Resolving reported issues
- Improving error handling
- Enhancing stability and reliability

### Testing

- Increasing test coverage
- Adding integration tests
- Improving test infrastructure

## Review Process

### What We Look For

When reviewing pull requests, we consider:

- **Correctness**: Does the code work as intended?
- **Style**: Does it follow Go idioms and project conventions?
- **Test Coverage**: Are there adequate tests?
- **Documentation**: Is the change properly documented?
- **Performance**: Does it have reasonable performance characteristics?
- **Compatibility**: Does it maintain backward compatibility where appropriate?

### Review Timeline

- Initial review: Typically within 1-3 business days
- Iteration: Usually within 1-2 days of changes
- Final review: Once all concerns are addressed

### Common Feedback

Reviewers often comment on:

- Code organization and structure
- Error handling patterns
- Test coverage
- Documentation completeness
- Performance considerations
- API design choices

## Community

### Getting Help

For questions about contributing:
- Open an issue in the repository
- Check the [documentation](https://github.com/Kingson4Wu/saturncli/wiki)

### Staying Updated

- Watch the repository for releases
- Check the [CHANGELOG](https://github.com/Kingson4Wu/saturncli/blob/main/CHANGELOG.md) for updates

## Special Thanks

Contributors help make Saturn CLI better for everyone! Your efforts are greatly appreciated, whether it's:

- Reporting bugs
- Suggesting features
- Improving documentation
- Writing code
- Answering questions

## Questions?

If you have questions about contributing that aren't covered in this guide, please open an issue and we'll be happy to help!

## Next Steps

Ready to contribute? Start by:

1. Looking at issues tagged with `good first issue`
2. Setting up your development environment
3. Joining discussions in the issue tracker
4. Submitting your first contribution!

---

*This guide was inspired by various open source contribution guides and adapted for the Saturn CLI project. If you find issues with this guide, please feel free to suggest improvements!*