# a2a-go

[![GoDoc](https://pkg.go.dev/badge/github.com/yeeaiclub/a2a-go)](https://pkg.go.dev/github.com/yeeaiclub/a2a-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/yeeaiclub/a2a-go/main?logo=codecov&logoColor=white)](https://codecov.io/gh/yeeaiclub/a2a-go/branch/main)
[![Status](https://img.shields.io/badge/Status-Under%20Development-orange.svg)](https://github.com/yeeaiclub/a2a-go)

A comprehensive Go SDK for the Agent-to-Agent (A2A) protocol with middleware support and full protocol method implementation.

> **⚠️ This project is currently under active development. APIs may change without notice.**

## Installation

```shell
go get github.com/yeeaiclub/a2a-go
```

## Features

### Core A2A Protocol Methods
- [x] **Authentication** - Secure agent authentication
- [x] **Send Task** - Submit tasks to other agents
- [x] **Get Task** - Retrieve task information and status
- [x] **Cancel Task** - Cancel running or pending tasks
- [x] **Stream Task** - Real-time task result streaming
- [x] **Set Push Notification** - Configure push notifications for tasks
- [x] **Get Push Notification** - Retrieve push notification configurations

### SDK Features
- [x] **Middleware Support** - Extensible middleware architecture for request/response processing
- [x] **Security Schemes** - Support for multiple authentication methods (API Key, Bearer, OAuth2, OpenID Connect)
- [x] **Context Management** - Flexible context handling with security configuration

## Quick Start

```go
package main

import (
    "log"
    "net/http"
    "github.com/yeeaiclub/a2a-go/sdk/client"
    "github.com/yeeaiclub/a2a-go/sdk/client/middleware"
    "github.com/yeeaiclub/a2a-go/sdk/types"
)

func main() {
    // Step 1: Create a new A2A client with the API endpoint
    a2aClient := client.NewClient(http.DefaultClient, "http://localhost:8080/api")

    // Step 2: Set up authentication middleware (optional, but recommended for secured endpoints)
    credential := middleware.NewInMemoryContextCredentials()
    credential.SetCredentials("session1", "apiKey", "your-api-key")
    a2aClient.Use(middleware.Intercept(credential))

    // Step 3: Build the message payload
    params := types.MessageSendParam{
        Message: &types.Message{
            TaskID: "1", // The task ID for this message
            Role:   types.User, // The sender's role (User or Agent)
            Parts: []types.Part{
                &types.TextPart{
                    Kind: "text", // The kind of part ("text", "data", etc.)
                    Text: "hello, world", // The actual message content
                },
            },
        },
    }

    // Step 4: Send the message to the server
    resp, err := a2aClient.SendMessage(params)
    if err != nil {
        log.Fatalf("SendMessage failed: %v", err)
    }

    // Step 5: Parse the response to get the task information
    task, err := types.MapTo[types.Task](resp.Result)
    if err != nil {
        log.Fatalf("Failed to parse response: %v", err)
    }
    log.Printf("Task ID: %s, Status: %s", task.Id, task.Status.State)
}
```