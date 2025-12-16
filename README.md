[![CI/CD Pipeline](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml/badge.svg)](https://github.com/Kingson4Wu/saturncli/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kingson4Wu/saturncli)](https://goreportcard.com/report/github.com/Kingson4Wu/saturncli)
![GitHub top language](https://img.shields.io/github/languages/top/Kingson4Wu/saturncli)
[![GitHub stars](https://img.shields.io/github/stars/Kingson4Wu/saturncli)](https://github.com/Kingson4Wu/saturncli/stargazers)
[![codecov](https://codecov.io/gh/Kingson4Wu/saturncli/branch/main/graph/badge.svg)](https://codecov.io/gh/Kingson4Wu/saturncli)
[![Go Reference](https://pkg.go.dev/badge/github.com/Kingson4Wu/saturncli.svg)](https://pkg.go.dev/github.com/Kingson4Wu/saturncli)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go#database)
[![LICENSE](https://img.shields.io/github/license/Kingson4Wu/saturncli.svg?style=flat-square)](https://github.com/Kingson4Wu/saturncli/blob/main/LICENSE)

English | [简体中文](https://github.com/Kingson4Wu/saturncli/blob/main/README-CN.md) | [deepwiki](https://deepwiki.com/Kingson4Wu/saturncli)

---

# Saturn CLI (Go)

Saturn CLI is a lightweight client/server toolkit that lets you trigger and monitor shell-style jobs from your Go services or the command line. The client communicates with a long-running daemon over Unix domain sockets (macOS/Linux) or HTTP (Windows), providing a fast, secure channel for orchestrating background work such as scheduled tasks in [Saturn](https://github.com/vipshop/Saturn).

This project is a Go CLI client implementation based on the [Saturn](https://github.com/vipshop/Saturn) distributed task scheduling system originally developed by VipShop. The Saturn CLI toolkit is designed to be embedded directly into existing services. You can register jobs programmatically, expose them via the bundled server, and control execution through the CLI or custom integrations built on top of the client package.

## Documentation

Complete documentation for Saturn CLI is available in the [documentation](./documentation/) directory and can be viewed online at our [documentation site](https://kingson4wu.github.io/saturncli/). The documentation includes:
- Getting started guides
- API references
- Architecture documentation
- Usage examples
- Best practices

## Highlights

- **Embeddable job runtime** – Register regular or stoppable jobs using a composable registry that can be scoped per service instance.
- **Socket-first transport** – Uses Unix domain sockets by default for low latency and predictable permissions, with automatic TCP fallback on Windows.
- **Graceful cancellation** – Stoppable jobs receive a dedicated quit channel and the client exposes `--stop` semantics and `CTRL+C` interception.
- **Structured CLI experience** – Repeatable `--param key=value` flags, helpful diagnostics, and consistent exit codes make automation scripts straightforward.
- **Production-ready ergonomics** – Context-aware HTTP clients, signal cleanup, and comprehensive integration tests keep long-lived processes healthy.

## Architecture Overview

![Saturn CLI design overview](https://github.com/Kingson4Wu/saturncli/blob/main/resource/img/design-overview-saturn-cli-go.png)

## Getting Started

### Prerequisites

- Go 1.19 or newer
- Unix-like system for socket transport (Windows is supported via TCP loopback)

### Build

```bash
make
```

This produces two binaries:

- `saturn_svr` – reference server demonstrating job registration
- `saturn_cli` – command-line client

### Run the demo

```bash
./saturn_svr
./saturn_cli --name hello --param id=33 --param ver=22
./saturn_cli --name hello_stoppable
./saturn_cli --name hello_stoppable --stop
```

Use `CTRL+C` while a stoppable job is running to trigger an interrupt with automatic stop propagation.

## Embedding in Your Service

```go
package main

import (
    "log"

    "github.com/Kingson4Wu/saturncli/server"
    "github.com/Kingson4Wu/saturncli/utils"
)

func main() {
    registry := server.NewRegistry()

    if err := registry.AddJob("hello", func(args map[string]string, signature string) bool {
        log.Printf("hello %v (run=%s)", args, signature)
        return true
    }); err != nil {
        log.Fatal(err)
    }

    if err := registry.AddStoppableJob("slow-task", func(args map[string]string, signature string, quit chan struct{}) bool {
        for {
            select {
            case <-quit:
                return true
            default:
                // do work...
            }
        }
    }); err != nil {
        log.Fatal(err)
    }

    server.NewServer(&utils.DefaultLogger{}, "/tmp/saturn.sock", server.WithRegistry(registry)).Serve()
}
```

Clients can be embedded as well:

```go
import (
    "github.com/Kingson4Wu/saturncli/client"
    "github.com/Kingson4Wu/saturncli/utils"
)

cli := client.NewClient(&utils.DefaultLogger{}, "/tmp/saturn.sock")
result := cli.Run(&client.Task{
    Name:   "hello",
    Params: map[string]string{"id": "42"},
})
```

`result` will contain one of `success`, `failure`, or `interrupt`.

## CLI Reference

```bash
saturn_cli [FLAGS]

Flags:
  --name string         Job name to execute (required)
  --args string         Legacy query string, merged with --param values
  --param key=value     Repeatable structured argument
  --stop                Send a stop signal instead of starting a job
  --signature string    Target a specific run when stopping
  --help                Show detailed usage
```

The CLI returns exit code `0` on success or interrupt, and `1` on failure.

## Server API Reference

- `server.NewRegistry()` – create an isolated registry
- `registry.AddJob(name, handler)` – register a synchronous job
- `registry.AddStoppableJob(name, handler)` – register a job that accepts a quit channel
- `server.NewServer(logger, sockPath, opts...)` – construct a server bound to a socket path
- `server.WithRegistry(registry)` – inject a custom registry (defaults to a package-level shared registry)

Use `client.Task.Params` for structured argument passing; `Task.Args` is preserved for existing integrations that already supply URL query strings.

## Testing

```bash
go test ./...
```

The suite includes end-to-end tests that spin up in-memory registries and verify stop semantics.

## Documentation

- [Project wiki](https://github.com/Kingson4Wu/saturncli/wiki)
- [Examples](https://github.com/Kingson4Wu/saturncli/tree/main/examples)

## Contributing

Contributions are welcome! Please read the [CONTRIBUTING](https://github.com/Kingson4Wu/saturncli/blob/main/CONTRIBUTING.md) guide before opening issues or pull requests.

## License

Saturn CLI is licensed under the terms of the [Apache 2.0 License](https://github.com/Kingson4Wu/saturncli/blob/main/LICENSE).
