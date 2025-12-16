---
sidebar_position: 10
title: Best Practices - Optimal Saturn CLI Usage for Go Applications
description: Recommended practices for using Saturn CLI effectively, safely, and efficiently in your Go applications. Optimize job execution and inter-process communication.
keywords: [saturn cli best practices, go cli best practices, job execution best practices, unix domain sockets best practices, golang cli optimization]
---

# Best Practices

This guide provides recommended practices for using Saturn CLI effectively, safely, and efficiently in your Go applications for optimal job execution and process management.

## Job Design Patterns

### Creating Well-Behaved Jobs

Design your jobs to be predictable, reliable, and safe:

```go
func wellBehavedJob(args map[string]string, signature string) bool {
    // 1. Validate inputs first
    if args["required_param"] == "" {
        log.Printf("[%s] Missing required parameter", signature)
        return false
    }
    
    // 2. Handle resources properly with defer
    resource, err := obtainResource(args["resource_id"])
    if err != nil {
        log.Printf("[%s] Failed to obtain resource: %v", signature, err)
        return false
    }
    defer resource.Close() // Always cleanup
    
    // 3. Do the actual work
    success := doWork(resource, args)
    
    // 4. Log completion status
    if success {
        log.Printf("[%s] Job completed successfully", signature)
    } else {
        log.Printf("[%s] Job failed", signature)
    }
    
    return success
}
```

### Stoppable Job Patterns

For long-running jobs that need cancellation support:

```go
func goodStoppableJob(args map[string]string, signature string, quit chan struct{}) bool {
    log.Printf("[%s] Starting stoppable job", signature)
    
    // 1. Break work into small, interruptible chunks
    for i := 0; ; i++ {
        select {
        case <-quit:
            log.Printf("[%s] Job stopped at iteration %d", signature, i)
            return true // Return true for successful cancellation
        default:
            // 2. Do small amount of work
            if err := doSmallWorkChunk(i); err != nil {
                log.Printf("[%s] Error in work chunk %d: %v", signature, i, err)
                return false
            }
            
            // 3. Yield periodically to check quit channel
            time.Sleep(100 * time.Millisecond)
        }
    }
}

// Better version using ticker for more consistent intervals
func betterStoppableJob(args map[string]string, signature string, quit chan struct{}) bool {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for i := 0; ; i++ {
        select {
        case <-quit:
            log.Printf("[%s] Job stopped at iteration %d", signature, i)
            return true
        case <-ticker.C:
            // Do work for this interval
            if err := doWorkForInterval(i); err != nil {
                log.Printf("[%s] Error: %v", signature, err)
                return false
            }
        }
    }
}
```

## Parameter Handling

### Secure Parameter Processing

Always validate and sanitize parameters to prevent injection attacks:

```go
func secureJob(args map[string]string, signature string) bool {
    // 1. Validate all parameters
    if !validateParameters(args) {
        log.Printf("[%s] Invalid parameters provided", signature)
        return false
    }
    
    // 2. Sanitize parameters before use
    filename := sanitizeFilename(args["filename"])
    if filename == "" {
        log.Printf("[%s] Invalid filename after sanitization", signature)
        return false
    }
    
    // 3. Use parameters safely
    data, err := readSafeFile(filename)
    if err != nil {
        log.Printf("[%s] Error reading file: %v", signature, err)
        return false
    }
    
    return processData(data)
}

func validateParameters(args map[string]string) bool {
    // Check required parameters
    required := []string{"user_id", "action"}
    for _, param := range required {
        if args[param] == "" {
            return false
        }
    }
    
    // Validate parameter formats
    if !isValidUserID(args["user_id"]) {
        return false
    }
    
    return true
}

func sanitizeFilename(filename string) string {
    // Prevent directory traversal
    if strings.Contains(filename, "../") || strings.Contains(filename, "..\\") {
        return ""
    }
    
    // Only allow safe characters
    if !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(filename) {
        return ""
    }
    
    return filepath.Clean(filename)
}
```

## Error Handling and Resilience

### Comprehensive Error Handling

Structure your jobs with proper error handling at every level:

```go
func resilientJob(args map[string]string, signature string, quit chan struct{}) bool {
    // 1. Input validation
    if err := validateJobArgs(args); err != nil {
        log.Printf("[%s] Input validation failed: %v", signature, err)
        return false
    }
    
    // 2. Resource acquisition with retry
    var resource *MyResource
    for attempt := 0; attempt < 3; attempt++ {
        select {
        case <-quit:
            return true
        default:
            var err error
            resource, err = acquireResource(args["resource_id"])
            if err == nil {
                break
            }
            log.Printf("[%s] Resource acquisition failed (attempt %d): %v", 
                      signature, attempt+1, err)
            time.Sleep(time.Duration(attempt+1) * time.Second)
        }
    }
    
    if resource == nil {
        log.Printf("[%s] Failed to acquire resource after retries", signature)
        return false
    }
    defer resource.Close()
    
    // 3. Main work with progress tracking
    totalWork := extractTotalWork(args)
    for i := 0; i < totalWork; i++ {
        select {
        case <-quit:
            log.Printf("[%s] Job cancelled at %d/%d", signature, i, totalWork)
            return true
        default:
            if err := doWorkUnit(resource, i); err != nil {
                log.Printf("[%s] Work unit %d failed: %v", signature, i, err)
                
                // Decide whether to continue or fail based on error type
                if isFatalError(err) {
                    return false
                }
                // For non-fatal errors, continue
            }
            
            // Track progress
            if i%10 == 0 { // Log every 10 units
                log.Printf("[%s] Progress: %d/%d", signature, i, totalWork)
            }
        }
    }
    
    log.Printf("[%s] Job completed successfully", signature)
    return true
}

func isFatalError(err error) bool {
    // Define which errors are fatal vs retryable
    return errors.Is(err, ErrInvalidData) || errors.Is(err, ErrPermissionDenied)
}
```

## Performance Optimization

### Efficient Resource Usage

Optimize for memory, CPU, and I/O efficiency:

```go
func efficientJob(args map[string]string, signature string, quit chan struct{}) bool {
    // 1. Pre-allocate slices when possible
    expectedCount := extractExpectedCount(args)
    results := make([]string, 0, expectedCount)
    
    // 2. Process in batches to limit memory usage
    batchSize := 100
    for start := 0; start < expectedCount; start += batchSize {
        select {
        case <-quit:
            return true
        default:
            // Process batch
            batchEnd := start + batchSize
            if batchEnd > expectedCount {
                batchEnd = expectedCount
            }
            
            batchResults, err := processBatch(start, batchEnd)
            if err != nil {
                log.Printf("[%s] Batch %d-%d failed: %v", signature, start, batchEnd, err)
                return false
            }
            
            results = append(results, batchResults...)
        }
    }
    
    // 3. Process results efficiently
    return saveResults(results)
}

// 4. Use buffered I/O for file operations
func processFileWithBuffering(inputPath, outputPath string) error {
    input, err := os.Open(inputPath)
    if err != nil {
        return err
    }
    defer input.Close()
    
    output, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer output.Close()
    
    // Use buffered reader/writer for better performance
    reader := bufio.NewReaderSize(input, 64*1024) // 64KB buffer
    writer := bufio.NewWriterSize(output, 64*1024) // 64KB buffer
    
    // Process data
    buffer := make([]byte, 8192) // 8KB processing buffer
    for {
        n, err := reader.Read(buffer)
        if n > 0 {
            if _, writeErr := writer.Write(buffer[:n]); writeErr != nil {
                return writeErr
            }
        }
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
    }
    
    return writer.Flush()
}
```

## Security Considerations

### Secure Job Implementation

Implement security best practices to protect your system:

```go
func secureJobWithValidation(args map[string]string, signature string) bool {
    // 1. Validate and sanitize all inputs
    userID := args["user_id"]
    if !isValidUserID(userID) {
        return false
    }
    
    action := args["action"]
    if !isValidAction(action) {
        return false
    }
    
    // 2. Implement access controls
    if !userHasPermission(userID, action) {
        log.Printf("[%s] Access denied for user %s to perform %s", signature, userID, action)
        return false
    }
    
    // 3. Use system calls safely (if needed)
    if action == "execute_command" {
        command := args["command"]
        if !isValidCommand(command) {
            log.Printf("[%s] Invalid command: %s", signature, command)
            return false
        }
        
        // Use only safe commands, never pass user input directly
        cmd := exec.Command("/usr/bin/safe-utility", sanitizeCommandArgs(command))
        output, err := cmd.CombinedOutput()
        
        if err != nil {
            log.Printf("[%s] Command failed: %v, output: %s", signature, err, output)
            return false
        }
        
        log.Printf("[%s] Command completed successfully", signature)
    }
    
    return true
}

func isValidCommand(cmd string) bool {
    // Only allow specific, safe commands
    allowedCommands := map[string]bool{
        "status": true,
        "backup": true,
        "report": true,
    }
    
    return allowedCommands[cmd]
}

func sanitizeCommandArgs(args string) []string {
    // Validate that command args don't contain dangerous shell characters
    dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "*"}
    for _, char := range dangerousChars {
        if strings.Contains(args, char) {
            return nil // Invalid args
        }
    }
    
    // Split and sanitize individual arguments
    parts := strings.Fields(args)
    sanitized := make([]string, 0, len(parts))
    
    for _, part := range parts {
        // Apply further sanitization as needed
        sanitized = append(sanitized, strings.TrimSpace(part))
    }
    
    return sanitized
}
```

## Logging and Monitoring

### Effective Logging Practices

Implement comprehensive and useful logging:

```go
func wellLoggedJob(args map[string]string, signature string, quit chan struct{}) bool {
    start := time.Now()
    log.Printf("[%s] Job started with args: %+v", signature, redactSensitiveInfo(args))
    
    defer func() {
        duration := time.Since(start)
        log.Printf("[%s] Job completed in %v", signature, duration)
    }()
    
    // Use structured logging with context
    log.Printf("[%s] Step 1: Initializing resources", signature)
    resource, err := initializeResource()
    if err != nil {
        log.Printf("[%s] ERROR: Failed to initialize resource: %v", signature, err)
        return false
    }
    defer func() {
        if closeErr := resource.Close(); closeErr != nil {
            log.Printf("[%s] WARNING: Error closing resource: %v", signature, closeErr)
        }
    }()
    
    log.Printf("[%s] Step 2: Processing main work", signature)
    for i := 0; i < 100; i++ {
        select {
        case <-quit:
            log.Printf("[%s] Job cancelled at step %d/100", signature, i)
            return true
        default:
            if i%20 == 0 { // Log progress every 20%
                progress := float64(i) / 100.0 * 100
                log.Printf("[%s] Progress: %.1f%%", signature, progress)
            }
            
            if err := processItem(i, resource); err != nil {
                log.Printf("[%s] WARNING: Error processing item %d: %v", signature, i, err)
                // Decide whether to continue or fail based on context
            }
        }
    }
    
    log.Printf("[%s] Step 3: Finalizing", signature)
    return finalizeWork(resource)
}

func redactSensitiveInfo(args map[string]string) map[string]string {
    // Don't log sensitive information
    safeArgs := make(map[string]string, len(args))
    for k, v := range args {
        switch k {
        case "password", "token", "secret", "key":
            safeArgs[k] = "[REDACTED]"
        default:
            safeArgs[k] = v
        }
    }
    return safeArgs
}
```

## Testing and Verification

### Writing Testable Jobs

Design jobs that can be easily tested and verified:

```go
// Interface for dependencies to enable mocking
type ResourceProvider interface {
    GetResource(id string) (Resource, error)
    Close() error
}

// Make your job function configurable for testing
type JobConfig struct {
    ResourceProvider ResourceProvider
    Logger           utils.Logger
    MaxRetries       int
}

func jobWithConfig(config JobConfig) server.JobHandler {
    return func(args map[string]string, signature string) bool {
        resource, err := config.ResourceProvider.GetResource(args["resource_id"])
        if err != nil {
            config.Logger.Errorf("[%s] Failed to get resource: %v", signature, err)
            return false
        }
        defer config.ResourceProvider.Close()
        
        // ... rest of job logic using config
        return doWorkWithResource(resource)
    }
}

// Example unit test
func TestJobWithConfig(t *testing.T) {
    // Create mock provider
    mockProvider := &MockResourceProvider{
        resource: &TestResource{},
        err:      nil,
    }
    
    config := JobConfig{
        ResourceProvider: mockProvider,
        Logger:           &TestLogger{},
        MaxRetries:       3,
    }
    
    jobHandler := jobWithConfig(config)
    
    result := jobHandler(map[string]string{"resource_id": "test"}, "test-sig")
    
    if !result {
        t.Error("Expected job to succeed")
    }
}
```

## Integration Patterns

### Service Integration Best Practices

When embedding Saturn CLI in services:

```go
type JobService struct {
    server      *server.ser
    logger      utils.Logger
    client      *client.cli
    initialized bool
}

func NewJobService(config Config) *JobService {
    service := &JobService{
        logger: config.Logger,
    }
    
    // Register all jobs during initialization
    if err := service.registerJobs(); err != nil {
        service.logger.Errorf("Failed to register jobs: %v", err)
        return nil
    }
    
    // Create server with proper registry
    service.server = server.NewServer(
        service.logger,
        config.SocketPath,
        server.WithRegistry(config.Registry),
    )
    
    // Create client for internal use
    service.client = client.NewClient(service.logger, config.SocketPath)
    
    service.initialized = true
    return service
}

func (js *JobService) registerJobs() error {
    jobs := map[string]interface{}{
        "data-processor":  js.handleDataProcessing,
        "file-archiver":   js.handleFileArchival,
        "report-generator": js.handleReportGeneration,
    }
    
    for name, jobFunc := range jobs {
        if handler, ok := jobFunc.(server.JobHandler); ok {
            if err := server.AddJob(name, handler); err != nil {
                return fmt.Errorf("failed to register job %s: %w", name, err)
            }
        } else if handler, ok := jobFunc.(server.StoppableJobHandler); ok {
            if err := server.AddStoppableJob(name, handler); err != nil {
                return fmt.Errorf("failed to register stoppable job %s: %w", name, err)
            }
        }
    }
    
    return nil
}

func (js *JobService) HealthCheck() bool {
    // Verify the service is running correctly
    result := js.client.Run(&client.Task{
        Name: "health-check",
        Params: map[string]string{"timestamp": time.Now().String()},
    })
    
    return result == base.SUCCESS
}
```

## Common Anti-Patterns to Avoid

### What NOT to Do

```go
// BAD: Resource leak
func badResourceLeak(args map[string]string, signature string) bool {
    file, _ := os.Open("important_file.txt") // No error handling
    // Do work with file...
    // FORGOT to close file! - Resource leak!
    return true
}

// GOOD: Proper resource handling
func goodResourceHandling(args map[string]string, signature string) bool {
    file, err := os.Open("important_file.txt")
    if err != nil {
        return false
    }
    defer file.Close() // Proper cleanup
    
    // Do work with file...
    return true
}

// BAD: Blocking without checking quit channel
func badBlockingJob(args map[string]string, signature string, quit chan struct{}) bool {
    time.Sleep(10 * time.Hour) // Blocks for 10 hours without checking quit!
    return true
}

// GOOD: Non-blocking with quit channel check
func goodNonBlockingJob(args map[string]string, signature string, quit chan struct{}) bool {
    for i := 0; i < 1000; i++ {
        select {
        case <-quit:
            return true
        case <-time.After(10 * time.Millisecond):
            // Check quit channel every 10ms
        }
        // Do small bit of work
    }
    return true
}

// BAD: No parameter validation
func badNoValidation(args map[string]string, signature string) bool {
    // Directly passing user input to system commands - SECURITY RISK!
    cmd := exec.Command("sh", "-c", args["command"])
    cmd.Run()
    return true
}

// GOOD: Input validation
func goodWithValidation(args map[string]string, signature string) bool {
    // Validate input before using
    if !isValidCommand(args["command"]) {
        return false
    }
    
    // Sanitize and use safely
    sanitizedCmd := sanitizeCommand(args["command"])
    cmd := exec.Command("/safe/whitelisted/command", sanitizedCmd)
    return cmd.Run() == nil
}
```

## Performance Monitoring

### Observability Best Practices

Include proper monitoring and metrics in your Saturn CLI implementations:

```go
type JobMetrics struct {
    totalJobs     int64
    successfulJobs int64
    failedJobs    int64
    totalDuration time.Duration
}

func monitoredJob(metrics *JobMetrics, args map[string]string, signature string) bool {
    start := time.Now()
    success := false
    
    defer func() {
        duration := time.Since(start)
        
        atomic.AddInt64(&metrics.totalJobs, 1)
        atomic.AddInt64(&metrics.totalDuration, duration.Nanoseconds())
        
        if success {
            atomic.AddInt64(&metrics.successfulJobs, 1)
        } else {
            atomic.AddInt64(&metrics.failedJobs, 1)
        }
    }()
    
    // Your job logic here
    success = performJob(args, signature)
    
    return success
}

func (jm *JobMetrics) GetStats() map[string]interface{} {
    total := atomic.LoadInt64(&jm.totalJobs)
    if total == 0 {
        return map[string]interface{}{"message": "no jobs executed yet"}
    }
    
    return map[string]interface{}{
        "total_jobs":     total,
        "successful":     atomic.LoadInt64(&jm.successfulJobs),
        "failed":         atomic.LoadInt64(&jm.failedJobs),
        "success_rate":   float64(atomic.LoadInt64(&jm.successfulJobs)) / float64(total) * 100,
        "avg_duration":   time.Duration(atomic.LoadInt64(&jm.totalDuration)/total) * time.Nanosecond,
    }
}
```

## Next Steps

- Review the [Troubleshooting](./troubleshooting.md) guide for common issue resolution
- Look at the [Examples](./examples.md) for practical implementations
- Check the [Architecture](./architecture.md) documentation for system design understanding
- Follow the [Contributing](./contributing.md) guidelines if you plan to modify Saturn CLI