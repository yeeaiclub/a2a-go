package execution

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
)

type AgentExecutor interface {
	Execute(ctx context.Context, reqContext RequestContext, queue *event.Queue)
	Cancel(ctx context.Context, requestContext RequestContext, queue *event.Queue)
}
