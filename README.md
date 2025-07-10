# a2a-go: Agent-to-Agent Protocol SDK for Go

[![GoDoc](https://pkg.go.dev/badge/github.com/yeeaiclub/a2a-go)](https://pkg.go.dev/github.com/yeeaiclub/a2a-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/yeeaiclub/a2a-go/main?logo=codecov&logoColor=white)](https://codecov.io/gh/yeeaiclub/a2a-go/branch/main)
[![Status](https://img.shields.io/badge/Status-Under%20Development-orange.svg)](https://github.com/yeeaiclub/a2a-go)

**a2a-go** is a comprehensive Go SDK implementation of the **Agent-to-Agent (A2A) protocol**, providing a robust foundation for building AI agent communication systems. This SDK enables seamless integration between AI agents with middleware support, authentication, and full protocol method implementation.

> **⚠️ This project is currently under active development. APIs may change without notice.**

## Overview

The **a2a-go** SDK implements the [A2A protocol specification](https://github.com/a2aproject/A2A) in Go, offering:

- **Complete A2A Protocol Support**: Full implementation of all A2A protocol methods
- **Go-Native Design**: Built specifically for Go applications with idiomatic patterns
- **Middleware Architecture**: Extensible middleware system for custom request/response processing
- **Multi-Authentication Support**: API Key, Bearer, OAuth2, and OpenID Connect
- **Real-time Streaming**: Support for streaming task results and events
- **Type Safety**: Comprehensive Go type definitions with full type safety

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
			TaskID: "1",        // The task ID for this message
			Role:   types.User, // The sender's role (User or Agent)
			Parts: []types.Part{
				&types.TextPart{
					Kind: "text",         // The kind of part ("text", "data", etc.)
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


## Use Cases

**a2a-go** is ideal for:

- **AI Agent Communication**: Build systems where multiple AI agents need to communicate and collaborate
- **Task Distribution**: Distribute and manage tasks across multiple agent instances
- **Real-time AI Systems**: Create real-time AI applications with streaming capabilities
- **Microservices Integration**: Integrate AI agents into existing microservices architectures
- **Multi-Agent Workflows**: Orchestrate complex workflows involving multiple specialized agents

## Security Schemes

Support for multiple authentication methods:

- **API Key** - Header or query parameter authentication
- **Bearer Token** - HTTP Bearer authentication
- **OAuth2** - OAuth 2.0 token authentication
- **OpenID Connect** - OpenID Connect token authentication

## Demo

For a complete working example, check out our demo repository:

```shell
git clone https://github.com/yeeaiclub/demo.git
```

## Documentation

- [GoDoc](https://pkg.go.dev/github.com/yeeaiclub/a2a-go) - API documentation
- [A2A Specification](https://github.com/a2aproject/A2A) - Protocol specification

## Related Projects

- [a2a-spec](https://github.com/a2aproject/A2A) - A2A protocol specification documents
- [a2a-python](https://github.com/a2aproject/a2a-python) - Python implementation of A2A protocol

## Contributing

We welcome contributions! Please read our contributing guidelines before submitting pull requests.

> **Note:** Since this project is under active development, we recommend checking the latest issues and discussions before contributing to understand the current development priorities.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---

**Keywords**: a2a-go, a2a, agent-to-agent, Go SDK, AI agents, protocol implementation, middleware, authentication, streaming, task management