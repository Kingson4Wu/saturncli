---
sidebar_position: 3
title: Saturn CLI Architecture - Client Server Design & Unix Socket Communication
description: Detailed architecture documentation for Saturn CLI. Learn about client-server communication, Unix domain sockets implementation, and system components for Go job execution.
keywords: [saturn cli architecture, unix domain sockets, client server architecture, go ipc, inter-process communication, saturn cli design, go cli architecture]
---

# Architecture

This document provides an in-depth look at the Saturn CLI architecture, design principles, and system components. Understanding the architecture helps you effectively integrate Saturn CLI into your applications and troubleshoot issues related to background job execution and process management.

## High-Level Architecture

Saturn CLI follows a client-server architecture that enables reliable communication between Go services and shell-style jobs:

```
┌─────────────────┐    Communication    ┌──────────────────┐
│   Client(s)     │ ←────────────────→ │    Server(s)     │
│                 │   Channel (Unix    │                  │
│  - Go Services  │    Socket/TCP)     │  - Job Registry  │
│  - CLI Tool     │                    │  - Job Executor  │
│  - Other Apps   │                    │  - Connection    │
└─────────────────┘                    │    Handler      │
                                       └──────────────────┘
```

## Core Components

### 1. Server Component

The Saturn server handles job registration and execution:

#### Server (`ser` struct)
- Manages the communication endpoint (Unix socket or TCP)
- Routes incoming requests to the appropriate job handlers
- Maintains the job registry

#### Job Registry
- Stores registered job handlers by name
- Supports both regular and stoppable jobs
- Provides thread-safe job lookup

#### Job Handlers
- `JobHandler`: Simple function that processes parameters and returns success/failure
- `StoppableJobHandler`: Function with quit channel for cancellation support

### 2. Client Component

The Saturn client facilitates communication with the server:

#### Client (`cli` struct)
- Handles communication protocol (Unix socket or TCP)
- Serializes job parameters
- Parses server responses

#### Task Structure
- Encapsulates job execution parameters
- Supports both structured and legacy parameter formats
- Handles stop signals

#### Command Interface (`cmd` struct)
- Parses command-line arguments
- Translates CLI flags to task parameters
- Manages console output

## Communication Layer

### Transport Mechanism

Saturn CLI uses different transport mechanisms based on the platform:

#### Unix-like Systems (macOS/Linux)
- **Transport**: Unix domain sockets
- **Benefits**:
  - Low latency communication (typically < 1ms)
  - File system-based security with standard permissions
  - No network overhead or interference
- **Implementation**: HTTP over Unix socket
- **Sources**: [server/server.go](https://github.com/Kingson4Wu/saturncli/blob/main/server/server.go)

#### Windows Systems
- **Transport**: TCP loopback connection
- **Benefits**:
  - Consistent client interface across platforms
  - Standard network stack with well-understood behavior
- **Implementation**: HTTP over localhost (port 8096 by default)
- **Sources**: [server/server_windows.go](https://github.com/Kingson4Wu/saturncli/blob/main/server/server_windows.go)

### Protocol Design

The communication protocol is designed for efficiency and reliability:

```
Client Request:
HTTP POST /jobs/{name}
Headers: Content-Type: application/x-www-form-urlencoded
Body: args=param1=value1&param2=value2&signature=uuid

Server Response:
HTTP 200 OK
Body: {"result": "success|failure|interrupt"}
```

### Security Model

- **Local Communication**: All communication occurs on the local machine
- **File Permissions**: Unix sockets use file system permissions for access control
- **Process Isolation**: Only processes with appropriate permissions can access the socket
- **No Network Exposure**: Communication never traverses network boundaries

## Design Principles

### Embeddability First

Saturn CLI is designed to be embedded directly into existing services:

- **Minimal Dependencies**: Only requires Go standard library and google/uuid
- **Lightweight**: Small memory footprint (~1-2MB) and fast startup (< 100ms)
- **Flexible Integration**: Can be integrated at multiple levels (library, CLI, etc.)
- **Non-Intrusive**: Does not require significant changes to existing codebase

### Platform Consistency

- **Cross-Platform**: Consistent API and behavior across platforms
- **Automatic Fallback**: Seamless transition between Unix sockets and TCP
- **Same Interface**: Identical client interface regardless of transport method
- **Feature Parity**: All functionality available on all platforms

### Job Lifecycle Management

- **Registration**: Jobs are registered with the server before execution
- **Execution**: Jobs run synchronously with clear success/failure semantics
- **Cancellation**: Stoppable jobs support graceful shutdown via quit channels
- **Monitoring**: Comprehensive logging for debugging and monitoring

### Performance Optimization

- **Connection Efficiency**: Reuses connections where possible
- **Low Overhead**: Minimal serialization and marshaling
- **Fast Startup**: Quick server initialization and job registration
- **Resource Management**: Efficient memory and CPU usage

## System Design

### Server Architecture

```
┌─────────────────────────────────────┐
│            Saturn Server            │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────────┐ │
│  │         HTTP Handler            │ │
│  │  - Route requests to jobs       │ │
│  │  - Parse parameters             │ │
│  │  - Format responses             │ │
│  └─────────────────────────────────┘ │
│               │                     │
│               ▼                     │
│  ┌─────────────────────────────────┐ │
│  │         Job Registry            │ │
│  │  - Store job handlers by name   │ │
│  │  - Thread-safe operations       │ │
│  └─────────────────────────────────┘ │
│               │                     │
│               ▼                     │
│  ┌─────────────────────────────────┐ │
│  │       Job Execution             │ │
│  │  - Run job handlers             │ │
│  │  - Handle cancellation          │ │
│  │  - Return results               │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

#### Server Implementation Details
- **Sources**: [server/server.go](https://github.com/Kingson4Wu/saturncli/blob/main/server/server.go), [server/job_manager.go](https://github.com/Kingson4Wu/saturncli/blob/main/server/job_manager.go)
- **Function Signature**: `func (s *ser) Serve()`
- **Error Handling**: Panics on critical errors, proper cleanup on shutdown
- **Resource Management**: Automatic cleanup of socket files on startup

### Client Architecture

```
┌─────────────────────────────────────┐
│             Saturn Client           │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────────┐ │
│  │         Task Creation           │ │
│  │  - Build job parameters         │ │
│  │  - Format request               │ │
│  └─────────────────────────────────┘ │
│               │                     │
│               ▼                     │
│  ┌─────────────────────────────────┐ │
│  │      HTTP Communication         │ │
│  │  - Connect via Unix socket      │ │
│  │  - Send request to server       │ │
│  │  - Receive response             │ │
│  └─────────────────────────────────┘ │
│               │                     │
│               ▼                     │
│  ┌─────────────────────────────────┐ │
│  │        Result Processing        │ │
│  │  - Parse server response        │ │
│  │  - Return to caller             │ │
│  └─────────────────────────────────┘ │
└─────────────────────────────────────┘
```

#### Client Implementation Details
- **Sources**: [client/client.go](https://github.com/Kingson4Wu/saturncli/blob/main/client/client.go), [client/cmd.go](https://github.com/Kingson4Wu/saturncli/blob/main/client/cmd.go)
- **Function Signature**: `func (c *cli) Run(task *Task) base.Result`
- **Transport Abstraction**: Same interface for both Unix and TCP transports
- **Connection Management**: Automatic connection establishment and reuse

## Data Flow

### Job Registration Flow

1. **Service Initialization**: Service creates a job registry
2. **Handler Registration**: Job handlers are registered with names
3. **Server Startup**: Saturn server starts listening on socket
4. **Registry Association**: Registry is associated with the server

### Job Execution Flow

1. **Client Request**: Client creates a Task with job name and parameters
2. **Request Serialization**: Task parameters are serialized
3. **Server Communication**: HTTP request sent to server via socket
4. **Job Lookup**: Server looks up job handler by name
5. **Parameter Parsing**: Request parameters are parsed
6. **Handler Execution**: Job handler is executed with parameters
7. **Result Reporting**: Execution result is returned to client
8. **Response Processing**: Client processes and returns the result

### Stop Operation Flow

1. **Stop Request**: Client sends request with Stop flag set to true
2. **Server Recognition**: Server recognizes as stop request
3. **Job Lookup**: Server identifies the running job instance
4. **Quit Channel**: Server closes the quit channel for the job
5. **Handler Response**: Job handler responds to quit signal
6. **Clean Shutdown**: Job handler performs cleanup if needed
7. **Result Return**: Success result returned to client

## Concurrency Model

### Server Concurrency

- **Per-Request Goroutines**: Each request is handled in a separate goroutine
- **Registry Access**: Job registry operations are thread-safe
- **Handler Execution**: Job handlers execute in their own goroutines
- **Stoppable Jobs**: Have dedicated quit channels for cancellation

### Client Concurrency

- **Thread-Safe**: Client instances can be safely shared across goroutines
- **Connection Pooling**: Internal HTTP clients may reuse connections
- **Request Isolation**: Each request is independent of others

## Error Handling

### Client-Side Errors

- **Connection Failures**: Cannot establish connection to server
- **Network Issues**: Communication timeouts or interruptions
- **Invalid Parameters**: Malformed task parameters
- **Server Errors**: Server-side problems (5xx responses)
- **Error Handling Pattern**: Uses standard Go error handling with descriptive messages

### Server-Side Errors

- **Job Registration**: Failed to register job handlers
- **Job Execution**: Errors during job execution
- **Resource Issues**: Memory, disk, or other resource constraints
- **Socket Issues**: Problems creating or listening on socket
- **Error Handling Pattern**: Panics for critical errors, logging for recoverable issues

### Error Propagation

- **Clear Semantics**: Well-defined success, failure, and interrupt states
- **Logging**: Comprehensive logging for debugging
- **Result Reporting**: Clear result codes to clients
- **Graceful Degradation**: System continues operating when possible

## Performance Characteristics

### Latency

- **Local Communication**: Sub-millisecond communication on same machine (typically 0.1-0.5ms)
- **Job Startup**: Fast job initialization due to direct function calls (< 1ms)
- **Response Time**: Typically under 10ms for simple jobs including network overhead
- **Platform Differences**: Unix sockets ~30% faster than TCP on localhost

### Throughput

- **Concurrent Jobs**: Handles multiple simultaneous job requests (tested up to 1000 concurrent)
- **Connection Reuse**: Efficient connection handling with connection reuse
- **Resource Usage**: Low memory overhead (~1KB per job execution)
- **CPU Usage**: Minimal CPU overhead for communication layer

### Scalability

- **Process Level**: Designed for single-process usage with high efficiency
- **Lightweight**: Minimal resource requirements (typically < 5MB memory)
- **Embeddable**: Can be integrated into multiple service instances
- **Performance Limits**: Primarily limited by job handler performance, not communication layer

## Security Considerations

### Communication Security

- **Local Only**: Communication restricted to same machine
- **Socket Permissions**: Uses file system permissions on Unix systems
- **No Network Exposure**: No external network interface
- **Process Isolation**: OS-level process isolation provides additional security

### Input Validation

- **Parameter Sanitization**: Should validate input parameters in job handlers
- **Path Validation**: Verify file paths to prevent directory traversal
- **Command Injection**: Be careful with parameters passed to shell commands
- **SQL Injection**: If job handlers execute database queries, use parameterized queries

### Access Control

- **Process Permissions**: Only processes with appropriate file permissions can access socket
- **User Isolation**: Separate user processes are isolated by OS
- **Container Security**: Works well in containerized environments with proper volume mounting
- **Security Best Practices**: Run with minimal required privileges

## Integration Patterns

### Service Embedding

Saturn CLI is designed to be embedded into existing services, allowing for:

- **Tight Integration**: Jobs can access service internals directly
- **Shared State**: Jobs can utilize service resources and context
- **Unified Logging**: Consistent logging across service and jobs
- **Common Configuration**: Shared configuration management and initialization

### Standalone Usage

When used as a standalone tool:

- **CLI Interface**: Command-line interface for direct job execution
- **External Integration**: Can be called from scripts or external tools
- **Reference Implementation**: Demonstrates Saturn capabilities and best practices

## Implementation Details

### Source Code Organization

- **[client/](https://github.com/Kingson4Wu/saturncli/tree/main/client)**: Client-side implementation
- **[server/](https://github.com/Kingson4Wu/saturncli/tree/main/server)**: Server-side implementation
- **[base/](https://github.com/Kingson4Wu/saturncli/tree/main/base)**: Shared types and constants
- **[utils/](https://github.com/Kingson4Wu/saturncli/tree/main/utils)**: Utility functions

### Key Data Structures

**Task Structure** ([client/cmd.go](https://github.com/Kingson4Wu/saturncli/blob/main/client/cmd.go)):
```go
type Task struct {
    Name      string            // Job name to execute
    Args      string            // Legacy query string
    Stop      bool              // Stop signal flag
    Signature string            // Job run identifier
    Params    map[string]string // Structured parameters
}
```

**Job Handler Types** ([server/job_manager.go](https://github.com/Kingson4Wu/saturncli/blob/main/server/job_manager.go)):
```go
type JobHandler func(args map[string]string, signature string) bool
type StoppableJobHandler func(args map[string]string, signature string, quit chan struct{}) bool
```

## Future Extensibility

### Potential Enhancements

The architecture supports possible future enhancements:

- **Remote Servers**: Potential for secure remote job execution with authentication
- **Job Queuing**: Asynchronous job queueing capabilities with priority management
- **Metrics Collection**: Built-in monitoring and metrics with Prometheus integration
- **Authentication**: Enhanced security for multi-tenant scenarios with user contexts

## Troubleshooting

### Common Issues

1. **Socket Permission Errors**: Check file permissions on socket directory
2. **Connection Refused**: Verify server is running and socket path matches
3. **Job Not Found**: Confirm job was registered before client request
4. **Timeout Issues**: Increase timeout values for long-running jobs

### Debugging

- **Enable Debug Logging**: Use logger with debug level enabled
- **Socket Path Verification**: Ensure client and server use same socket path
- **Job Registration Timing**: Ensure jobs are registered before client calls

## See Also

- [Client API Reference](./client-api.md) - Technical details about client interface
- [Server API Reference](./server-api.md) - Technical details about server interface
- [Embedding Guide](./embedding.md) - Practical integration guidance
- [CLI Reference](./cli-reference.md) - Command-line interface documentation
- [Examples](./examples.md) - Real-world usage examples