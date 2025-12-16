---
sidebar_position: 7
title: CLI Reference - Saturn Command Line Interface Guide
description: Complete command-line reference for Saturn CLI. Learn all available options, commands, and usage patterns for executing jobs from the command line.
keywords: [saturn cli reference, command line interface, go cli commands, saturn cli options, cli job execution, saturn cli usage]
---

# CLI Reference

This document provides a complete reference for the Saturn command-line interface, including all available options, commands, and usage patterns for executing background jobs.

## Command Syntax

The basic syntax for the Saturn CLI is:

```bash
saturn_cli [FLAGS]
```

## Available Flags

### --name (Required)
Specifies the name of the job to execute on the Saturn server.

- **Type**: String
- **Required**: Yes
- **Shorthand**: None

**Example:**
```bash
saturn_cli --name hello
```

### --param (Repeatable)
Provides structured arguments to the job in key-value format. Multiple `--param` flags can be specified.

- **Type**: Key-Value pair (format: key=value)
- **Required**: No
- **Shorthand**: None

**Example:**
```bash
saturn_cli --name hello --param id=33 --param version=2.0
```

### --args
Provides legacy query string arguments. These are merged with `--param` values.

- **Type**: String (URL query format: key1=value1&key2=value2)
- **Required**: No
- **Shorthand**: None

**Example:**
```bash
saturn_cli --name hello --args 'id=33&version=2.0'
```

### --stop
Sends a stop signal to an already running job instead of starting a new one.

- **Type**: Boolean
- **Required**: No (only when stopping jobs)
- **Shorthand**: None

**Example:**
```bash
saturn_cli --name long_running_job --stop
```

### --signature
Targets a specific job run when stopping. Used in conjunction with `--stop`.

- **Type**: String
- **Required**: No
- **Shorthand**: None

**Example:**
```bash
saturn_cli --name long_running_job --stop --signature "job-12345"
```

### --help
Displays detailed usage information for the Saturn CLI.

- **Type**: Boolean
- **Required**: No
- **Shorthand**: None

**Example:**
```bash
saturn_cli --help
```

## Parameter Handling

Saturn CLI supports two methods of providing parameters to jobs:

### Structured Parameters (--param)
The preferred method for passing parameters:

```bash
saturn_cli --name my_job --param key1=value1 --param key2=value2
```

This creates a parameter map: `{"key1": "value1", "key2": "value2"}`

### Legacy Query String (--args)
For backward compatibility with existing systems:

```bash
saturn_cli --name my_job --args 'key1=value1&key2=value2'
```

### Combined Approach
Both parameter formats can be used together, with `--param` values taking precedence:

```bash
saturn_cli --name my_job --args 'common=value' --param key1=value1 --param common=override
```

Results in: `{"common": "override", "key1": "value1"}`

## Exit Codes

The Saturn CLI returns the following exit codes:

- **0**: Success or Interrupt
  - Job completed successfully
  - Job was interrupted (e.g., by `--stop` or `Ctrl+C`)
- **1**: Failure
  - Job failed during execution
  - Invalid command-line arguments
  - Communication error with the server

## Examples

### Basic Job Execution
Execute a simple job with parameters:

```bash
saturn_cli --name hello --param id=42 --param message=Greeting
```

### Job with Multiple Parameters
Run a job with multiple parameters:

```bash
saturn_cli --name data_processor --param input_file=/path/to/input --param output_dir=/path/to/output --param format=json
```

### Stopping a Running Job
Stop a currently running job:

```bash
saturn_cli --name long_running_task --stop
```

### Targeting Specific Job Instance
Stop a specific instance of a job using its signature:

```bash
saturn_cli --name batch_processor --stop --signature "batch-20231201-001"
```

### Combining Args and Params
Use both legacy args and structured params:

```bash
saturn_cli --name my_job --args 'legacy_param=value' --param new_param=new_value
```

### Interrupting Jobs with Ctrl+C
Start a stoppable job and interrupt it using `Ctrl+C`:

```bash
saturn_cli --name hello_stoppable
# Then press Ctrl+C to send interrupt signal
```

## Server Connection

The Saturn CLI connects to the Saturn server using:

- **Unix-like systems**: Unix domain socket (default path set during server configuration)
- **Windows**: TCP connection to localhost

The socket path or TCP address is configured when starting the Saturn server and must match the client configuration.

## Configuration

### Socket Path
The Saturn client requires knowledge of the server's socket path, which is typically configured when starting the Saturn server. The client uses this path to establish communication.

### Logging
The Saturn CLI writes diagnostic information to stderr to avoid interfering with job output:

- Success messages: "Execution Success"
- Interrupt messages: "Execution Interrupted" 
- Failure messages: "Execution Failure"

## Integration with Scripts

The predictable exit codes and structured parameter support make Saturn CLI suitable for use in automation scripts:

```bash
#!/bin/bash

# Run a job and check the result
if saturn_cli --name backup --param source=/data --param destination=/backup; then
    echo "Backup completed successfully"
else
    echo "Backup failed"
    exit 1
fi

# Run multiple jobs in sequence
for dataset in dataset1 dataset2 dataset3; do
    if ! saturn_cli --name process --param dataset=$dataset; then
        echo "Processing failed for $dataset"
        exit 1
    fi
done
```

## Troubleshooting

### Common Issues

#### Missing Required Name Parameter
**Error**: Command fails without providing job name
**Solution**: Always specify the `--name` parameter

```bash
saturn_cli --name my_job  # Correct
saturn_cli                # Incorrect - will fail
```

#### Invalid Parameter Format
**Error**: Parameter flag used without key=value format
**Solution**: Use proper key=value format

```bash
saturn_cli --name my_job --param key=value  # Correct
saturn_cli --name my_job --param value      # Incorrect
```

#### Connection Issues
**Error**: Cannot connect to Saturn server
**Solution**: Verify server is running and socket path is correct

#### Job Not Found
**Error**: Requested job name doesn't exist on the server
**Solution**: Ensure the job has been registered with the Saturn server

### Debugging Tips

1. **Check Server Status**: Verify the Saturn server is running and accessible
2. **Verify Socket Path**: Ensure client and server use the same socket path
3. **Validate Job Names**: Confirm job names are registered on the server
4. **Review Job Parameters**: Check that all required parameters are provided

## Advanced Usage

### Script Integration
Combine Saturn CLI with shell scripting for complex workflows:

```bash
#!/bin/bash

# Execute a sequence of jobs with error handling
execute_job_sequence() {
    local job_data=("$@")
    
    for job in "${job_data[@]}"; do
        echo "Executing job: $job"
        if ! saturn_cli --name "$job" --param timestamp="$(date -Iseconds)"; then
            echo "Job $job failed, stopping sequence"
            return 1
        fi
    done
    
    return 0
}

# Usage
job_sequence=("validate" "process" "archive")
if execute_job_sequence "${job_sequence[@]}"; then
    echo "All jobs completed successfully"
else
    echo "Job sequence failed"
fi
```

## See Also

- [Getting Started Guide](./getting-started.md) - For basic installation and usage
- [Server API Reference](./server-api.md) - For server-side configuration
- [Client API Reference](./client-api.md) - For programmatic client usage
- [Architecture](./architecture.md) - For system design information