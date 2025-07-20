# a2a-go：面向 Go 的 Agent-to-Agent 协议 SDK

[![GoDoc](https://pkg.go.dev/badge/github.com/yeeaiclub/a2a-go)](https://pkg.go.dev/github.com/yeeaiclub/a2a-go)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Codecov](https://img.shields.io/codecov/c/github/yeeaiclub/a2a-go/main?logo=codecov&logoColor=white)](https://codecov.io/gh/yeeaiclub/a2a-go/branch/main)
[![Status](https://img.shields.io/badge/Status-Under%20Development-orange.svg)](https://github.com/yeeaiclub/a2a-go)

**a2a-go** 是一套完整的 **Agent-to-Agent (A2A) 协议** Go 语言 SDK 实现，为构建 AI Agent 通信系统提供了坚实的基础。该 SDK 支持中间件、认证和完整协议方法，助力 AI Agent 之间的无缝集成。

> **⚠️ 本项目正在积极开发中，API 可能随时变更。**

## 概述

**a2a-go** SDK 基于 [A2A 协议规范](https://github.com/a2aproject/A2A) 实现，主要特性包括：

### 核心 A2A 协议方法
- [x] **认证** - 安全的 Agent 认证
- [x] **发送 message** - 向其他 Agent 发送消息
- [x] **获取任务** - 查询任务信息与状态
- [x] **取消任务** - 取消运行中或待处理的任务
- [x] **任务流式传输** - 实时获取任务结果流
- [x] **设置推送通知** - 配置任务推送通知
- [x] **获取推送通知** - 查询推送通知配置

### SDK 特性
- [x] **中间件支持** - 可扩展的请求/响应中间件架构
- [x] **多种安全方案** - 支持 API Key、Bearer、OAuth2、OpenID Connect 等多种认证方式
- [x] **上下文管理** - 灵活的上下文与安全配置管理


## 安装

```shell
go get github.com/yeeaiclub/a2a-go
```

## 快速开始

## 客户端示例

// 初始化带有自定义超时的 HTTP 客户端（可根据需要调整 config.Timeout）
```go
httpClient := &http.Client{ Timeout: config.Timeout }

// 使用 HTTP 客户端和 API 地址创建 a2a-go 客户端实例
// 请将 "http://localhost:8080/api" 替换为你的实际服务端地址
a2aClient := client.NewClient(httpClient, "http://localhost:8080/api")
```

> 上述代码演示了如何初始化 a2a-go 客户端。你需要传入自定义的 `http.Client`（可设置超时、代理等参数）以及 a2a 服务端的 API 地址。

// 使用 a2a-go 客户端发送消息示例
```go
resp, err := client.SendMessage(types.MessageSendParam{
    Message: &types.Message{
        TaskID: taskID,           // 消息所属的任务 ID
        Role:   types.User,       // 发送方角色（如 User、Agent）
        Parts: []types.Part{
            // 消息内容部分，这里以文本消息为例
            &types.TextPart{Kind: "text", Text: message},
        },
    },
})
```

> 上述代码演示了如何发送一条消息：
> - `TaskID`：指定消息关联的任务。**如果是新建任务，该字段可以为空，服务端会自动生成 TaskID。**
> - `Role`：指定消息发送者的角色（如 User、Agent）。
> - `Parts`：支持多种消息内容（如文本、图片等），此处为文本消息。
> 
> `resp` 变量包含服务端响应，`err` 用于错误处理。

如果你想以流式的方式来从服务端返回对应的消息，你可以使用 SendMessageStream 这个方法:
```go
events := make(chan events, 10)

err := client.SendMessageStream(types.MessageSendParam{
    Message: &types.Message{
        TaskID: taskID,           // 消息所属的任务 ID
        Role:   types.User,       // 发送方角色（如 User、Agent）
        Parts: []types.Part{
            // 消息内容部分，这里以文本消息为例
            &types.TextPart{Kind: "text", Text: message},
        },
    },
}, events)
```


## 服务端

一个 a2a-server 其实是包含了四个部分 taskStore，executor， queueManager， updater
其中 taskStore 是用来存储和更新对应的 task 的信息和 状态，queueManager 用来管理和 task 相关的队列创建和销毁的，而 executor 我们稍后介绍

```go
store := tasks.NewInMemoryTaskStore()
manager := a2a.NewQueueManager()

defaultHandler := handler.NewDefaultHandler(store, a2a.NewExecutor(), handler.WithQueueManger(manager))
server := handler.NewServer("/card", "/api", agentCard, defaultHandler)
server.Start(8080)
```

Executor 模块提供了两个核心函数，execute和 cancel

其中 Execute 函数负责根据用户提供的上下文执行指定任务。而cancel 则是取消对应的 task 的执行

在执行过程中，该函数会通过内置的状态更新机制实时反馈任务进度，具体实现如下：

1. 状态更新机制：通过 updater 辅助工具类实现任务状态的追踪和更新
2. 消息队列交互：所有状态变更和进度信息都会实时写入事件队列（event.Queue）
3. 结果返回：最终执行结果通过队列返回给调用方

典型实现示例：

```go
func (e *Executor) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
    // 初始化任务状态追踪器
    u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)

    // 创建初始状态消息
    message := u.NewAgentMessage([]types.Part{
        &types.TextPart{Text: "start the work", Kind: types.PartTypeText},
    })

    // 更新任务状态
    u.StartWork(updater.WithMessage(message))
    u.Complete()
    
    return nil
}
```

## 文档

- [GoDoc](https://pkg.go.dev/github.com/yeeaiclub/a2a-go) - API 文档
- [A2A 规范](https://github.com/a2aproject/A2A) - 协议规范


## 贡献

欢迎贡献代码！在提交 Pull Request 前请阅读贡献指南。

> **注意：** 由于本项目处于活跃开发阶段，建议在贡献前查看最新的 issues 和讨论，以了解当前开发重点。

## 许可证

本项目基于 [Apache 2.0 License](LICENSE) 开源。

---
**关键词**：a2a-go, a2a, agent-to-agent, Go SDK, AI agents, 协议实现, 中间件, 认证, 流式传输, 任务管理

