---
sidebar_position: 9
title: Troubleshooting Saturn CLI - Common Issues & Solutions
description: Solutions to common issues when using Saturn CLI. Troubleshoot connection problems, job execution errors, and configuration issues.
keywords: [saturn cli troubleshooting, go cli errors, unix domain sockets issues, saturn cli problems, job execution errors, cli debugging]
---

# Troubleshooting

This guide covers common issues you may encounter when using Saturn CLI and how to resolve them for reliable job execution.

## Connection Issues

### "Connection Refused" Error

**Symptoms**: Client fails to connect to server with "connection refused" error.

**Causes and Solutions**:
1. **Server Not Running**: Ensure the Saturn server is running before starting client requests
   ```bash
   # First, start the server
   ./saturn_svr
   
   # Then run client commands in another terminal
   ./saturn_cli --name hello
   ```

2. **Socket Path Mismatch**: Verify that client and server use the same socket path
   ```go
   // Server side - make sure to use the same path as client
   server.NewServer(&utils.DefaultLogger{}, "/tmp/myservice.sock")
   
   // Client side - must match server path
   cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/myservice.sock")
   ```

3. **File Permissions**: Check that your user has permission to access the socket file
   ```bash
   # Check socket directory permissions
   ls -la /tmp/
   
   # Ensure directory is writable by your user
   chmod 755 /tmp
   ```

### Socket File Permission Errors

**Symptoms**: Errors creating or accessing socket files.

**Solutions**:
1. **Check Directory Permissions**: Ensure the socket directory is writable
   ```bash
   # Check permissions
   ls -la /tmp/
   
   # Create directory if needed with proper permissions
   mkdir -p /tmp && chmod 755 /tmp
   ```

2. **Use Writable Path**: Use a path where your user has write permissions
   ```go
   // Instead of using /root/ (if you're not root)
   server.NewServer(&utils.DefaultLogger{}, "/tmp/saturn.sock")
   
   // Or use user-specific directory
   server.NewServer(&utils.DefaultLogger{}, "/home/user/.saturn.sock")
   ```

3. **Clean Up Old Sockets**: Remove any stale socket files from previous runs
   ```bash
   rm -f /tmp/saturn.sock
   ```

## Job Registration Issues

### "Job Not Found" Error

**Symptoms**: Client returns error indicating job name doesn't exist.

**Causes and Solutions**:
1. **Registration Timing**: Ensure the job is registered before client attempts to run it
   ```go
   // Register jobs BEFORE starting the server
   if err := server.AddJob("hello", func(args map[string]string, signature string) bool {
       fmt.Println("Hello job executed")
       return true
   }); err != nil {
       log.Fatal(err)
   }
   
   // Server can now accept requests for "hello" job
   server.NewServer(&utils.DefaultLogger{}, "/tmp/saturn.sock").Serve()
   ```

2. **Case Sensitivity**: Job names are case-sensitive
   ```bash
   # This job name...
   ./saturn_cli --name Hello  # Wrong case
   
   # Should match exactly what's registered
   ./saturn_cli --name hello  # Correct
   ```

## Platform-Specific Issues

### Windows TCP Connection Issues

**Symptoms**: Server runs but client can't connect on Windows.

**Solutions**:
1. **Port Conflicts**: The default port (8096) might be in use
   ```go
   // On Windows, server uses TCP instead of Unix sockets
   // Port is hardcoded in server_windows.go
   // Make sure port 8096 is available
   netstat -an | grep 8096  // Check if port is in use
   ```

2. **Firewall Issues**: Windows firewall might block the connection
   ```go
   // Saturn CLI uses localhost only, so should be allowed by default
   // But if issues persist, check Windows Firewall settings
   ```

## Stoppable Job Issues

### Jobs Don't Stop When Requested

**Symptoms**: `--stop` flag or `Ctrl+C` doesn't terminate a stoppable job.

**Causes and Solutions**:
1. **Missing Quit Channel Check**: Job handler must regularly check the quit channel
   ```go
   // Incorrect - doesn't check quit channel
   func badHandler(args map[string]string, signature string, quit chan struct{}) bool {
       for i := 0; i < 1000000; i++ {
           // Does work without checking quit
           fmt.Printf("Working... %d\n", i)
       }
       return true
   }
   
   // Correct - checks quit channel regularly
   func goodHandler(args map[string]string, signature string, quit chan struct{}) bool {
       for i := 0; i < 1000000; i++ {
           select {
           case <-quit:
               fmt.Printf("Job %s stopped at iteration %d\n", signature, i)
               return true
           default:
               // Do a small amount of work
               fmt.Printf("Working... %d\n", i)
               
               // Small sleep to yield to quit check
               time.Sleep(10 * time.Millisecond)
           }
       }
       return true
   }
   ```

2. **Blocking Operations**: Long-running operations should be interruptible
   ```go
   // Use context for I/O operations to make them cancellable
   func cancellableIOHandler(args map[string]string, signature string, quit chan struct{}) bool {
       ctx, cancel := context.WithCancel(context.Background())
       
       // Start a goroutine to cancel context when quit is received
       go func() {
           <-quit
           cancel()
       }()
       
       // Use ctx in your I/O operations
       // ... perform cancellable I/O
       
       return true
   }
   ```

## Parameter Issues

### Parameters Not Passed Correctly

**Symptoms**: Job handlers receive empty or unexpected parameter values.

**Causes and Solutions**:
1. **Parameter Format**: Ensure correct format for different parameter types
   ```bash
   # Correct structured parameters
   ./saturn_cli --name job --param key1=value1 --param key2=value2
   
   # Correct legacy format
   ./saturn_cli --name job --args 'key1=value1&key2=value2'
   
   # Both can be used together (params take precedence)
   ./saturn_cli --name job --args 'common=old' --param common=new --param other=also_new
   # Results in: {"common": "new", "other": "also_new"}
   ```

2. **Special Characters**: URL-encode special characters in parameters
   ```bash
   # If your parameter value contains special characters
   ./saturn_cli --name job --param message="Hello%20World"  # For "Hello World"
   ```

## Performance Issues

### High Memory Usage

**Symptoms**: Memory usage grows over time or is higher than expected.

**Solutions**:
1. **Resource Cleanup**: Ensure job handlers properly clean up resources
   ```go
   func wellBehavedJob(args map[string]string, signature string) bool {
       // Open resources
       file, err := os.Open("path/to/file")
       if err != nil {
           return false
       }
       defer file.Close()  // Proper cleanup
       
       // ... use file
       
       return true
   }
   ```

2. **Goroutine Leaks**: Don't leave goroutines running after job completion
   ```go
   func leakyJob(args map[string]string, signature string) bool {
       // This goroutine might run after job completes
       go func() {
           // Work that might outlive the function
       }()
       return true  // Function returns but goroutine continues
   }
   
   func nonLeakyJob(args map[string]string, signature string) bool {
       done := make(chan bool, 1)
       go func() {
           // Work that should complete before function returns
           // ... do work
           done <- true
       }()
       
       <-done  // Wait for goroutine to complete
       return true
   }
   ```

## Debugging Tips

### Enable Detailed Logging

Use a logger that provides more detailed output for debugging:

```go
type DebugLogger struct {
    *log.Logger
}

func (l *DebugLogger) Debug(msg string) {
    l.Logger.Printf("[DEBUG] %s", msg)
}

func (l *DebugLogger) Debugf(format string, args ...interface{}) {
    l.Logger.Printf("[DEBUG] "+format, args...)
}

// Use the debug logger
debugLogger := &DebugLogger{log.New(os.Stdout, "", log.LstdFlags)}
server.NewServer(debugLogger, "/tmp/debug.sock").Serve()
```

### Test Job Registration Directly

Verify job registration works as expected:

```go
// Test that your job can be registered and executed
func TestJobRegistration(t *testing.T) {
    // Register a test job
    executed := false
    var executedArgs map[string]string
    
    testJob := func(args map[string]string, signature string) bool {
        executed = true
        executedArgs = args
        return true
    }
    
    if err := server.AddJob("test-job", testJob); err != nil {
        t.Fatalf("Failed to register test job: %v", err)
    }
    
    // Create and run client
    cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/test.sock")
    result := cli.Run(&client.Task{
        Name:   "test-job",
        Params: map[string]string{"key": "value"},
    })
    
    if result != base.SUCCESS {
        t.Errorf("Job execution failed, got: %v", result)
    }
    
    if !executed {
        t.Error("Job was not executed")
    }
}
```

### Common Exit Codes and Their Meanings

- **0**: Success or Interrupt (job completed successfully or was interrupted)
- **1**: Failure (job returned false, or client/server communication error)

## Getting Help

If you encounter issues not covered in this guide:

1. **Check GitHub Issues**: Look for similar problems in the [issue tracker](https://github.com/Kingson4Wu/saturncli/issues)
2. **Enable Debug Logging**: Add more logging to understand the execution flow
3. **Verify Installation**: Ensure you're using the correct versions of client and server
4. **Test with Examples**: Verify basic functionality with the example code
5. **Create Minimal Reproduction**: Create a minimal example that reproduces the issue

## Next Steps

- Review the [Architecture](./architecture.md) documentation for deeper understanding
- Check the [Examples](./examples.md) for correct usage patterns
- Look at the [API Reference](./client-api.md) for detailed function documentation