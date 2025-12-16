---
sidebar_position: 6
title: Embedding Guide - Integrate Saturn CLI into Go Services
description: Complete guide on embedding Saturn CLI into existing Go services. Learn integration patterns and best practices for seamless job execution in your applications.
keywords: [saturn cli embedding, embed go cli, integrate saturn cli, go service integration, embed job execution, saturn cli integration]
---

# Embedding Guide

Saturn CLI is designed to be embedded directly into your existing services. This guide provides patterns and best practices for seamless integration of job execution capabilities into your Go applications.

## Why Embed Saturn CLI?

Saturn CLI excels at scenarios where you need to:
- Execute shell-style commands as background jobs
- Integrate external tools with your Go services
- Provide job scheduling capabilities to your applications
- Enable graceful cancellation of long-running tasks
- Maintain fast, secure communication channels with background workers

## Basic Integration Pattern

Here's the recommended pattern for embedding Saturn CLI into your services:

### 1. Service Integration Template

```go
package main

import (
    "context"
    "log"
    "sync"

    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

// MyService demonstrates how to embed Saturn CLI into an existing service
type MyService struct {
    server   *server.ser
    registry *server.Registry
    logger   utils.Logger
    wg       sync.WaitGroup
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewMyService(logger utils.Logger) *MyService {
    ctx, cancel := context.WithCancel(context.Background())
    
    service := &MyService{
        logger: logger,
        ctx:    ctx,
        cancel: cancel,
    }
    
    // Create a dedicated registry for this service
    service.registry = server.NewRegistry()
    
    // Register jobs specific to this service
    service.setupJobs()
    
    // Initialize the Saturn server
    service.server = server.NewServer(
        logger,
        "/tmp/myservice_saturn.sock",  // Customize socket path
        server.WithRegistry(service.registry),
    )
    
    return service
}

func (s *MyService) setupJobs() {
    // Register service-specific jobs
    if err := s.registry.AddJob("data-processing", s.handleDataProcessing); err != nil {
        s.logger.Errorf("Failed to register data-processing job: %v", err)
    }
    
    if err := s.registry.AddStoppableJob("continuous-monitoring", s.handleContinuousMonitoring); err != nil {
        s.logger.Errorf("Failed to register continuous-monitoring job: %v", err)
    }
}

func (s *MyService) handleDataProcessing(args map[string]string, signature string) bool {
    s.logger.Infof("Starting data processing job %s with args: %v", signature, args)
    
    // Extract and validate parameters
    dataset := args["dataset"]
    if dataset == "" {
        s.logger.Errorf("Missing dataset parameter for job %s", signature)
        return false
    }
    
    // Simulate data processing
    s.logger.Infof("Processing dataset: %s in job %s", dataset, signature)
    
    // Add your actual data processing logic here
    // ...
    
    s.logger.Infof("Data processing job %s completed", signature)
    return true
}

func (s *MyService) handleContinuousMonitoring(args map[string]string, signature string, quit chan struct{}) bool {
    s.logger.Infof("Starting continuous monitoring job %s", signature)
    
    intervalStr := args["interval"]
    if intervalStr == "" {
        intervalStr = "5" // default to 5 seconds
    }
    
    // Parse interval (implement proper parsing for production use)
    // interval, _ := strconv.Atoi(intervalStr)
    
    // Monitoring loop that can be stopped via quit channel
    ticker := s.newTickerWithQuit(quit, 5) // every 5 seconds
    
    for {
        select {
        case <-quit:
            s.logger.Infof("Monitoring job %s received quit signal", signature)
            return true
        case <-ticker.C:
            s.logger.Debugf("Performing monitoring check in job %s", signature)
            // Perform monitoring activity
            s.performHealthCheck()
        }
    }
}

func (s *MyService) newTickerWithQuit(quit chan struct{}, seconds int) *tickerWithQuit {
    t := time.NewTicker(time.Duration(seconds) * time.Second)
    return &tickerWithQuit{t: t, quit: quit}
}

type tickerWithQuit struct {
    t    *time.Ticker
    quit chan struct{}
}

func (twq *tickerWithQuit) C() <-chan time.Time {
    return twq.t.C
}

func (s *MyService) performHealthCheck() {
    // Implement your health checking logic here
    s.logger.Debugf("Health check performed")
}

func (s *MyService) Start() {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.logger.Info("Starting Saturn server...")
        s.server.Serve()
    }()
}

func (s *MyService) Stop() {
    s.logger.Info("Stopping Saturn server...")
    s.cancel() // Cancel context
    s.wg.Wait() // Wait for goroutines to finish
}

// Client usage example
func (s *MyService) TriggerDataProcessing(dataset string) error {
    // Create a client to communicate with our own server
    cli := client.NewClient(s.logger, "/tmp/myservice_saturn.sock")
    
    result := cli.Run(&client.Task{
        Name: "data-processing",
        Params: map[string]string{
            "dataset": dataset,
            "source":  "api",
        },
    })
    
    switch result {
    case base.SUCCESS:
        s.logger.Info("Data processing job completed successfully")
        return nil
    case base.FAILURE:
        return errors.New("data processing job failed")
    case base.INTERRUPT:
        s.logger.Info("Data processing job was interrupted")
        return nil
    default:
        return errors.New("unknown job result")
    }
}

func main() {
    service := NewMyService(&utils.DefaultLogger{})
    
    // Start the service
    service.Start()
    
    // The server is now running and accepting job requests
    // You can trigger jobs via client or external commands
    
    // Graceful shutdown example with signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigChan // Wait for signal
    log.Println("Shutting down service...")
    
    service.Stop()
}
```

## Advanced Integration Patterns

### Pattern 1: Job Factory with Configuration

```go
type JobFactory struct {
    config Config
    logger utils.Logger
}

type Config struct {
    Timeout    time.Duration
    MaxRetries int
    QueueSize  int
}

func (jf *JobFactory) CreateDatabaseBackupJob(dbConfig DBConfig) server.JobHandler {
    return func(args map[string]string, signature string) bool {
        jf.logger.Infof("Starting database backup %s", signature)
        
        ctx, cancel := context.WithTimeout(context.Background(), jf.config.Timeout)
        defer cancel()
        
        // Use dbConfig to perform database backup
        // ...
        
        // Log progress and handle errors appropriately
        jf.logger.Infof("Database backup %s completed", signature)
        return true
    }
}

func (jf *JobFactory) CreateFileProcessorJob(outputDir string) server.StoppableJobHandler {
    return func(args map[string]string, signature string, quit chan struct{}) bool {
        jf.logger.Infof("Starting file processor %s", signature)
        
        // Process files with ability to quit
        files := strings.Split(args["files"], ",")
        
        for i, file := range files {
            select {
            case <-quit:
                jf.logger.Infof("File processor %s stopped at file %d", signature, i)
                return true
            default:
                // Process individual file
                jf.processFile(file, signature)
            }
        }
        
        jf.logger.Infof("File processor %s completed", signature)
        return true
    }
}
```

### Pattern 2: Service Container with Multiple Saturn Instances

```go
type ServiceContainer struct {
    apiServer     *http.Server
    jobServer     *server.ser
    notificationSrv *server.ser
    logger        utils.Logger
    cancel        context.CancelFunc
}

func NewServiceContainer(config ServiceConfig) *ServiceContainer {
    ctx, cancel := context.WithCancel(context.Background())
    
    container := &ServiceContainer{
        logger: &utils.DefaultLogger{},
        cancel: cancel,
    }
    
    // Create separate Saturn servers for different purposes
    container.jobServer = container.createJobServer(config.JobSocketPath)
    container.notificationSrv = container.createNotificationServer(config.NotificationSocketPath)
    
    return container
}

func (sc *ServiceContainer) createJobServer(socketPath string) *server.ser {
    registry := server.NewRegistry()
    
    // Register job-specific handlers
    registry.AddJob("batch-process", sc.handleBatchProcess)
    registry.AddStoppableJob("monitor-resources", sc.handleResourceMonitoring)
    
    return server.NewServer(sc.logger, socketPath, server.WithRegistry(registry))
}

func (sc *ServiceContainer) createNotificationServer(socketPath string) *server.ser {
    registry := server.NewRegistry()
    
    // Register notification-specific handlers
    registry.AddJob("send-email", sc.handleSendEmail)
    registry.AddJob("push-notification", sc.handlePushNotification)
    
    return server.NewServer(sc.logger, socketPath, server.WithRegistry(registry))
}

func (sc *ServiceContainer) Start() {
    go sc.jobServer.Serve()
    go sc.notificationSrv.Serve()
    // Start other services...
}

func (sc *ServiceContainer) Stop() {
    sc.cancel()
    // Close servers gracefully
}
```

## Configuration Management

### Environment-Based Configuration

```go
type SaturnConfig struct {
    SocketPath    string
    JobTimeout    time.Duration
    MaxWorkers    int
    EnableMetrics bool
}

func LoadSaturnConfigFromEnv() SaturnConfig {
    return SaturnConfig{
        SocketPath:    getEnvOrDefault("SATURN_SOCKET_PATH", "/tmp/saturn.sock"),
        JobTimeout:    parseDurationEnv("SATURN_JOB_TIMEOUT", 5*time.Minute),
        MaxWorkers:    parseIntEnv("SATURN_MAX_WORKERS", 10),
        EnableMetrics: parseBoolEnv("SATURN_ENABLE_METRICS", true),
    }
}

func ApplyConfigToServer(server *server.ser, config SaturnConfig) {
    // Apply configuration to server - this is conceptual,
    // real implementation would depend on Saturn's configuration API
}
```

## Logging and Monitoring

### Custom Logger Integration

```go
// Example of integrating with a structured logger like logrus
type LogrusAdapter struct {
    logger *logrus.Logger
}

func (l *LogrusAdapter) Info(msg string) {
    l.logger.Info(msg)
}

func (l *LogrusAdapter) Infof(format string, args ...interface{}) {
    l.logger.Infof(format, args...)
}

func (l *LogrusAdapter) Error(msg string) {
    l.logger.Error(msg)
}

func (l *LogrusAdapter) Errorf(format string, args ...interface{}) {
    l.logger.Errorf(format, args...)
}

func (l *LogrusAdapter) Warn(msg string) {
    l.logger.Warn(msg)
}

func (l *LogrusAdapter) Warnf(format string, args ...interface{}) {
    l.logger.Warnf(format, args...)
}

func (l *LogrusAdapter) Debug(msg string) {
    l.logger.Debug(msg)
}

func (l *LogrusAdapter) Debugf(format string, args ...interface{}) {
    l.logger.Debugf(format, args...)
}
```

## Testing Embedded Saturn

### Unit Testing Job Handlers

```go
func TestMyJobHandler(t *testing.T) {
    // Mock logger for testing
    var logBuffer bytes.Buffer
    mockLogger := &TestLogger{writer: &logBuffer}
    
    // Test data
    args := map[string]string{
        "input": "test_data",
        "param": "value",
    }
    
    // Create handler
    handler := createMyJobHandler(mockLogger)
    
    // Execute handler
    result := handler(args, "test_signature")
    
    // Assertions
    if !result {
        t.Error("Expected job to succeed")
    }
    
    // Verify log output
    if !strings.Contains(logBuffer.String(), "expected log message") {
        t.Error("Expected log message not found")
    }
}

func TestStoppableJobHandler(t *testing.T) {
    // Similar test for stoppable job with quit channel
    quit := make(chan struct{})
    args := map[string]string{"test": "value"}
    
    // Use a timeout to avoid hanging tests
    done := make(chan bool)
    
    go func() {
        result := myStoppableHandler(args, "test_sig", quit)
        done <- result
    }()
    
    // Close quit channel after a delay to test stopping behavior
    time.Sleep(100 * time.Millisecond)
    close(quit)
    
    select {
    case result := <-done:
        // Verify that the job handled the quit signal correctly
        if !result {
            t.Error("Expected job to return true after stopping")
        }
    case <-time.After(1 * time.Second):
        t.Error("Test timed out - job didn't respond to quit signal")
    }
}
```

## Security Considerations

### Socket Permissions

When deploying in production, consider the following security aspects:

1. **Socket File Permissions**: Restrict socket file access to authorized processes only
2. **Parameter Validation**: Always validate parameters passed to job handlers to prevent injection attacks
3. **Timeout Enforcement**: Implement timeouts in job handlers to prevent resource exhaustion
4. **Resource Limits**: Limit resource consumption (memory, disk, CPU) within job handlers

### Parameter Validation Example

```go
func validateJobParameters(args map[string]string) error {
    // Validate that all required params are present
    for _, required := range []string{"input", "output"} {
        if args[required] == "" {
            return fmt.Errorf("missing required parameter: %s", required)
        }
    }
    
    // Sanitize input parameters to prevent command injection
    for key, value := range args {
        if strings.Contains(value, "../") || strings.Contains(value, "..\\") {
            return fmt.Errorf("unsafe path parameter detected in %s: %s", key, value)
        }
        
        // Additional validation based on parameter meaning
        // ...
    }
    
    return nil
}
```

## Performance Optimization

### Connection Pooling for Clients

```go
type ClientPool struct {
    mu    sync.RWMutex
    pool  []*client.cli
    max   int
    count int
    factory func() *client.cli
}

func NewClientPool(maxClients int, factory func() *client.cli) *ClientPool {
    return &ClientPool{
        max:     maxClients,
        factory: factory,
    }
}

func (cp *ClientPool) GetClient() *client.cli {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    if cp.count < cp.max && len(cp.pool) == 0 {
        cp.count++
        return cp.factory()
    }
    
    if len(cp.pool) > 0 {
        client := cp.pool[len(cp.pool)-1]
        cp.pool = cp.pool[:len(cp.pool)-1]
        return client
    }
    
    // Block until a client becomes available
    // In practice, you might want to implement a timeout
    for len(cp.pool) == 0 {
        cp.mu.Unlock()
        time.Sleep(10 * time.Millisecond)
        cp.mu.Lock()
    }
    
    client := cp.pool[len(cp.pool)-1]
    cp.pool = cp.pool[:len(cp.pool)-1]
    return client
}

func (cp *ClientPool) ReturnClient(cli *client.cli) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    
    if len(cp.pool) < cp.max {
        cp.pool = append(cp.pool, cli)
    } else {
        // Pool is full, discard the extra client
        cp.count--
    }
}
```

## Migration Strategies

### Version Management

When upgrading Saturn CLI in your services:

1. **Deploy Gradually**: Roll out Saturn updates to subsets of your service fleet
2. **Maintain Compatibility**: Ensure older clients can communicate with newer servers during transition periods
3. **Monitor Metrics**: Track job execution times, success rates, and error patterns
4. **Version Pinning**: In production environments, pin specific Saturn versions to ensure stability

## Troubleshooting

### Common Issues and Solutions

1. **Socket Permission Errors**: Ensure your service has appropriate permissions to create socket files in the specified directory
2. **Job Timeouts**: Increase timeout values or optimize job implementation
3. **Memory Leaks**: Monitor for unreleased resources in job handlers
4. **Connection Failures**: Verify socket paths and ensure server is running

## Next Steps

- Review the [Architecture](./architecture.md) documentation to understand system design
- Check out [Examples](./examples.md) for real-world usage patterns
- Consult the [Client API](./client-api.md) and [Server API](./server-api.md) references for detailed method information