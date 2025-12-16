---
sidebar_position: 1
---

# Setup Guide

This guide covers the complete setup process for Saturn CLI, including prerequisites, installation, and initial configuration.

## Prerequisites

- Go 1.19 or newer
- Unix-like system for socket transport (macOS/Linux) with Unix domain socket support
- On Windows systems, TCP transport is used instead
- Git for cloning the repository (if installing from source)

## Installation Methods

You have several options for installing Saturn CLI:

### Method 1: Go Install (Recommended for using as library)

If you want to use Saturn CLI as a library in your Go projects:

```bash
go mod init your-project-name
go get github.com/Kingson4Wu/saturncli
```

### Method 2: Build from Source

To build both client and server binaries:

```bash
git clone https://github.com/Kingson4Wu/saturncli.git
cd saturncli
make
```

This produces:
- `saturn_svr` – Reference server demonstrating job registration
- `saturn_cli` – Command-line client

Alternative build command if Make is not available:
```bash
go build -o saturn_svr ./examples/server/server.go
go build -o saturn_cli ./examples/client/client.go
```

### Method 3: Download Pre-built Binaries

Check the [GitHub releases page](https://github.com/Kingson4Wu/saturncli/releases) for pre-built binaries for your platform.

## Configuration

### Socket Path Configuration

The communication between client and server happens through a socket path. By default:

- **Unix-like systems**: Uses Unix domain sockets (e.g., `/tmp/saturn.sock`)
- **Windows**: Uses TCP loopback on localhost

You can specify the socket path when creating both clients and servers:

```go
// Server side
server := server.NewServer(&utils.DefaultLogger{}, "/custom/path/saturn.sock")

// Client side  
cli := client.NewClient(&utils.DefaultLogger{}, "/custom/path/saturn.sock")
```

### Platform-Specific Considerations

#### Unix-like Systems (macOS/Linux)
- Uses Unix domain sockets by default for optimal performance
- Requires appropriate file system permissions for socket creation
- Socket files are automatically cleaned up on server startup

#### Windows Systems
- Automatically falls back to TCP connections
- Uses localhost with a configurable port (default: 8096)
- No special file permissions needed

## Verification

After installation, verify that everything is working correctly:

1. Start the example server:
```bash
./saturn_svr
```

2. In another terminal, run a test job:
```bash
./saturn_cli --name hello --param id=33 --param ver=22
```

3. You should see output indicating the job was successful.

## Troubleshooting Common Issues

### Socket Permission Errors
If you encounter permission errors, ensure your user has write access to the socket directory:
```bash
ls -la /tmp/  # Check permissions of socket directory
```

### Connection Refused
If you get connection refused errors, ensure:
- The Saturn server is running
- The socket path matches between client and server
- On Windows, ensure the TCP port is available

### Build Issues
If you encounter build issues, ensure:
- Go 1.19+ is installed (`go version`)
- All dependencies are available (`go mod tidy`)
- Git is available for dependency resolution

## Next Steps

Now that you have Saturn CLI set up, you can:

- Proceed to [Quick Start](./quick-start.md) for immediate hands-on experience
- Check out the [Architecture](../architecture.md) to understand how everything works together
- Review the [API Reference](../client-api.md) for detailed technical information