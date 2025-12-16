---
sidebar_position: 3
---

# Getting Started Examples

This page provides practical examples to help you get familiar with Saturn CLI concepts and usage patterns.

## Simple Job Example

Let's start with a basic job registration and execution:

```go
// server.go
package main

import (
    "fmt"
    "log"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Register a simple job
    if err := server.AddJob("echo", func(args map[string]string, signature string) bool {
        message := args["message"]
        fmt.Printf("[%s] Echo: %s\n", signature, message)
        return true
    }); err != nil {
        log.Fatal("Failed to register job:", err)
    }
    
    fmt.Println("Server running on /tmp/echo.sock")
    server.NewServer(&utils.DefaultLogger{}, "/tmp/echo.sock").Serve()
}
```

Execute with:
```bash
go run server.go
```

In another terminal:
```bash
./saturn_cli --name echo --param message="Hello, World!"
```

## Stoppable Job with Progress

A more advanced example showing a job that reports progress and can be stopped:

```go
// progress_server.go
package main

import (
    "fmt"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    if err := server.AddStoppableJob("progress-task", func(args map[string]string, signature string, quit chan struct{}) bool {
        steps := 20
        if args["steps"] != "" {
            fmt.Sscanf(args["steps"], "%d", &steps)
        }
        
        fmt.Printf("[%s] Starting task with %d steps\n", signature, steps)
        
        completed := 0
        for i := 1; i <= steps; i++ {
            select {
            case <-quit:
                fmt.Printf("[%s] Task cancelled at step %d of %d (completed %d)\n", 
                    signature, i, steps, completed)
                return true
            case <-time.After(200 * time.Millisecond): // Simulate work
                completed = i
                if i%5 == 0 { // Report progress every 5 steps
                    fmt.Printf("[%s] Progress: %d/%d steps completed\n", signature, i, steps)
                }
            }
        }
        
        fmt.Printf("[%s] Task completed successfully (%d steps)\n", signature, steps)
        return true
    }); err != nil {
        panic(err)
    }
    
    server.NewServer(&utils.DefaultLogger{}, "/tmp/progress.sock").Serve()
}
```

Test it:
```bash
# Start the server
go run progress_server.go

# In another terminal, start the task
./saturn_cli --name progress-task --param steps=30

# In a third terminal, stop the task before completion
./saturn_cli --name progress-task --stop
```

## File Processing Example

Here's a practical example of processing files with error handling:

```go
// file_processor.go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    if err := server.AddStoppableJob("file-processor", func(args map[string]string, signature string, quit chan struct{}) bool {
        inputDir := args["input_dir"]
        outputDir := args["output_dir"]
        
        if inputDir == "" || outputDir == "" {
            fmt.Printf("[%s] ERROR: input_dir and output_dir are required\n", signature)
            return false
        }
        
        // Verify directories exist
        if _, err := os.Stat(inputDir); os.IsNotExist(err) {
            fmt.Printf("[%s] ERROR: Input directory does not exist: %s\n", signature, inputDir)
            return false
        }
        
        if err := os.MkdirAll(outputDir, 0755); err != nil {
            fmt.Printf("[%s] ERROR: Failed to create output directory: %v\n", signature, err)
            return false
        }
        
        fmt.Printf("[%s] Processing files from %s to %s\n", signature, inputDir, outputDir)
        
        // Find all .txt files
        files, err := filepath.Glob(filepath.Join(inputDir, "*.txt"))
        if err != nil {
            fmt.Printf("[%s] ERROR: Failed to find files: %v\n", signature, err)
            return false
        }
        
        processed := 0
        for _, file := range files {
            select {
            case <-quit:
                fmt.Printf("[%s] Processing cancelled, %d files processed\n", signature, processed)
                return true
            default:
                if processFile(file, outputDir, signature) {
                    processed++
                    fmt.Printf("[%s] Processed: %s\n", signature, filepath.Base(file))
                }
            }
        }
        
        fmt.Printf("[%s] Completed processing %d files\n", signature, processed)
        return true
    }); err != nil {
        panic(err)
    }
    
    server.NewServer(&utils.DefaultLogger{}, "/tmp/file-processor.sock").Serve()
}

func processFile(inputPath, outputDir, signature string) bool {
    // Simulate file processing (in real implementation, you'd read, process, and write the file)
    time.Sleep(500 * time.Millisecond) // Simulate processing time
    
    filename := filepath.Base(inputPath)
    nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
    outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_processed.txt", nameWithoutExt))
    
    // In a real implementation, you would actually process the file content
    fmt.Printf("[%s] Would write processed content to: %s\n", signature, outputPath)
    
    return true
}
```

## Embedding in a Web Service

This example shows how to embed Saturn CLI in a web service:

```go
// web_service.go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

type WebService struct {
    saturnClient *client.cli
}

type JobRequest struct {
    Name   string            `json:"name"`
    Params map[string]string `json:"params"`
}

type JobResponse struct {
    Success bool   `json:"success"`
    Result  string `json:"result"`
}

func NewWebService() *WebService {
    return &WebService{
        saturnClient: client.NewClient(&utils.DefaultLogger{}, "/tmp/web_service.sock"),
    }
}

func (ws *WebService) JobHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req JobRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    result := ws.saturnClient.Run(&client.Task{
        Name:   req.Name,
        Params: req.Params,
    })
    
    response := JobResponse{
        Success: result == base.SUCCESS,
        Result:  string(result),
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // Set up Saturn jobs
    if err := server.AddJob("data-processor", func(args map[string]string, signature string) bool {
        fmt.Printf("[%s] Processing data with args: %+v\n", signature, args)
        // Simulate processing
        // In real implementation, perform actual data processing
        return true
    }); err != nil {
        log.Fatal("Failed to register job:", err)
    }
    
    // Start Saturn server in a goroutine
    go func() {
        fmt.Println("Starting Saturn server...")
        server.NewServer(&utils.DefaultLogger{}, "/tmp/web_service.sock").Serve()
    }()
    
    // Give server time to start
    time.Sleep(100 * time.Millisecond)
    
    // Set up web API
    webService := NewWebService()
    http.HandleFunc("/job", webService.JobHandler)
    
    fmt.Println("Starting web service on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Testing Your Examples

Create a simple test for your Saturn jobs:

```go
// job_test.go
package main

import (
    "bytes"
    "fmt"
    "strings"
    "testing"
    "time"
)

type TestLogger struct {
    buffer *bytes.Buffer
}

func (tl *TestLogger) Info(msg string)    { tl.buffer.WriteString("INFO: " + msg + "\n") }
func (tl *TestLogger) Infof(format string, args ...interface{}) { tl.buffer.WriteString("INFO: " + fmt.Sprintf(format, args...) + "\n") }
func (tl *TestLogger) Error(msg string)   { tl.buffer.WriteString("ERROR: " + msg + "\n") }
func (tl *TestLogger) Errorf(format string, args ...interface{}) { tl.buffer.WriteString("ERROR: " + fmt.Sprintf(format, args...) + "\n") }
func (tl *TestLogger) Debug(msg string)   { tl.buffer.WriteString("DEBUG: " + msg + "\n") }
func (tl *TestLogger) Debugf(format string, args ...interface{}) { tl.buffer.WriteString("DEBUG: " + fmt.Sprintf(format, args...) + "\n") }
func (tl *TestLogger) Warn(msg string)    { tl.buffer.WriteString("WARN: " + msg + "\n") }
func (tl *TestLogger) Warnf(format string, args ...interface{})  { tl.buffer.WriteString("WARN: " + fmt.Sprintf(format, args...) + "\n") }

func TestEchoJob(t *testing.T) {
    var logBuffer bytes.Buffer
    testLogger := &TestLogger{buffer: &logBuffer}
    
    // Test parameters
    args := map[string]string{"message": "Test message"}
    signature := "test-signature"
    
    // Create job handler function
    handler := func(jobArgs map[string]string, jobSignature string) bool {
        message := jobArgs["message"]
        testLogger.Infof("Echo: %s (Signature: %s)", message, jobSignature)
        return true
    }
    
    // Execute the job
    result := handler(args, signature)
    
    if !result {
        t.Error("Expected job to return true")
    }
    
    logOutput := logBuffer.String()
    if !strings.Contains(logOutput, "Test message") {
        t.Errorf("Expected log to contain 'Test message', got: %s", logOutput)
    }
    
    if !strings.Contains(logOutput, "test-signature") {
        t.Errorf("Expected log to contain signature, got: %s", logOutput)
    }
}

func TestStoppableJob(t *testing.T) {
    quit := make(chan struct{})
    done := make(chan bool, 1)
    
    handler := func(args map[string]string, signature string, quitChan chan struct{}) bool {
        iterations := 0
        maxIterations := 10
        
        for iterations < maxIterations {
            select {
            case <-quitChan:
                return true
            default:
                iterations++
                time.Sleep(10 * time.Millisecond) // Small sleep for the test
            }
        }
        
        return true
    }
    
    args := map[string]string{}
    signature := "test-stop"
    
    // Run the handler in a goroutine
    go func() {
        result := handler(args, signature, quit)
        done <- result
    }()
    
    // Allow some iterations
    time.Sleep(50 * time.Millisecond)
    
    // Close the quit channel to stop the job
    close(quit)
    
    // Wait for the handler to finish
    select {
    case result := <-done:
        if !result {
            t.Error("Expected stoppable job to return true after being stopped")
        }
    case <-time.After(1 * time.Second):
        t.Error("Test timed out")
    }
}
```

## Running Examples

To run any of these examples:

1. Create a new directory for the example
2. Create the Go file with the example code
3. Run `go mod init example_name`
4. Run `go get github.com/Kingson4Wu/saturncli`
5. Build with `go build -o example_name main.go`
6. Run the server: `./example_name`
7. In another terminal, use the Saturn CLI or create a client to trigger jobs

## Source Code References

- [Full server example](https://github.com/Kingson4Wu/saturncli/blob/main/examples/server/server.go)
- [Full client example](https://github.com/Kingson4Wu/saturncli/blob/main/examples/client/client.go)
- [Client implementation](https://github.com/Kingson4Wu/saturncli/blob/main/client/client.go)
- [Server implementation](https://github.com/Kingson4Wu/saturncli/blob/main/server/server.go)

## Next Steps

- Review the [API Reference](../client-api.md) for detailed function documentation
- Check out the [Embedding Guide](../embedding.md) for more integration patterns
- Look at the [Architecture](../architecture.md) to understand the system design