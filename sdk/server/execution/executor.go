package execution

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
)

// AgentExecutor Implementations of this interface contain the core logic of the agent,
// executing tasks based on requests and publishing updates to an event queue.
type AgentExecutor interface {
	// Execute the agent's logic for a given request context
	Execute(ctx context.Context, reqContext *RequestContext, queue *event.Queue) error
	// Cancel request the agent to cancel an ongoing task
	Cancel(ctx context.Context, requestContext *RequestContext, queue *event.Queue) error
}
