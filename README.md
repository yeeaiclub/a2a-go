# a2a-go

Agent-to-Agent Protocol Implementation for Go

## how to use?

### 1. define Agent Card:

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
### 2. define agent executor



## install