---
sidebar_position: 2
title: Getting Started with Saturn CLI - Installation & Setup Guide
description: Learn how to install, set up, and run your first Saturn CLI jobs. Complete guide for beginners to get started with Go job execution and Unix domain sockets.
keywords: [saturn cli getting started, go cli setup, install saturn cli, unix domain sockets tutorial, go job execution setup, cli automation]
---

# Getting Started with Saturn CLI - Installation & Setup Guide

This guide will walk you through installing Saturn CLI, building the binaries, and running your first example to execute background jobs in Go applications.

## Prerequisites

- Go 1.19 or newer
- Unix-like system for socket transport (macOS/Linux) with Unix domain socket support
- On Windows systems, TCP transport is used instead

## Installation

### Method 1: Clone and Build from Source

```bash
git clone https://github.com/Kingson4Wu/saturncli.git
cd saturncli
```

### Method 2: Use Go Modules

Add Saturn CLI as a module to your existing Go project:

```bash
go mod init your-project-name
go get github.com/Kingson4Wu/saturncli
```

## Building the Binaries

To build both the server and client binaries:

```bash
make
```

This produces:

- `saturn_svr` – reference server demonstrating job registration
- `saturn_cli` – command-line client

Alternative build command if Make is not available:

```bash
go build -o saturn_svr ./examples/server/server.go
go build -o saturn_cli ./examples/client/client.go
```

## Running the Demo

1. Start the server:

```bash
./saturn_svr
```

2. In another terminal, run a simple job:

```bash
./saturn_cli --name hello --param id=33 --param ver=22
```

3. Run a stoppable job:

```bash
./saturn_cli --name hello_stoppable
```

4. Stop the running job:

```bash
./saturn_cli --name hello_stoppable --stop
```

Alternatively, press `Ctrl+C` while the stoppable job is running to trigger an interrupt with automatic stop propagation.

## Using Saturn CLI as a Library

Saturn CLI is designed to be embedded directly into your existing services. Here's how to register jobs programmatically:

### Setting up a Server with Custom Jobs

```go
package main

import (
    "log"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    registry := server.NewRegistry()

    // Register a regular job
    if err := registry.AddJob("hello", func(args map[string]string, signature string) bool {
        log.Printf("Hello %v (run=%s)", args, signature)
        return true
    }); err != nil {
        log.Fatal(err)
    }

    // Register a stoppable job
    if err := registry.AddStoppableJob("slow-task", func(args map[string]string, signature string, quit chan struct{}) bool {
        log.Printf("Starting slow task %v (run=%s)", args, signature)
        
        // Simulate work with the ability to be stopped
        for i := 0; i < 100; i++ {
            select {
            case <-quit:
                log.Println("Task received quit signal, exiting gracefully")
                return true
            default:
                log.Printf("Processing item %d...", i)
                // Simulate work
                // time.Sleep(1 * time.Second) 
            }
        }
        return true
    }); err != nil {
        log.Fatal(err)
    }

    // Start the server
    server.NewServer(&utils.DefaultLogger{}, "/tmp/saturn.sock", server.WithRegistry(registry)).Serve()
}
```

### Creating a Client to Trigger Jobs

```go
package main

import (
    "fmt"
    
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/saturn.sock")
    
    result := cli.Run(&client.Task{
        Name:   "hello",
        Params: map[string]string{"id": "42", "user": "admin"},
    })
    
    switch result {
    case base.SUCCESS:
        fmt.Println("Job executed successfully")
    case base.INTERRUPT:
        fmt.Println("Job was interrupted")
    default:
        fmt.Println("Job failed")
    }
}
```

## Understanding the Communication Protocol

Saturn CLI uses different transport mechanisms depending on the platform:

- **Unix-like systems (macOS/Linux)**: Unix domain sockets for low-latency, secure communication
- **Windows**: TCP loopback connection for compatibility

The communication is structured as follows:

- Clients connect to servers via the registered socket/TCP endpoint
- Jobs are identified by name and can accept parameters
- Results are returned with success/failure/interrupt status

## Configuration Options

When creating a server, you can customize:

- Socket path location
- Registry to use
- Logger implementation

## Next Steps

- Learn about [Server Architecture](./architecture.md)
- Explore the [Client API Reference](./client-api.md)
- Check out more [Examples](./examples.md)