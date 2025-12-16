---
sidebar_position: 12
title: Testing Strategies - Saturn CLI Job and Server Testing
description: Comprehensive testing strategies for Saturn CLI. Learn how to test your own jobs and verify Saturn CLI functionality for reliable job execution.
keywords: [saturn cli testing, go cli testing, job execution testing, saturn cli tests, go background job testing, cli testing strategies]
---

# Testing

This guide covers comprehensive testing strategies for Saturn CLI - both testing your own jobs and testing Saturn CLI itself for reliable job execution.

## Testing Overview

Saturn CLI follows Go's standard testing practices with additional patterns for testing concurrent and network-dependent components.

### Testing Philosophy

- **Fast Feedback**: Tests should run quickly to provide immediate feedback
- **Reliability**: Tests should be deterministic and not flaky
- **Maintainability**: Tests should be easy to understand and modify
- **Coverage**: Aim for high test coverage, especially around critical paths
- **Integration**: Test both individual components and their integration

## Unit Testing Job Handlers

### Testing Regular Jobs

```go
package main

import (
    "bytes"
    "fmt"
    "strings"
    "testing"
)

// Example job handler to test
func exampleJob(args map[string]string, signature string) bool {
    if args["fail"] == "true" {
        return false
    }
    
    expected := args["expected"]
    actual := args["actual"]
    
    return expected == actual
}

func TestExampleJob(t *testing.T) {
    tests := []struct {
        name     string
        args     map[string]string
        expected bool
    }{
        {
            name: "matching values return true",
            args: map[string]string{
                "expected": "value",
                "actual":   "value",
            },
            expected: true,
        },
        {
            name: "different values return false",
            args: map[string]string{
                "expected": "value1",
                "actual":   "value2",
            },
            expected: false,
        },
        {
            name: "fail flag returns false",
            args: map[string]string{
                "fail": "true",
            },
            expected: false,
        },
        {
            name:     "empty args return false",
            args:     map[string]string{},
            expected: true, // because expected and actual are both empty strings
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := exampleJob(tt.args, "test-signature")
            if result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Testing Stoppable Jobs

Testing stoppable jobs requires special handling for the quit channel:

```go
package main

import (
    "sync"
    "testing"
    "time"
)

func exampleStoppableJob(args map[string]string, signature string, quit chan struct{}) bool {
    maxIterations := 10
    if args["max"] != "" {
        fmt.Sscanf(args["max"], "%d", &maxIterations)
    }
    
    for i := 0; i < maxIterations; i++ {
        select {
        case <-quit:
            return true // Successfully stopped
        case <-time.After(10 * time.Millisecond):
            // Continue working
        }
    }
    
    return true // Completed normally
}

func TestStoppableJobStop(t *testing.T) {
    quit := make(chan struct{})
    done := make(chan bool, 1)
    
    args := map[string]string{"max": "100"} // High number to ensure it would run long
    
    // Run the job in a goroutine
    go func() {
        result := exampleStoppableJob(args, "test-sig", quit)
        done <- result
    }()
    
    // Allow some time for the job to start
    time.Sleep(50 * time.Millisecond)
    
    // Signal to quit
    close(quit)
    
    // Wait for job to complete with timeout
    select {
    case result := <-done:
        if !result {
            t.Error("Expected job to return true when stopped")
        }
    case <-time.After(2 * time.Second):
        t.Error("Job did not respond to quit signal in time")
    }
}

func TestStoppableJobComplete(t *testing.T) {
    quit := make(chan struct{})
    done := make(chan bool, 1)
    
    args := map[string]string{"max": "3"} // Small number to complete quickly
    
    go func() {
        result := exampleStoppableJob(args, "test-sig", quit)
        done <- result
    }()
    
    // Wait for job to complete with timeout
    select {
    case result := <-done:
        if !result {
            t.Error("Expected job to return true when completed normally")
        }
    case <-time.After(2 * time.Second):
        t.Error("Job did not complete in time")
    }
}
```

## Testing with Mock Dependencies

### Mocking External Dependencies

```go
package main

import (
    "fmt"
    "testing"
)

// Interface to make code testable
type ResourceService interface {
    GetData(id string) (string, error)
    SaveData(id, data string) error
}

// Real implementation
type RealResourceService struct{}

func (r *RealResourceService) GetData(id string) (string, error) {
    // Real implementation that might access database, API, etc.
    return fmt.Sprintf("data-for-%s", id), nil
}

func (r *RealResourceService) SaveData(id, data string) error {
    // Real implementation
    return nil
}

// Job that uses external service
func jobWithDependency(service ResourceService, args map[string]string, signature string) bool {
    id := args["id"]
    if id == "" {
        return false
    }
    
    data, err := service.GetData(id)
    if err != nil {
        return false
    }
    
    return service.SaveData(id, data) == nil
}

// Mock implementation for testing
type MockResourceService struct {
    GetDataFunc func(id string) (string, error)
    SaveDataFunc func(id, data string) error
}

func (m *MockResourceService) GetData(id string) (string, error) {
    if m.GetDataFunc != nil {
        return m.GetDataFunc(id)
    }
    return "", nil
}

func (m *MockResourceService) SaveData(id, data string) error {
    if m.SaveDataFunc != nil {
        return m.SaveDataFunc(id, data)
    }
    return nil
}

func TestJobWithDependencySuccess(t *testing.T) {
    mockService := &MockResourceService{
        GetDataFunc: func(id string) (string, error) {
            if id == "valid-id" {
                return "test-data", nil
            }
            return "", fmt.Errorf("not found")
        },
        SaveDataFunc: func(id, data string) error {
            if id == "valid-id" && data == "test-data" {
                return nil
            }
            return fmt.Errorf("save failed")
        },
    }
    
    args := map[string]string{"id": "valid-id"}
    result := jobWithDependency(mockService, args, "test-sig")
    
    if !result {
        t.Error("Expected job to succeed with valid dependencies")
    }
}

func TestJobWithDependencyFailure(t *testing.T) {
    mockService := &MockResourceService{
        GetDataFunc: func(id string) (string, error) {
            return "", fmt.Errorf("service unavailable")
        },
    }
    
    args := map[string]string{"id": "any-id"}
    result := jobWithDependency(mockService, args, "test-sig")
    
    if result {
        t.Error("Expected job to fail when dependency fails")
    }
}
```

## Integration Testing

### Testing Client-Server Integration

```go
package main

import (
    "sync"
    "testing"
    "time"

    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/base"
    "github.com/Kingson4Wu/saturncli/utils"
)

func TestClientServerIntegration(t *testing.T) {
    socketPath := "/tmp/test_integration.sock"
    
    // Register a test job
    testJobExecuted := false
    testJobMutex := sync.Mutex{}
    
    jobHandler := func(args map[string]string, signature string) bool {
        testJobMutex.Lock()
        testJobExecuted = true
        testJobMutex.Unlock()
        return true
    }
    
    if err := server.AddJob("integration-test", jobHandler); err != nil {
        t.Fatalf("Failed to register test job: %v", err)
    }
    
    // Start server in goroutine
    serverStarted := make(chan bool)
    go func() {
        serverStarted <- true
        server.NewServer(&utils.DefaultLogger{}, socketPath).Serve()
    }()
    
    // Give server time to start
    <-serverStarted
    time.Sleep(100 * time.Millisecond)
    
    // Create client and run job
    cli := client.NewClient(&utils.DefaultLogger{}, socketPath)
    result := cli.Run(&client.Task{
        Name: "integration-test",
        Params: map[string]string{
            "test": "value",
        },
    })
    
    if result != base.SUCCESS {
        t.Errorf("Expected SUCCESS, got %v", result)
    }
    
    // Verify job was executed
    testJobMutex.Lock()
    executed := testJobExecuted
    testJobMutex.Unlock()
    
    if !executed {
        t.Error("Expected test job to be executed")
    }
}

func TestStoppableJobIntegration(t *testing.T) {
    socketPath := "/tmp/test_stoppable_integration.sock"
    
    // Track if job was cancelled
    jobWasCancelled := false
    
    stoppableJob := func(args map[string]string, signature string, quit chan struct{}) bool {
        for i := 0; i < 100; i++ { // Long-running job
            select {
            case <-quit:
                jobWasCancelled = true
                return true
            case <-time.After(10 * time.Millisecond):
                // Do work
            }
        }
        return true
    }
    
    if err := server.AddStoppableJob("stop-test", stoppableJob); err != nil {
        t.Fatalf("Failed to register stoppable test job: %v", err)
    }
    
    // Start server
    go server.NewServer(&utils.DefaultLogger{}, socketPath).Serve()
    
    // Give server time to start
    time.Sleep(100 * time.Millisecond)
    
    // Start the job
    cli := client.NewClient(&utils.DefaultLogger{}, socketPath)
    go func() {
        result := cli.Run(&client.Task{
            Name: "stop-test",
        })
        if result != base.INTERRUPT { // Expected since we'll stop it
            t.Errorf("Expected INTERRUPT, got %v", result)
        }
    }()
    
    // Let it run a bit
    time.Sleep(50 * time.Millisecond)
    
    // Stop the job
    stopResult := cli.Run(&client.Task{
        Name: "stop-test",
        Stop: true,
    })
    
    if stopResult != base.SUCCESS {
        t.Errorf("Expected SUCCESS from stop, got %v", stopResult)
    }
    
    if !jobWasCancelled {
        t.Error("Expected job to be cancelled")
    }
}
```

## Testing with Test Servers

### Creating Test-Specific Servers

```go
package main

import (
    "sync"
    "time"

    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

type TestServer struct {
    socketPath string
    server     *server.ser
    started    chan bool
    wg         sync.WaitGroup
}

func NewTestServer(socketPath string) *TestServer {
    return &TestServer{
        socketPath: socketPath,
        started:    make(chan bool),
    }
}

func (ts *TestServer) Start() {
    ts.wg.Add(1)
    go func() {
        defer ts.wg.Done()
        ts.started <- true
        ts.server = server.NewServer(&utils.DefaultLogger{}, ts.socketPath)
        ts.server.Serve()
    }()
    
    // Wait for server to indicate it's started
    <-ts.started
    // Give it a moment to fully initialize
    time.Sleep(50 * time.Millisecond)
}

func (ts *TestServer) Stop() {
    // In tests, we usually just let the process end
    // For more sophisticated cleanup, you might need signal handling
    ts.wg.Wait()
}

func (ts *TestServer) RegisterJob(name string, handler server.JobHandler) error {
    return server.AddJob(name, handler)
}

func (ts *TestServer) RegisterStoppableJob(name string, handler server.StoppableJobHandler) error {
    return server.AddStoppableJob(name, handler)
}

// Example usage in test
func TestWithTestServer(t *testing.T) {
    testServer := NewTestServer("/tmp/test_with_server.sock")
    
    // Register test jobs
    jobExecuted := make(chan bool, 1)
    testServer.RegisterJob("test-job", func(args map[string]string, signature string) bool {
        jobExecuted <- true
        return true
    })
    
    testServer.Start()
    
    // Now test using the server
    // ... test code here
    
    // Verify expectations
    select {
    case <-jobExecuted:
        // Success
    case <-time.After(1 * time.Second):
        t.Error("Test job was not executed in time")
    }
    
    testServer.Stop()
}
```

## Testing Error Conditions

### Testing Error Handling

```go
func TestJobErrorHandling(t *testing.T) {
    tests := []struct {
        name          string
        args          map[string]string
        expectedError bool
    }{
        {
            name:          "missing required parameter",
            args:          map[string]string{},
            expectedError: true, // Assuming the job returns false for missing params
        },
        {
            name: "with required parameter",
            args: map[string]string{
                "required": "value",
            },
            expectedError: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := jobThatValidatesArgs(tt.args, "test-sig")
            
            if tt.expectedError && result {
                t.Error("Expected job to fail but it succeeded")
            }
            if !tt.expectedError && !result {
                t.Error("Expected job to succeed but it failed")
            }
        })
    }
}
```

## Performance Testing

### Benchmarking Jobs

```go
// benchmark_test.go
package main

import (
    "testing"
)

func BenchmarkSimpleJob(b *testing.B) {
    args := map[string]string{
        "test": "value",
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = simpleJob(args, "benchmark-sig")
    }
}

func BenchmarkStoppableJob(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        quit := make(chan struct{})
        done := make(chan bool, 1)
        
        go func() {
            result := fastStoppableJob(map[string]string{}, "benchmark-sig", quit)
            done <- result
        }()
        
        close(quit)
        <-done
    }
}

func BenchmarkConcurrentJobs(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            result := simpleJob(map[string]string{}, "parallel-sig")
            if !result {
                b.Error("Job failed in parallel benchmark")
            }
        }
    })
}
```

### Memory Profiling

```go
// Use Go's built-in memory profiling
func TestMemoryUsage(t *testing.T) {
    // This test can be run with memory profiling:
    // go test -memprofile=mem.prof -run=TestMemoryUsage
    // go tool pprof mem.prof
    
    // Simulate work that might use memory
    for i := 0; i < 1000; i++ {
        result := processLargeData()
        if result == nil {
            t.Errorf("Result was nil at iteration %d", i)
        }
    }
}
```

## Testing Utilities and Helpers

### Test Helpers

```go
package main

import (
    "sync"
    "time"
)

// TestJobTracker helps track job execution for testing
type TestJobTracker struct {
    executedJobs map[string][]map[string]string
    mu          sync.Mutex
    executionCount map[string]int
}

func NewTestJobTracker() *TestJobTracker {
    return &TestJobTracker{
        executedJobs:   make(map[string][]map[string]string),
        executionCount: make(map[string]int),
    }
}

func (tjt *TestJobTracker) TrackExecution(jobName string, args map[string]string) {
    tjt.mu.Lock()
    defer tjt.mu.Unlock()
    
    tjt.executedJobs[jobName] = append(tjt.executedJobs[jobName], args)
    tjt.executionCount[jobName]++
}

func (tjt *TestJobTracker) GetExecutionCount(jobName string) int {
    tjt.mu.Lock()
    defer tjt.mu.Unlock()
    
    return tjt.executionCount[jobName]
}

func (tjt *TestJobTracker) GetExecutions(jobName string) []map[string]string {
    tjt.mu.Lock()
    defer tjt.mu.Unlock()
    
    executions := make([]map[string]string, len(tjt.executedJobs[jobName]))
    copy(executions, tjt.executedJobs[jobName])
    return executions
}

// Example of using the tracker
func testableJob(tracker *TestJobTracker, args map[string]string, signature string) bool {
    tracker.TrackExecution("testable-job", args)
    
    // Do actual work
    return true
}

func TestJobWithTracker(t *testing.T) {
    tracker := NewTestJobTracker()
    
    // Run job multiple times
    for i := 0; i < 5; i++ {
        args := map[string]string{"iteration": string(rune(i + '0'))}
        result := testableJob(tracker, args, "test-sig")
        if !result {
            t.Errorf("Job failed on iteration %d", i)
        }
    }
    
    // Verify tracking worked
    if count := tracker.GetExecutionCount("testable-job"); count != 5 {
        t.Errorf("Expected 5 executions, got %d", count)
    }
    
    executions := tracker.GetExecutions("testable-job")
    if len(executions) != 5 {
        t.Errorf("Expected 5 executions recorded, got %d", len(executions))
    }
}
```

## Continuous Testing

### Test Organization and Running

```bash
# Run all tests
go test ./...

# Run tests in verbose mode
go test -v ./...

# Run with race detection (important for concurrent code)
go test -race ./...

# Run specific test package
go test ./server/... -v

# Run specific test
go test -run TestClientServerIntegration -v

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Environment Setup

```go
// Use build tags for test-specific setup
//go:build test
// +build test

package main

import (
    "os"
    "path/filepath"
    "testing"
)

var testTempDir string

func TestMain(m *testing.M) {
    // Setup
    var err error
    testTempDir, err = os.MkdirTemp("", "saturn-test-*")
    if err != nil {
        panic(err)
    }
    
    // Run tests
    exitCode := m.Run()
    
    // Teardown
    os.RemoveAll(testTempDir)
    
    os.Exit(exitCode)
}

func getTestSocketPath(name string) string {
    return filepath.Join(testTempDir, name+".sock")
}

func TestWithTempSockets(t *testing.T) {
    socketPath := getTestSocketPath("my-test")
    
    // Use socketPath for test server
    // ... test code here
}
```

## Testing Best Practices

### What to Test

1. **Happy Path**: Normal operation scenarios
2. **Error Conditions**: Invalid inputs, resource failures
3. **Edge Cases**: Empty inputs, maximum values, boundary conditions
4. **Concurrent Access**: Multiple clients accessing server simultaneously
5. **Cancellation**: Proper stoppable job behavior
6. **Resource Cleanup**: Proper handling of resources
7. **Parameter Validation**: All parameter validation logic

### What NOT to Test

1. **Implementation Details**: Focus on behavior, not internal implementation
2. **Third-party Dependencies**: Mock external services
3. **Timing-Dependent Code**: Use channels and synchronization instead
4. **Non-Deterministic Behavior**: Ensure tests are repeatable

### Test Naming Conventions

Follow the pattern: `Test[Feature][Scenario][ExpectedResult]`

Examples:
- `TestJobRegistrationSuccess`
- `TestStoppableJobStop`
- `TestJobWithInvalidParamsFails`
- `TestConcurrentJobExecution`

## Troubleshooting Tests

### Common Testing Issues

#### Flaky Tests
```go
// Fix timing-dependent tests with proper synchronization
func TestJobWithWait(t *testing.T) {
    done := make(chan bool, 1)
    
    go func() {
        result := slowJob()
        done <- result
    }()
    
    select {
    case result := <-done:
        if !result {
            t.Error("Job failed")
        }
    case <-time.After(5 * time.Second): // Reasonable timeout
        t.Error("Job did not complete in time")
    }
}
```

#### Resource Conflicts
```go
// Use unique resources per test
func getUniqueSocketPath(t *testing.T) string {
    return fmt.Sprintf("/tmp/test-%s.sock", t.Name())
}

func TestUniqueResources(t *testing.T) {
    socketPath := getUniqueSocketPath(t)
    
    // Now each test has its own socket path
    // ... test code
}
```

## Next Steps

- Review the [Development Setup](./development-setup.md) guide for environment configuration
- Check the [Architecture](./architecture.md) documentation for system understanding
- Look at the [Best Practices](./best-practices.md) guide for quality testing patterns
- Follow the [Contributing](./contributing.md) guidelines if contributing code