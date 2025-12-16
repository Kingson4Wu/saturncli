---
sidebar_position: 2
---

# Quick Start

Get up and running with Saturn CLI in just a few minutes. This guide provides a hands-on introduction to the core functionality.

## Prerequisites

Make sure you have Saturn CLI [installed and set up](./setup.md) before proceeding.

## Step 1: Start the Server

First, start the Saturn server that will handle job execution:

```bash
./saturn_svr
```

This creates a server listening on the default socket path (`/tmp/saturn.sock` on Unix-like systems).

## Step 2: Execute Your First Job

In another terminal, run a simple job using the command-line interface:

```bash
./saturn_cli --name hello --param id=42 --param message="Hello Saturn!"
```

You should see output showing the job execution in the server terminal and "Execution Success" in the client terminal.

## Step 3: Register Your Own Job

Now let's create a custom job. Create a new file called `custom_server.go`:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Register a custom job
    if err := server.AddJob("greet", func(args map[string]string, signature string) bool {
        name := args["name"]
        if name == "" {
            name = "Anonymous"
        }
        fmt.Printf("Hello, %s! This is job run: %s\n", name, signature)
        return true
    }); err != nil {
        log.Fatal("Failed to register job:", err)
    }
    
    fmt.Println("Server started with 'greet' job. Press Ctrl+C to stop.")
    server.NewServer(&utils.DefaultLogger{}, "/tmp/custom.sock").Serve()
}
```

Build and run this server:
```bash
go run custom_server.go
```

## Step 4: Run Your Custom Job

In another terminal, run your custom job:
```bash
./saturn_cli --name greet --param name=Alice
```

## Step 5: Try a Stoppable Job

Let's create a long-running job that can be stopped:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    if err := server.AddStoppableJob("count", func(args map[string]string, signature string, quit chan struct{}) bool {
        maxCount := 10
        fmt.Printf("Starting counter job: %s\n", signature)
        
        for i := 1; i <= maxCount; i++ {
            select {
            case <-quit:
                fmt.Printf("Counter job %s stopped at count %d\n", signature, i)
                return true
            case <-time.After(1 * time.Second):
                fmt.Printf("Count: %d\n", i)
            }
        }
        
        fmt.Printf("Counter job %s completed\n", signature)
        return true
    }); err != nil {
        panic(err)
    }
    
    server.NewServer(&utils.DefaultLogger{}, "/tmp/count.sock").Serve()
}
```

Run this server, then in another terminal start the counting job:
```bash
./saturn_cli --name count
```

While it's running, stop it from a third terminal:
```bash
./saturn_cli --name count --stop
```

Or simply press `Ctrl+C` in the terminal where you started the counting job.

## Step 6: Embed in Your Application

Here's how to use Saturn CLI from within your Go application:

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
    cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/custom.sock")
    
    // Run a job from code
    result := cli.Run(&client.Task{
        Name: "greet",
        Params: map[string]string{
            "name": "Programmatic Client",
        },
    })
    
    switch result {
    case base.SUCCESS:
        fmt.Println("Job executed successfully!")
    case base.FAILURE:
        fmt.Println("Job failed")
    case base.INTERRUPT:
        fmt.Println("Job was interrupted")
    }
}
```

## Understanding the Core Concepts

### Jobs
- **Regular Jobs**: Execute once and return a result
- **Stoppable Jobs**: Can be cancelled using a quit channel

### Parameters
- Use `--param key=value` for structured parameters
- Use `--args "key1=value1&key2=value2"` for legacy query string format
- Parameters are passed as a `map[string]string` to job handlers

### Communication
- Client and server communicate via Unix domain socket (Unix) or TCP (Windows)
- Both must use the same socket path/endpoint
- Results are returned with standardized success/failure/interrupt codes

## Next Steps

- Explore more [Examples](../examples.md) for real-world scenarios
- Learn about [Architecture](../architecture.md) to understand how it works
- Check the [API Reference](../client-api.md) for detailed technical information
- Try the [Embedding Guide](../embedding.md) to integrate Saturn into your services