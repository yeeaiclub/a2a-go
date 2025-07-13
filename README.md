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

1. [demo](https://github.com/yeeaiclub/demo)
2. [agent-platform](https://github.com/yeeaiclub/agent-platform)



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