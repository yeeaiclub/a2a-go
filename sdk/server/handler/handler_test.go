package handler

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/server/execution"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e Executor) Execute(ctx context.Context, reqContext *execution.RequestContext, queue *event.Queue) error {
	return nil
}

func (e Executor) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	return nil
}
