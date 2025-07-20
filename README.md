# a2a-go: Agent-to-Agent Protocol SDK for Go

[![GoDoc](https://pkg.go.dev/badge/github.com/yeeaiclub/a2a-go)](https://pkg.go.dev/github.com/yeeaiclub/a2a-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/yeeaiclub/a2a-go/main?logo=codecov&logoColor=white)](https://codecov.io/gh/yeeaiclub/a2a-go/branch/main)
[![Status](https://img.shields.io/badge/Status-Under%20Development-orange.svg)](https://github.com/yeeaiclub/a2a-go)

**a2a-go** is a comprehensive Go SDK implementation of the **Agent-to-Agent (A2A) protocol**, providing a robust foundation for building AI agent communication systems. This SDK enables seamless integration between AI agents with middleware support, authentication, and full protocol method implementation.

> **⚠️ This project is currently under active development. APIs may change without notice.**

## Overview

The **a2a-go** SDK implements the [A2A protocol specification](https://github.com/a2aproject/A2A) in Go, offering:

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


## Installation

```shell
go get github.com/yeeaiclub/a2a-go
```

## Quick Start

## client

```go
// Initialize an HTTP client with a custom timeout (you can adjust config.Timeout as needed)
httpClient := &http.Client{ Timeout: config.Timeout }

// Create a new a2a-go client instance using the HTTP client and the API base URL
// Replace "http://localhost:8080/api" with your actual server endpoint
a2aClient := client.NewClient(httpClient, "http://localhost:8080/api")
```

> The above code demonstrates how to set up the a2a-go client. You need to provide a custom `http.Client` (for timeout, proxy, etc.) and the API endpoint of your a2a server.

```go
// Example: Sending a message using the a2a-go client
resp, err := client.SendMessage(types.MessageSendParam{
    Message: &types.Message{
        TaskID: taskID,           // The ID of the task this message belongs to
        Role:   types.User,       // The sender's role (e.g., User, Agent)
        Parts: []types.Part{
            // Message content parts; here we use a text message as an example
            &types.TextPart{Kind: "text", Text: message},
        },
    },
})
```

> The above code shows how to send a message:
> - `TaskID`: Specifies which task the message is associated with. **This field can be empty if you are starting a new task; the server will generate a new TaskID automatically if not provided.**
> - `Role`: Indicates the sender's role (such as User or Agent).
> - `Parts`: Supports multiple content types (text, image, etc.); here, a text message is used.
> 
> The `resp` variable contains the server's response, and `err` is used for error handling.

If you want to receive messages from the server in a streaming fashion, you can use the `SendMessageStream` method:
```go
// Create a channel to receive events (buffer size 10 as an example)
events := make(chan events, 10)

err := client.SendMessageStream(types.MessageSendParam{
    Message: &types.Message{
        TaskID: taskID,           // The ID of the task this message belongs to
        Role:   types.User,       // The sender's role (e.g., User, Agent)
        Parts: []types.Part{
            // Message content parts; here we use a text message as an example
            &types.TextPart{Kind: "text", Text: message},
        },
    },
}, events)
```

> The above code demonstrates how to use `SendMessageStream` to receive server responses as a stream. Events will be sent to the `events` channel as they arrive. This is useful for real-time or incremental message processing.

## Server

An a2a-server essentially consists of four components: taskStore, executor, queueManager, and updater.
- **taskStore**: Used to store and update task information and status.
- **queueManager**: Manages the creation and destruction of queues related to tasks.
- **executor**: (explained below)
- **updater**: Assists with task status tracking and updates.

Example of basic server setup:

```go
store := tasks.NewInMemoryTaskStore()
manager := a2a.NewQueueManager()

defaultHandler := handler.NewDefaultHandler(store, a2a.NewExecutor(), handler.WithQueueManger(manager))
server := handler.NewServer("/card", "/api", agentCard, defaultHandler)
server.Start(8080)
```

The Executor module provides two core functions: `Execute` and `Cancel`.

- The `Execute` function is responsible for executing the specified task based on the user-provided context.
- The `Cancel` function cancels the execution of the corresponding task.

During execution, the function uses a built-in status update mechanism to provide real-time feedback on task progress. The implementation includes:

1. **Status update mechanism**: Uses the `updater` utility to track and update task status.
2. **Message queue interaction**: All status changes and progress information are written to the event queue (`event.Queue`) in real time.
3. **Result return**: The final execution result is returned to the caller via the queue.

Typical implementation example:

```go
func (e *Executor) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
    // Initialize the task status tracker
    u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)

    // Create the initial status message
    message := u.NewAgentMessage([]types.Part{
        &types.TextPart{Text: "start the work", Kind: types.PartTypeText},
    })

    // Update task status
    u.StartWork(updater.WithMessage(message))
    u.Complete()
    
    return nil
}
```

## Documentation

- [GoDoc](https://pkg.go.dev/github.com/yeeaiclub/a2a-go) - API documentation
- [A2A Specification](https://github.com/a2aproject/A2A) - Protocol specification


## Contributing

We welcome contributions! Please read our contributing guidelines before submitting pull requests.

> **Note:** Since this project is under active development, we recommend checking the latest issues and discussions before contributing to understand the current development priorities.

## License

This project is licensed under the [Apache 2.0 License](LICENSE).

---
**Keywords**: a2a-go, a2a, agent-to-agent, Go SDK, AI agents, protocol implementation, middleware, authentication, streaming, task management