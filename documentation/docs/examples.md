---
sidebar_position: 8
title: Saturn CLI Examples - Real World Go Job Execution Use Cases
description: Practical examples of Saturn CLI usage in various scenarios. From basic implementations to advanced patterns for background job execution in Go.
keywords: [saturn cli examples, go job execution examples, unix domain sockets examples, background job examples, golang cli examples]
---

# Examples

This page provides practical examples of Saturn CLI usage in various scenarios, from basic implementations to advanced patterns for executing background jobs in Go applications.

## Basic Examples

### Simple Job Execution

The most basic example of using Saturn CLI is to register a simple job and execute it:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    // Register a simple job
    if err := server.AddJob("hello", func(args map[string]string, signature string) bool {
        name := args["name"]
        if name == "" {
            name = "World"
        }
        fmt.Printf("Hello, %s! (Run: %s)\n", name, signature)
        return true
    }); err != nil {
        log.Fatal(err)
    }
    
    // Start the server
    server.NewServer(&utils.DefaultLogger{}, "/tmp/hello.sock").Serve()
}
```

To run this example:
1. Build: `go build -o hello_server main.go`
2. Start server: `./hello_server`
3. In another terminal: `./saturn_cli --name hello --param name=Alice`

### Stoppable Job Example

A job that can be cancelled gracefully:

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
    if err := server.AddStoppableJob("countdown", func(args map[string]string, signature string, quit chan struct{}) bool {
        duration := 10 // default to 10 seconds
        if args["duration"] != "" {
            // In a real app, parse this safely
            fmt.Sscanf(args["duration"], "%d", &duration)
        }
        
        fmt.Printf("Starting countdown for %d seconds (Run: %s)\n", duration, signature)
        
        for i := duration; i >= 0; i-- {
            select {
            case <-quit:
                fmt.Printf("Countdown stopped early at %d seconds (Run: %s)\n", i, signature)
                return true
            case <-time.After(1 * time.Second):
                fmt.Printf("%d... ", i)
                if i == 0 {
                    fmt.Println("Liftoff!")
                }
            }
        }
        return true
    }); err != nil {
        log.Fatal(err)
    }
    
    server.NewServer(&utils.DefaultLogger{}, "/tmp/countdown.sock").Serve()
}
```

## Real-World Examples

### File Processing Service

A service that processes files with progress tracking:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

type FileProcessor struct {
    logger utils.Logger
}

func NewFileProcessor(logger utils.Logger) *FileProcessor {
    return &FileProcessor{logger: logger}
}

func (fp *FileProcessor) ProcessFilesJob(args map[string]string, signature string, quit chan struct{}) bool {
    inputDir := args["input_dir"]
    outputDir := args["output_dir"]
    
    if inputDir == "" || outputDir == "" {
        fp.logger.Errorf("Missing required parameters for job %s", signature)
        return false
    }
    
    // Validate directories exist
    if _, err := os.Stat(inputDir); os.IsNotExist(err) {
        fp.logger.Errorf("Input directory does not exist: %s", inputDir)
        return false
    }
    
    if _, err := os.Stat(outputDir); os.IsNotExist(err) {
        fp.logger.Errorf("Output directory does not exist: %s", outputDir)
        return false
    }
    
    fp.logger.Infof("Starting file processing job %s from %s to %s", signature, inputDir, outputDir)
    
    // Get all .txt files in the input directory
    files, err := filepath.Glob(filepath.Join(inputDir, "*.txt"))
    if err != nil {
        fp.logger.Errorf("Error finding files: %v", err)
        return false
    }
    
    if len(files) == 0 {
        fp.logger.Infof("No .txt files found in %s", inputDir)
        return true
    }
    
    processedCount := 0
    for _, file := range files {
        select {
        case <-quit:
            fp.logger.Infof("File processing job %s cancelled after processing %d files", signature, processedCount)
            return true
        default:
            if fp.processFile(file, outputDir, signature) {
                processedCount++
            }
        }
    }
    
    fp.logger.Infof("File processing job %s completed. Processed %d files", signature, processedCount)
    return true
}

func (fp *FileProcessor) processFile(inputPath, outputDir, signature string) bool {
    // Simulate file processing
    fp.logger.Debugf("Processing file: %s for job %s", inputPath, signature)
    
    // In a real implementation, you would:
    // 1. Read the file
    // 2. Process its contents
    // 3. Write to the output directory
    // 4. Handle errors appropriately
    
    // Simulate processing time
    time.Sleep(500 * time.Millisecond)
    
    // Generate output filename
    filename := filepath.Base(inputPath)
    nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
    outputPath := filepath.Join(outputDir, fmt.Sprintf("%s_processed.txt", nameWithoutExt))
    
    fp.logger.Debugf("Would write processed output to: %s", outputPath)
    
    return true
}

func main() {
    logger := &utils.DefaultLogger{}
    processor := NewFileProcessor(logger)
    
    registry := server.NewRegistry()
    
    if err := registry.AddStoppableJob("process_files", processor.ProcessFilesJob); err != nil {
        log.Fatal(err)
    }
    
    server.NewServer(logger, "/tmp/file_processor.sock", server.WithRegistry(registry)).Serve()
}
```

### Data Backup Service

A backup service that can be stopped gracefully:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "path/filepath"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

type BackupService struct {
    logger utils.Logger
}

func NewBackupService(logger utils.Logger) *BackupService {
    return &BackupService{logger: logger}
}

func (bs *BackupService) BackupJob(args map[string]string, signature string, quit chan struct{}) bool {
    source := args["source"]
    destination := args["destination"]
    compress := args["compress"] == "true"
    
    if source == "" || destination == "" {
        bs.logger.Errorf("Missing source or destination for backup job %s", signature)
        return false
    }
    
    // Check if source exists
    if _, err := os.Stat(source); os.IsNotExist(err) {
        bs.logger.Errorf("Source does not exist: %s", source)
        return false
    }
    
    bs.logger.Infof("Starting backup job %s: %s -> %s (compress: %t)", 
                     signature, source, destination, compress)
    
    // Create backup directory if it doesn't exist
    if err := os.MkdirAll(destination, 0755); err != nil {
        bs.logger.Errorf("Failed to create destination directory: %v", err)
        return false
    }
    
    // Simulate backup process
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()
    
    backupSteps := []string{"scanning", "copying", "verifying", "finalizing"}
    for _, step := range backupSteps {
        select {
        case <-quit:
            bs.logger.Infof("Backup job %s cancelled during %s step", signature, step)
            return true
        case <-ticker.C:
            bs.logger.Infof("Backup job %s: %s...", step)
            // Simulate work for each step
            time.Sleep(1 * time.Second)
        }
    }
    
    // In a real implementation, you would actually copy the files
    timestamp := time.Now().Format("2006-01-02-15-04-05")
    backupName := fmt.Sprintf("backup_%s_%s", filepath.Base(source), timestamp)
    backupPath := filepath.Join(destination, backupName)
    
    bs.logger.Infof("Created backup: %s", backupPath)
    bs.logger.Infof("Backup job %s completed successfully", signature)
    
    return true
}

func main() {
    logger := &utils.DefaultLogger{}
    backupService := NewBackupService(logger)
    
    registry := server.NewRegistry()
    
    if err := registry.AddStoppableJob("backup_data", backupService.BackupJob); err != nil {
        log.Fatal(err)
    }
    
    server.NewServer(logger, "/tmp/backup_service.sock", server.WithRegistry(registry)).Serve()
}
```

### Monitoring Service

A service that continuously monitors system resources:

```go
package main

import (
    "fmt"
    "log"
    "math/rand"
    "time"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

type MonitorService struct {
    logger utils.Logger
}

func NewMonitorService(logger utils.Logger) *MonitorService {
    return &MonitorService{logger: logger}
}

func (ms *MonitorService) SystemMonitorJob(args map[string]string, signature string, quit chan struct{}) bool {
    interval := 5 // default to 5 seconds
    if args["interval"] != "" {
        fmt.Sscanf(args["interval"], "%d", &interval)
    }
    
    ms.logger.Infof("Starting system monitor job %s with %d second intervals", signature, interval)
    
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    defer ticker.Stop()
    
    iteration := 0
    for {
        select {
        case <-quit:
            ms.logger.Infof("System monitor job %s stopped after %d iterations", signature, iteration)
            return true
        case <-ticker.C:
            iteration++
            // Simulate monitoring metrics
            cpu := rand.Float64() * 100 // Random CPU usage 0-100%
            memory := rand.Float64() * 100 // Random memory usage 0-100%
            disk := rand.Float64() * 100 // Random disk usage 0-100%
            
            ms.logger.Infof("Monitor %s - Iteration %d: CPU=%.2f%%, Memory=%.2f%%, Disk=%.2f%%", 
                            signature, iteration, cpu, memory, disk)
            
            // In a real implementation, you would:
            // 1. Collect actual system metrics
            // 2. Send them to a monitoring system
            // 3. Trigger alerts if thresholds are exceeded
        }
    }
}

func (ms *MonitorService) HealthCheckJob(args map[string]string, signature string) bool {
    service := args["service"]
    if service == "" {
        service = "unknown"
    }
    
    // Simulate health check logic
    ms.logger.Infof("Performing health check for service: %s (Run: %s)", service, signature)
    
    // Simulate checking service status
    // In real implementation, you would check actual service status
    time.Sleep(500 * time.Millisecond)
    
    // Simulate 95% success rate
    success := rand.Float64() < 0.95
    
    if success {
        ms.logger.Infof("Health check for service %s passed", service)
        return true
    } else {
        ms.logger.Errorf("Health check for service %s failed", service)
        return false
    }
}

func main() {
    logger := &utils.DefaultLogger{}
    monitorService := NewMonitorService(logger)
    
    registry := server.NewRegistry()
    
    if err := registry.AddStoppableJob("system_monitor", monitorService.SystemMonitorJob); err != nil {
        log.Fatal(err)
    }
    
    if err := registry.AddJob("health_check", monitorService.HealthCheckJob); err != nil {
        log.Fatal(err)
    }
    
    server.NewServer(logger, "/tmp/monitor_service.sock", server.WithRegistry(registry)).Serve()
}
```

## Advanced Patterns

### Multi-Registry Service

A service that uses multiple registries for different purposes:

```go
package main

import (
    "log"
    
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    logger := &utils.DefaultLogger{}
    
    // Create separate registries for different concerns
    criticalJobs := server.NewRegistry()
    backgroundJobs := server.NewRegistry()
    
    // Register critical jobs
    if err := criticalJobs.AddJob("validate_data", func(args map[string]string, signature string) bool {
        log.Printf("Running critical validation job: %s", signature)
        // Implement critical validation logic
        return true
    }); err != nil {
        log.Fatal(err)
    }
    
    // Register background jobs
    if err := backgroundJobs.AddStoppableJob("cleanup", func(args map[string]string, signature string, quit chan struct{}) bool {
        log.Printf("Running background cleanup: %s", signature)
        // Implement cleanup logic that can be stopped
        return true
    }); err != nil {
        log.Fatal(err)
    }
    
    // You could run multiple servers or combine them
    // For this example, we'll use a single server with the critical jobs registry
    server.NewServer(logger, "/tmp/multi_service.sock", server.WithRegistry(criticalJobs)).Serve()
}
```

### Client Usage Examples

Examples of how to use the client programmatically:

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/example.sock")
    
    // Example 1: Simple job execution
    fmt.Println("Executing simple job...")
    result := cli.Run(&client.Task{
        Name:   "hello",
        Params: map[string]string{"name": "Saturn User"},
    })
    
    switch result {
    case base.SUCCESS:
        fmt.Println("Job completed successfully")
    case base.FAILURE:
        fmt.Println("Job failed")
    case base.INTERRUPT:
        fmt.Println("Job was interrupted")
    }
    
    // Example 2: Job with multiple parameters
    fmt.Println("\nExecuting job with multiple parameters...")
    result = cli.Run(&client.Task{
        Name: "process_data",
        Params: map[string]string{
            "input_file":  "/path/to/input",
            "output_file": "/path/to/output",
            "format":      "json",
            "validate":    "true",
        },
    })
    
    if result == base.SUCCESS {
        fmt.Println("Data processing completed")
    } else {
        fmt.Println("Data processing failed")
    }
    
    // Example 3: Stopping a job
    fmt.Println("\nStopping a running job...")
    result = cli.Run(&client.Task{
        Name:      "long_running_task",
        Stop:      true,
        Signature: "task-to-stop", // Specify which instance to stop
    })
    
    fmt.Printf("Stop result: %v\n", result)
    
    // Example 4: Error handling
    fmt.Println("\nExecuting non-existent job...")
    result = cli.Run(&client.Task{
        Name: "non_existent_job",
    })
    
    if result == base.FAILURE {
        fmt.Println("As expected, job failed because it doesn't exist")
    }
}
```

## Integration Examples

### REST API Integration

Integrating Saturn CLI with a REST API:

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/utils"
)

type APIServer struct {
    saturnClient *client.cli
}

type JobRequest struct {
    Name   string            `json:"name"`
    Params map[string]string `json:"params"`
}

type JobResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func NewAPIServer() *APIServer {
    return &APIServer{
        saturnClient: client.NewClient(&utils.DefaultLogger{}, "/tmp/api_integration.sock"),
    }
}

func (api *APIServer) JobHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req JobRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    result := api.saturnClient.Run(&client.Task{
        Name:   req.Name,
        Params: req.Params,
    })
    
    resp := JobResponse{
        Success: result == base.SUCCESS,
        Message: "Job executed",
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func main() {
    apiServer := NewAPIServer()
    
    http.HandleFunc("/jobs", apiServer.JobHandler)
    
    log.Println("Starting API server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Testing Examples

### Unit Testing Saturn Jobs

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "strings"
    "testing"
)

type TestLogger struct {
    buffer *bytes.Buffer
}

func (tl *TestLogger) Info(msg string) {
    tl.buffer.WriteString("INFO: " + msg + "\n")
}

func (tl *TestLogger) Infof(format string, args ...interface{}) {
    tl.buffer.WriteString("INFO: " + fmt.Sprintf(format, args...) + "\n")
}

func (tl *TestLogger) Error(msg string) {
    tl.buffer.WriteString("ERROR: " + msg + "\n")
}

func (tl *TestLogger) Errorf(format string, args ...interface{}) {
    tl.buffer.WriteString("ERROR: " + fmt.Sprintf(format, args...) + "\n")
}

func (tl *TestLogger) Debug(msg string) {
    tl.buffer.WriteString("DEBUG: " + msg + "\n")
}

func (tl *TestLogger) Debugf(format string, args ...interface{}) {
    tl.buffer.WriteString("DEBUG: " + fmt.Sprintf(format, args...) + "\n")
}

func (tl *TestLogger) Warn(msg string) {
    tl.buffer.WriteString("WARN: " + msg + "\n")
}

func (tl *TestLogger) Warnf(format string, args ...interface{}) {
    tl.buffer.WriteString("WARN: " + fmt.Sprintf(format, args...) + "\n")
}

// Example test
func TestHelloJob(t *testing.T) {
    var logBuffer bytes.Buffer
    testLogger := &TestLogger{buffer: &logBuffer}
    
    // Create the job handler
    handler := func(args map[string]string, signature string) bool {
        name := args["name"]
        if name == "" {
            name = "World"
        }
        testLogger.Infof("Hello, %s! (Run: %s)", name, signature)
        return true
    }
    
    // Test with parameters
    args := map[string]string{"name": "TestUser"}
    signature := "test-signature"
    result := handler(args, signature)
    
    if !result {
        t.Error("Expected job to return true")
    }
    
    logOutput := logBuffer.String()
    if !strings.Contains(logOutput, "Hello, TestUser!") {
        t.Errorf("Expected log to contain 'Hello, TestUser!', got: %s", logOutput)
    }
    
    if !strings.Contains(logOutput, "test-signature") {
        t.Errorf("Expected log to contain signature, got: %s", logOutput)
    }
}

func TestStoppableJob(t *testing.T) {
    quit := make(chan struct{})
    done := make(chan bool, 1)
    
    handler := func(args map[string]string, signature string, quit chan struct{}) bool {
        counter := 0
        for {
            select {
            case <-quit:
                return true
            default:
                counter++
                if counter > 5 { // Limit iterations to avoid hanging tests
                    return true
                }
            }
        }
    }
    
    args := map[string]string{}
    signature := "test-stop"
    
    // Run the handler in a goroutine
    go func() {
        result := handler(args, signature, quit)
        done <- result
    }()
    
    // Allow some iterations
    time.Sleep(100 * time.Millisecond)
    
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

## Best Practices Demonstrated

These examples demonstrate several best practices:

1. **Error Handling**: Always validate inputs and handle errors gracefully
2. **Resource Management**: Properly manage file handles, connections, and other resources
3. **Cancellation Support**: Use quit channels in stoppable jobs for graceful shutdown
4. **Logging**: Use structured logging for debugging and monitoring
5. **Parameter Validation**: Validate parameters to prevent injection attacks
6. **Testing**: Write comprehensive tests for your Saturn jobs
7. **Documentation**: Document your custom job interfaces and expected parameters

## Running the Examples

To run any of these examples:

1. Create a new directory for the example
2. Create a `main.go` file with the example code
3. Run `go mod init example_name`
4. Run `go get github.com/Kingson4Wu/saturncli`
5. Build with `go build -o example_name main.go`
6. Create the necessary socket directory: `mkdir -p /tmp`
7. Run the server: `./example_name`
8. In another terminal, use the Saturn CLI to trigger jobs

## See Also

- [Getting Started Guide](./getting-started.md) - Basic setup and usage
- [Embedding Guide](./embedding.md) - How to integrate Saturn into your services
- [Client API Reference](./client-api.md) and [Server API Reference](./server-api.md) - Technical API documentation
- [Architecture](./architecture.md) - Understanding the system design