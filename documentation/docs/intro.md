---
sidebar_position: 0
title: Saturn CLI Documentation - Go Job Execution Toolkit
description: Complete documentation for Saturn CLI, a lightweight client/server toolkit for job execution in Go. Execute background jobs, manage processes, and enable inter-process communication.
keywords: [saturn cli, go job execution, background jobs, unix domain sockets, ipc, golang toolkit, cli automation]
---

# Welcome to Saturn CLI Documentation - Go Job Execution Toolkit

Saturn CLI is a lightweight client/server toolkit that lets you trigger and monitor shell-style jobs from your Go services or the command line. The client communicates with a long-running daemon over Unix domain sockets (macOS/Linux) or HTTP (Windows), providing a fast, secure channel for orchestrating background work and process management.

## About This Project

This project is a Go CLI client implementation based on the [Saturn](https://github.com/vipshop/Saturn) distributed task scheduling system originally developed by VipShop. Saturn CLI provides a lightweight alternative for executing jobs from Go applications using Unix domain sockets for efficient local communication.

## Key Features

- **Embeddable job runtime** – Register regular or stoppable jobs using a composable registry that can be scoped per service instance
- **Socket-first transport** – Uses Unix domain sockets by default for low latency and predictable permissions, with automatic TCP fallback on Windows
- **Graceful cancellation** – Stoppable jobs receive a dedicated quit channel and the client exposes `--stop` semantics and `CTRL+C` interception
- **Structured CLI experience** – Repeatable `--param key=value` flags, helpful diagnostics, and consistent exit codes make automation scripts straightforward
- **Production-ready ergonomics** – Context-aware HTTP clients, signal cleanup, and comprehensive integration tests keep long-lived processes healthy

## How to Navigate This Documentation

This documentation is organized for different types of users. Choose the path that best matches your needs:

### For New Users
- Start with our [Getting Started](./getting-started.md) guide to understand basic concepts and setup
- Try the [Quick Start](./getting-started/quick-start.md) for immediate hands-on experience
- Explore [Basic Examples](./getting-started/examples.md) to see common use cases

### For Developers
- Review [Architecture](./architecture.md) to understand the system design
- Read [API References](./client-api.md) for detailed technical specifications
- Follow the [Embedding Guide](./embedding.md) to integrate Saturn into your services

### For Contributors
- Check out [Contributing](./contributing.md) for development setup
- Review [Development Setup](./development-setup.md) for building and testing
- Look at [Testing](./testing.md) for guidelines on writing tests

## Highlights

- **Embeddable job runtime** – Register regular or stoppable jobs using a composable registry that can be scoped per service instance.
- **Socket-first transport** – Uses Unix domain sockets by default for low latency and predictable permissions, with automatic TCP fallback on Windows.
- **Graceful cancellation** – Stoppable jobs receive a dedicated quit channel and the client exposes `--stop` semantics and `CTRL+C` interception.
- **Structured CLI experience** – Repeatable `--param key=value` flags, helpful diagnostics, and consistent exit codes make automation scripts straightforward.
- **Production-ready ergonomics** – Context-aware HTTP clients, signal cleanup, and comprehensive integration tests keep long-lived processes healthy.

## Architecture Overview

![Saturn CLI design overview](https://github.com/Kingson4Wu/saturncli/blob/main/resource/img/design-overview-saturn-cli-go.png)

## Quick Start

To get started with Saturn CLI, check out our [Getting Started](./getting-started.md) guide.