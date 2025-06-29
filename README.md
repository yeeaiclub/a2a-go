# a2a-go

Agent-to-Agent Protocol Implementation for Go

## how to use?

### server
#### 1. define Agent Card:

```go
var mockAgentCard = types.AgentCard{
	Name:        "test agent",
	Description: "a test agent for test",
	URL:         "http://localhost:8080/",
	Version:     "1.0",
	Capabilities: &types.AgentCapabilities{
		Streaming:              true,
		StateTransitionHistory: true,
	},
	Skills: []types.AgentSkill{
		{
			Id:          "test-skills",
			Name:        "test skill",
			Description: "a test skill for unit test",
			InputModes:  []string{"text/plain"},
		},
	},
}
```
#### 2. define agent executor

```go
type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	u.Complete(updater.WithFinal(true))
	return nil
}

func (e *Executor) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	u.Complete(updater.WithFinal(true))
	return nil
}
```
#### 3. create a task store
#### 4. start a server

### client

## install