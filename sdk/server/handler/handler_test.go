package handler

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/server/execution"
	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/types"
)

type Executor struct {
}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e Executor) Execute(ctx context.Context, reqContext *execution.RequestContext, queue *event.Queue) error {
	updater := tasks.NewTaskUpdater(queue, reqContext.TaskId, reqContext.ContextId)
	parts := []types.Part{&types.TextPart{Kind: "text", Text: "I am yumosx !"}}
	msg := updater.NewAgentMessage(parts)
	updater.UpdateStatus(types.COMPLETED, tasks.WithMessage(&msg), tasks.WithFinal(true))
	return nil
}

func (e Executor) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	return nil
}
