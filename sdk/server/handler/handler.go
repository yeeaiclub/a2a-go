// Copyright 2025 yeeaiclub
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"fmt"

	"github.com/yeeaiclub/a2a-go/internal/errs"
	"github.com/yeeaiclub/a2a-go/sdk/server"
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/execution"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/aggregator"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// Handler defines the interface for handling all agent API requests.
type Handler interface {
	// OnGetTask retrieves a task by its ID.
	OnGetTask(ctx *server.CallContext, params types.TaskQueryParams) (*types.Task, error)
	// OnMessageSend starts the agent execution for the message and waits for the final result.
	OnMessageSend(ctx *server.CallContext, params types.MessageSendParam) (types.Event, error)
	// OnMessageSendStream starts the agent execution and yields events as a stream.
	OnMessageSendStream(ctx *server.CallContext, params types.MessageSendParam) <-chan types.StreamEvent
	// OnCancelTask attempts to cancel the task managed by the agent executor.
	OnCancelTask(ctx *server.CallContext, params types.TaskIdParams) (*types.Task, error)
	// OnSetTaskPushNotificationConfig sets the push notification configuration for a task.
	OnSetTaskPushNotificationConfig(ctx *server.CallContext, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error)
	// OnGetTaskPushNotificationConfig retrieves the push notification configuration for a task.
	OnGetTaskPushNotificationConfig(ctx *server.CallContext, params types.TaskIdParams) (*types.TaskPushNotificationConfig, error)
	// OnResubscribeToTask resubscribes to task events and returns a stream of events.
	OnResubscribeToTask(ctx *server.CallContext, params types.TaskIdParams) <-chan types.StreamEvent
}

// DefaultHandler provides a default implementation of the Handler interface.
type DefaultHandler struct {
	manger           *manager.TaskManager         // Task manager for task lifecycle
	store            tasks.TaskStore              // Task storage backend
	queueManger      event.QueueManager           // Event queue manager
	executor         execution.AgentExecutor      // Agent execution engine
	resultAggregator *aggregator.ResultAggregator // Aggregates results from event queue
	pushNotifier     tasks.PushNotifier           // Push notification handler
}

// NewDefaultHandler creates a new DefaultHandler with optional configuration.
func NewDefaultHandler(store tasks.TaskStore, executor execution.AgentExecutor, opts ...HandlerOption) *DefaultHandler {
	handler := &DefaultHandler{store: store, executor: executor}
	for _, opt := range opts {
		opt.Option(handler)
	}

	return handler
}

// OnGetTask handles task retrieval requests.
func (d *DefaultHandler) OnGetTask(ctx *server.CallContext, params types.TaskQueryParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// OnMessageSend handles synchronous message send requests and waits for the result.
func (d *DefaultHandler) OnMessageSend(ctx *server.CallContext, params types.MessageSendParam) (types.Event, error) {
	taskManager := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(params.Message.TaskID),
		manager.WithContextId(params.Message.ContextID),
		manager.WithInitMessage(params.Message),
	)

	task, err := taskManager.GetTask(ctx)
	if err != nil {
		return nil, err
	}

	if task != nil {
		if d.IsTerminalTaskSates(task.Status.State) {
			return nil, fmt.Errorf("task %s is in terminal state: %s", task.Id, task.Status.State)
		}
		task = taskManager.UpdateWithMessage(params.Message, task)
		if d.shouldAddPushInfo(params) {
			err = d.pushNotifier.SetInfo(ctx, task.Id, params.Configuration.PushNotificationConfig)
			if err != nil {
				return nil, err
			}
		}
	}

	reqContext, err := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(params.Message.TaskID),
		execution.WithContextId(params.Message.ContextID),
		execution.WithTask(task),
		execution.WithServerContext(ctx),
	)
	if err != nil {
		return nil, err
	}

	queue, err := d.queueManger.CreateOrTap(ctx, reqContext.TaskId)
	if err != nil {
		return nil, err
	}

	d.execute(ctx, reqContext, queue)
	resultAggregator := aggregator.NewResultAggregator(
		taskManager,
		aggregator.WithBatchSize(10),
	)
	ev, err := resultAggregator.ConsumeAndBreakOnInterrupt(ctx, queue)
	if err != nil {
		return nil, err
	}

	if ev != nil && ev.Type() == "task" && ev.GetTaskId() != reqContext.TaskId {
		return nil, errs.ErrTaskIdMissingMatch
	}
	return ev, nil
}

// OnMessageSendStream handles streaming message send requests, returning a channel of events.
func (d *DefaultHandler) OnMessageSendStream(ctx *server.CallContext, params types.MessageSendParam) <-chan types.StreamEvent {
	errorStream := func(err error) <-chan types.StreamEvent {
		ch := make(chan types.StreamEvent, 1)
		ch <- types.StreamEvent{Type: types.EventError, Err: err}
		close(ch)
		return ch
	}

	taskManager := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(params.Message.TaskID),
		manager.WithContextId(params.Message.ContextID),
		manager.WithInitMessage(params.Message),
	)

	task, err := taskManager.GetTask(ctx)
	if err != nil {
		return errorStream(err)
	}

	reqContext, err := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(params.Message.TaskID),
		execution.WithContextId(params.Message.ContextID),
		execution.WithTask(task),
		execution.WithServerContext(ctx),
	)
	if err != nil {
		return errorStream(err)
	}

	queue, err := d.queueManger.CreateOrTap(ctx, reqContext.TaskId)
	if err != nil {
		return errorStream(err)
	}

	d.execute(ctx, reqContext, queue)

	resultAggregator := aggregator.NewResultAggregator(taskManager, aggregator.WithBatchSize(10))
	return resultAggregator.ConsumeAndEmit(ctx, queue)
}

// OnCancelTask handles task cancellation requests.
func (d *DefaultHandler) OnCancelTask(ctx *server.CallContext, params types.TaskIdParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, errs.ErrTaskNotFound
	}

	if task.Id == "" {
		return nil, errs.ErrTaskNotFound
	}

	taskManager := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(task.Id),
		manager.WithContextId(task.ContextId),
	)

	rg := aggregator.NewResultAggregator(taskManager, aggregator.WithBatchSize(10))
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		return nil, err
	}
	if queue == nil {
		queue = event.NewQueue(0)
		defer queue.Close()
	}

	reqCtx, err := execution.NewRequestContext(
		execution.WithTaskId(task.Id),
		execution.WithContextId(task.ContextId),
		execution.WithTask(task),
	)
	if err != nil {
		return nil, err
	}

	d.cancel(ctx, reqCtx, queue)
	result, err := rg.ConsumeAll(ctx, queue)
	if err != nil {
		return nil, err
	}

	if result.Type() == "task" {
		return result.(*types.Task), nil
	}
	return nil, errs.ErrInValidResponse
}

// OnSetTaskPushNotificationConfig sets push notification configuration for a task.
func (d *DefaultHandler) OnSetTaskPushNotificationConfig(ctx *server.CallContext, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error) {
	if d.pushNotifier == nil {
		return nil, errs.ErrUnsupportedOperation
	}

	task, err := d.store.Get(ctx, params.TaskId)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errs.ErrTaskNotFound
	}

	err = d.pushNotifier.SetInfo(ctx, params.TaskId, params.Config)
	if err != nil {
		return nil, err
	}
	return &params, nil
}

// OnGetTaskPushNotificationConfig retrieves push notification configuration for a task.
func (d *DefaultHandler) OnGetTaskPushNotificationConfig(ctx *server.CallContext, params types.TaskIdParams) (*types.TaskPushNotificationConfig, error) {
	if d.pushNotifier == nil {
		return nil, errs.ErrUnsupportedOperation
	}

	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errs.ErrTaskNotFound
	}

	config, err := d.pushNotifier.GetInfo(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &types.TaskPushNotificationConfig{TaskId: params.Id, Config: config}, nil
}

// OnResubscribeToTask handles resubscription to task events, returning a channel of events.
func (d *DefaultHandler) OnResubscribeToTask(ctx *server.CallContext, params types.TaskIdParams) <-chan types.StreamEvent {
	errorStream := func(err error) <-chan types.StreamEvent {
		ch := make(chan types.StreamEvent, 1)
		ch <- types.StreamEvent{Type: types.EventError, Err: err}
		close(ch)
		return ch
	}

	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return errorStream(err)
	}
	if task == nil {
		return errorStream(errs.ErrTaskNotFound)
	}

	taskManager := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(task.Id),
		manager.WithContextId(task.ContextId),
	)
	resultAggregator := aggregator.NewResultAggregator(taskManager, aggregator.WithBatchSize(10))
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		return errorStream(err)
	}
	return resultAggregator.ConsumeAndEmit(ctx, queue)
}

// execute runs the agent executor in a goroutine and closes the queue on completion.
func (d *DefaultHandler) execute(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) {
	go func() {
		defer queue.Close()
		err := d.executor.Execute(ctx, reqCtx, queue)
		if err != nil {
			queue.EnqueueError(err)
		}
	}()
}

// cancel requests cancellation of a running task.
func (d *DefaultHandler) cancel(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) {
	go func() {
		err := d.executor.Cancel(ctx, reqCtx, queue)
		if err != nil {
			queue.EnqueueError(err)
		}
	}()
}

// shouldAddPushInfo checks if push notification info should be added for the request.
func (d *DefaultHandler) shouldAddPushInfo(params types.MessageSendParam) bool {
	return d.pushNotifier != nil && params.Configuration != nil && params.Configuration.PushNotificationConfig != nil
}

// IsTerminalTaskSates returns true if the task state is terminal (completed, canceled, failed, rejected).
func (d *DefaultHandler) IsTerminalTaskSates(state types.TaskState) bool {
	return state == types.COMPLETED || state == types.CANCELED || state == types.FAILED || state == types.REJECTED
}

// HandlerOption allows customizing DefaultHandler via functional options.
type HandlerOption interface {
	Option(d *DefaultHandler)
}

// HandlerOptionFunc is a function type for HandlerOption.
type HandlerOptionFunc func(d *DefaultHandler)

func (fn HandlerOptionFunc) Option(d *DefaultHandler) {
	fn(d)
}

// WithTaskManger sets a custom TaskManager for the handler.
func WithTaskManger(taskManger *manager.TaskManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.manger = taskManger
	})
}

// WithQueueManger sets a custom QueueManager for the handler.
func WithQueueManger(queueManger event.QueueManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.queueManger = queueManger
	})
}

// WithResultAggregator sets a custom ResultAggregator for the handler.
func WithResultAggregator(rg *aggregator.ResultAggregator) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.resultAggregator = rg
	})
}

// WithPushNotifier sets a custom PushNotifier for the handler.
func WithPushNotifier(pushNotifier tasks.PushNotifier) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.pushNotifier = pushNotifier
	})
}
