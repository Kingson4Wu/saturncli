---
sidebar_position: 3
title: Client API Reference - Go Library for Saturn CLI Job Execution
description: Complete API reference for Saturn CLI client package. Learn how to trigger jobs programmatically in your Go applications using Unix domain sockets communication.
keywords: [saturn cli client api, go client library, unix domain sockets api, go job execution api, saturn cli programming interface, golang client]
---

# Client API Reference

The Saturn CLI client package provides programmatic access to the Saturn server. This document details the client API and how to use it in your Go applications for executing background jobs and managing processes.

## Package Overview

The client package allows you to trigger jobs on a Saturn server from your Go code. It handles the communication protocol, parameter serialization, and result parsing.

Import the client package:

```go
import "github.com/Kingson4Wu/saturncli/client"
```

## Creating a Client

### NewClient

Creates a new client instance with the specified logger and socket path:

```go
func NewClient(logger utils.Logger, sockPath string) *cli
```

**Parameters:**
- `logger`: A logger implementation that satisfies the `utils.Logger` interface
- `sockPath`: Path to the Unix domain socket or TCP address of the server

**Example:**

```go
cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/saturn.sock")
```

## The Task Structure

The `Task` structure represents a job to be executed on the server:

```go
type Task struct {
    Name      string            // Required: Name of the job to execute
    Args      string            // Optional: Legacy query string (merged with Params)
    Stop      bool              // Optional: Send stop signal instead of starting job
    Signature string            // Optional: Target specific run when stopping
    Params    map[string]string // Optional: Structured parameters for the job
}
```

### Task Fields

- `Name`: Required. Identifies the job to execute on the server.
- `Args`: Optional. Legacy query string format (e.g., "id=33&ver=22") for backward compatibility. Merged with Params.
- `Stop`: Optional. When true, sends a stop signal to the job instead of running it.
- `Signature`: Optional. When stopping, targets a specific job run by signature.
- `Params`: Optional. Structured key-value parameters for the job (e.g., `map[string]string{"id": "42"}`).

## Running Tasks

### Run

Executes a task on the Saturn server and returns the result:

```go
func (c *cli) Run(task *Task) base.Result
```

**Parameters:**
- `task`: Pointer to the `Task` structure defining the job to run

**Returns:**
- `base.Result`: One of `base.SUCCESS`, `base.FAILURE`, or `base.INTERRUPT`

**Example:**

```go
result := cli.Run(&client.Task{
    Name: "hello",
    Params: map[string]string{
        "id": "42",
        "version": "2.0",
    },
})
```

## Result Types

Results are returned as constants from the `base` package:

- `base.SUCCESS`: Job completed successfully
- `base.FAILURE`: Job failed during execution
- `base.INTERRUPT`: Job was interrupted (e.g., by stop signal)

## Command-Line Interface Wrapper

The `cmd` package provides a command-line interface wrapper around the client functionality:

### NewCmd

Creates a new command-line wrapper:

```go
func NewCmd(logger utils.Logger, sockPath string) *cmd
```

**Example:**

```go
cmd := client.NewCmd(&utils.DefaultLogger{}, "/tmp/saturn.sock")
cmd.RunWithArgs(os.Args[1:])  // Pass command-line arguments
```

## Complete Usage Example

Here's a complete example showing how to use the client in your Go application:

```go
package main

import (
    "fmt"
    
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Create a client
    cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/saturn.sock")
    
    // Define a task
    task := &client.Task{
        Name: "hello",
        Params: map[string]string{
            "id": "42",
            "message": "greeting",
        },
    }
    
    // Execute the task
    result := cli.Run(task)
    
    // Handle the result
    switch result {
    case base.SUCCESS:
        fmt.Println("Job executed successfully!")
    case base.INTERRUPT:
        fmt.Println("Job was interrupted")
    case base.FAILURE:
        fmt.Println("Job failed")
    default:
        fmt.Println("Unknown result")
    }
    
    // To stop a job
    stopTask := &client.Task{
        Name:      "long_running_job",
        Stop:      true,
        Signature: "specific-job-signature",  // Optional: target specific job run
    }
    
    stopResult := cli.Run(stopTask)
    fmt.Printf("Stop result: %v\n", stopResult)
}
```

## Best Practices

1. **Error Handling**: Always check the return result of `Run()` calls.

2. **Connection Management**: The client handles connection establishment and teardown automatically.

3. **Parameter Validation**: Validate task parameters before sending them to the server.

4. **Logging**: Use the logger consistently for debugging and monitoring.

5. **Concurrency**: The client is safe for concurrent use by multiple goroutines.

## See Also

- [Server API Reference](./server-api.md) - For creating custom Saturn servers
- [Embedding Guide](./embedding.md) - For integrating Saturn into your services
- [CLI Reference](./cli-reference.md) - For command-line usage