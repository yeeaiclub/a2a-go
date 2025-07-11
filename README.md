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

## agent Card

Agent Cards define the capabilities and metadata of your AI agent.
```go
card := types.AgentCard{
	Name:        "Weather Agent",
	Description: "Helps with weather",
	URL:         "http://localhost:10001",
	Version:     "1.0.0",
	Capabilities: &types.AgentCapabilities{
		Streaming:              true,
		StateTransitionHistory: false,
    },
	DefaultInputModes:  []string{"text"},
	DefaultOutputModes: []string{"text"},
	Skills: []types.AgentSkill{
		{
			Id:          "weather_search",
			Name:        "Search weather",
			Description: "Helps with weather in city, or states",
			Tags:        []string{"weather"},
			Examples:    []string{"weather in LA, CA"},
        },
    },
}
```

## client

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

## executor

```go
import (
	"context"
	"fmt"

	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/execution"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/updater"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// WeatherAgentExecutorProducer implements the agent executor for weather agent
type WeatherAgentExecutorProducer struct {
	store        tasks.TaskStore
	weatherAgent WeatherAgent
}

// WeatherAgent interface for weather service
type WeatherAgent interface {
	Chat(message string) (string, error)
}

// NewExecutor creates a new weather agent executor
func NewExecutor(store tasks.TaskStore) *WeatherAgentExecutorProducer {
	return &WeatherAgentExecutorProducer{
		store:        store,
		weatherAgent: nil, // your implement
	}
}

// Execute implements the agent execution logic
func (m *WeatherAgentExecutorProducer) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)

	// mark the task as submitted and start working on it
	if requestContext.Task == nil {
		u.Submit()
	}
	u.StartWork()

	// extract the text from the message
	userMessage := m.extractTextFromMessage(requestContext.Params.Message)

	// call the weather agent with the user's message
	response, err := m.weatherAgent.Chat(userMessage)
	if err != nil {
		u.Failed()
		return err
	}

	// create the response part
	responsePart := &types.TextPart{
		Kind: "text",
		Text: response,
	}
	parts := []types.Part{responsePart}

	// add the response as an artifact and complete the task
	u.AddArtifact(parts)
	u.Complete()

	return nil
}

// Cancel implements the task cancellation logic
func (m *WeatherAgentExecutorProducer) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	task := requestContext.Task

	if task == nil {
		return fmt.Errorf("task not found")
	}

	if task.Status.State == types.CANCELED {
		// task already cancelled
		return fmt.Errorf("task already cancelled")
	}

	if task.Status.State == types.COMPLETED {
		// task already completed
		return fmt.Errorf("task already completed")
	}

	// cancel the task
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	u.UpdateStatus(types.CANCELED, updater.WithFinal(true))

	return nil
}

// extractTextFromMessage extracts text content from message parts
func (m *WeatherAgentExecutorProducer) extractTextFromMessage(message *types.Message) string {
	if message == nil || message.Parts == nil {
		return ""
	}

	str := ""
	for _, part := range message.Parts {
		if textPart, ok := part.(*types.TextPart); ok {
			str += textPart.Text
		}
	}
	return str
}
```

## server

This section provides a simple example of how to start and run an a2a-go server locally. 

Below is a basic example of starting the server:

```go
package main

import (
    "log"
    
	"github.com/yeeaiclub/a2a-go/sdk/server/handler"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func main() {
	// Create in-memory task store for storing and managing task states
	store := tasks.NewInMemoryTaskStore()
	// Save an example task to the store (this is just for demonstration, may not be needed in actual use)
	store.Save(context.Background(), &types.Task{Id: "1"})
	
    // TODO: Create your custom executor
    // The executor is responsible for handling specific business logic, such as calling external APIs, processing data, etc.
    executor = NewExecutor(store)
	
	// Create queue manager for handling asynchronous task queues
	queue := NewQueueManager()
	
	// Create default handler, configuring task store, executor, and queue manager
	defaultHandler := handler.NewDefaultHandler(mem, executor, handler.WithQueueManger(queue))

	// Create server instance
	// Parameter descriptions:
	// - "/card": Agent card access path
	// - "/api": API interface access path  
	// - card: Agent card configuration
	// - defaultHandler: Request handler
	server := handler.NewServer("/card", "/api", card, defaultHandler)
	
	// Start the server
	log.Println("Starting server on port 8080...")
	log.Println("Agent card available at: http://localhost:8080/card")
	log.Println("API available at: http://localhost:8080/api")
	if err := server.Start(8080); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
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