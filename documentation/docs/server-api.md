---
sidebar_position: 4
title: Server API Reference - Go Server Toolkit for Job Execution
description: Complete API reference for Saturn CLI server package. Learn how to create job execution servers and register jobs in your Go applications.
keywords: [saturn cli server api, go server library, job execution server, saturn cli server, go background job server, golang server api]
---

# Server API Reference

The Saturn CLI server package provides the infrastructure for creating job execution servers. This document details the server API and how to use it in your Go applications for managing background processes.

## Package Overview

The server package enables you to create custom Saturn job servers that can execute tasks sent by clients. It handles job registration, execution, and lifecycle management.

Import the server package:

```go
import "github.com/Kingson4Wu/saturncli/server"
```

## Creating a Server

### NewServer

Creates a new server instance with the specified logger and socket path:

```go
func NewServer(logger utils.Logger, sockPath string, opts ...Option) *ser
```

**Parameters:**
- `logger`: A logger implementation that satisfies the `utils.Logger` interface
- `sockPath`: Path to the Unix domain socket or TCP address for the server
- `opts`: Optional configuration options (see [Options](#options) below)

**Example:**

```go
server := server.NewServer(&utils.DefaultLogger{}, "/tmp/saturn.sock")
server.Serve()
```

## Options

### WithRegistry

Allows you to specify a custom job registry instead of using the default global registry:

```go
func WithRegistry(registry *Registry) Option
```

**Example:**

```go
registry := server.NewRegistry()
server := server.NewServer(
    &utils.DefaultLogger{}, 
    "/tmp/saturn.sock", 
    server.WithRegistry(registry),
)
```

## Server Methods

### Serve

Starts the server and begins listening for client connections. This method blocks the calling goroutine:

```go
func (s *ser) Serve()
```

## Job Registration

Saturn supports two types of jobs: regular jobs and stoppable jobs.

### Global Registration Functions

The server package provides global registration functions that operate on the default global registry:

#### AddJob

Registers a regular (non-stoppable) job globally:

```go
func AddJob(name string, handler JobHandler) error
```

#### AddStoppableJob

Registers a stoppable job globally:

```go
func AddStoppableJob(name string, handler StoppableJobHandler) error
```

### Registry-Based Registration

For more control, create a specific registry:

#### NewRegistry

Creates a new isolated job registry:

```go
func NewRegistry() *Registry
```

#### Registry.AddJob

Registers a regular job in the specified registry:

```go
func (r *Registry) AddJob(name string, handler JobHandler) error
```

**Parameters:**
- `name`: Unique identifier for the job
- `handler`: Function that implements the job logic

#### Registry.AddStoppableJob

Registers a stoppable job in the specified registry:

```go
func (r *Registry) AddStoppableJob(name string, handler StoppableJobHandler) error
```

**Parameters:**
- `name`: Unique identifier for the job  
- `handler`: Function that implements the stoppable job logic

## Job Handler Types

### JobHandler

Function type for regular jobs:

```go
type JobHandler func(args map[string]string, signature string) bool
```

**Parameters:**
- `args`: Map of parameters passed from the client
- `signature`: Unique identifier for this job run

**Return Value:**
- `bool`: `true` for success, `false` for failure

### StoppableJobHandler

Function type for stoppable jobs:

```go
type StoppableJobHandler func(args map[string]string, signature string, quit chan struct{}) bool
```

**Parameters:**
- `args`: Map of parameters passed from the client
- `signature`: Unique identifier for this job run
- `quit`: Channel that signals when the job should stop

**Return Value:**
- `bool`: `true` for success, `false` for failure

## Job Handler Patterns

### Regular Job Pattern

Regular jobs execute and return without the ability to be stopped externally:

```go
func myJobHandler(args map[string]string, signature string) bool {
    // Process arguments
    id := args["id"]
    version := args["version"]
    
    // Perform job work
    log.Printf("Executing job %s with id=%s, version=%s", signature, id, version)
    
    // Do the work...
    
    // Return true for success, false for failure
    return true
}
```

### Stoppable Job Pattern

Stoppable jobs must regularly check the quit channel to allow graceful shutdown:

```go
func myStoppableJobHandler(args map[string]string, signature string, quit chan struct{}) bool {
    // Process arguments
    id := args["id"]
    
    // Use a ticker to periodically check for quit signal
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for i := 0; ; i++ {
        select {
        case <-quit:
            log.Printf("Job %s received quit signal at iteration %d", signature, i)
            return true // Indicate successful clean shutdown
        case <-ticker.C:
            // Do work for this iteration
            log.Printf("Processing iteration %d for job %s", i, signature)
            
            // Simulate work
            time.Sleep(time.Second)
        }
    }
}
```

## Complete Usage Example

Here's a complete example showing how to create a custom Saturn server:

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Create a registry to organize jobs
    registry := server.NewRegistry()
    
    // Register a regular job
    if err := registry.AddJob("hello", func(args map[string]string, signature string) bool {
        log.Printf("Hello job running with args: %+v, signature: %s", args, signature)
        
        // Simulate some processing time
        time.Sleep(1 * time.Second)
        
        // Optionally return false if job fails
        if args["fail"] == "true" {
            log.Printf("Failing job %s intentionally", signature)
            return false
        }
        
        return true
    }); err != nil {
        log.Fatal("Failed to add job:", err)
    }
    
    // Register a stoppable job
    if err := registry.AddStoppableJob("countdown", func(args map[string]string, signature string, quit chan struct{}) bool {
        count := 10
        if args["count"] != "" {
            fmt.Sscanf(args["count"], "%d", &count)
        }
        
        log.Printf("Countdown job %s starting for %d seconds", signature, count)
        
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for i := 0; i < count; i++ {
            select {
            case <-quit:
                log.Printf("Countdown job %s stopped early at step %d", signature, i)
                return true
            case <-ticker.C:
                log.Printf("Countdown job %s: %d remaining", signature, count-i)
            }
        }
        
        log.Printf("Countdown job %s completed normally", signature)
        return true
    }); err != nil {
        log.Fatal("Failed to add stoppable job:", err)
    }
    
    // Create and start the server with our custom registry
    server.NewServer(
        &utils.DefaultLogger{}, 
        "/tmp/my_custom_saturn.sock",
        server.WithRegistry(registry),
    ).Serve()
}
```

## Server Lifecycle

1. Create a registry (optional, defaults to global registry)
2. Register your jobs with handlers
3. Create a server instance
4. Call `Serve()` to start the server
5. The server listens for client connections until terminated
6. Unix socket files are cleaned up automatically on startup

## Platform-Specific Behavior

- **Unix-like systems (macOS/Linux)**: Uses Unix domain sockets
- **Windows**: Falls back to TCP connections on localhost

## Best Practices

1. **Unique Job Names**: Ensure job names are unique within each registry to prevent conflicts.

2. **Proper Cleanup**: Stoppable jobs must properly handle the quit channel to allow graceful shutdown.

3. **Error Handling**: Return `false` from job handlers when errors occur, `true` for success.

4. **Logging**: Use the provided logger for consistent logging across your application.

5. **Thread Safety**: Job handlers may be called concurrently, so ensure thread-safe operations.

6. **Resource Management**: Clean up resources properly in job handlers to prevent leaks.

7. **Validation**: Validate input parameters in job handlers to prevent unexpected behavior.

## See Also

- [Client API Reference](./client-api.md) - For using the Saturn client in your applications
- [Embedding Guide](./embedding.md) - For integrating Saturn into your services
- [Architecture](./architecture.md) - For understanding the system design