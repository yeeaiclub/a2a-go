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

	"github.com/google/uuid"
	"github.com/yeeaiclub/a2a-go/internal/errs"
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/execution"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/aggregator"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// Handler a2a request handler interface
type Handler interface {
	OnGetTask(ctx context.Context, params types.TaskQueryParams) (*types.Task, error)
	// OnMessageSend start the agent execution for the message and waits for the final result
	OnMessageSend(ctx context.Context, params types.MessageSendParam) (types.Event, error)
	// OnMessageSendStream start the agent execution and yields events
	OnMessageSendStream(ctx context.Context, params types.MessageSendParam) <-chan types.StreamEvent
	// OnCancelTask attempts to cancel the task manged by the agentExecutor
	OnCancelTask(ctx context.Context, params types.TaskIdParams) (*types.Task, error)
	OnSetTaskPushNotificationConfig(ctx context.Context, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error)
	OnGetTaskPushNotificationConfig(ctx context.Context, params types.TaskIdParams) (*types.TaskPushNotificationConfig, error)
	OnResubscribeToTask(ctx context.Context, params types.TaskIdParams) <-chan types.StreamEvent
}

type DefaultHandler struct {
	manger           *manager.TaskManager
	store            tasks.TaskStore
	queueManger      event.QueueManager
	executor         execution.AgentExecutor
	resultAggregator *aggregator.ResultAggregator
	pushNotifier     tasks.PushNotifier
}

func NewDefaultHandler(store tasks.TaskStore, executor execution.AgentExecutor, opts ...HandlerOption) *DefaultHandler {
	handler := &DefaultHandler{store: store, executor: executor}
	for _, opt := range opts {
		opt.Option(handler)
	}

	return handler
}

func (d *DefaultHandler) OnGetTask(ctx context.Context, params types.TaskQueryParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (d *DefaultHandler) OnMessageSend(ctx context.Context, params types.MessageSendParam) (types.Event, error) {
	taskManager := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(params.Message.TaskID),
		manager.WithContextId(params.Message.ContextID),
		manager.WithInitMessage(params.Message),
	)

	task, err := d.store.Get(ctx, params.Message.TaskID)
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

	if task == nil {
		task = &types.Task{Id: uuid.New().String()}
	}

	reqContext := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(task.Id),
		execution.WithContextId(params.Message.ContextID),
		execution.WithTask(task),
	)

	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
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

	if ev.EventType() == "task" && ev.GetTaskId() != task.Id {
		return nil, errs.ErrTaskIdMissingMatch
	}
	return ev, nil
}

func (d *DefaultHandler) OnMessageSendStream(ctx context.Context, params types.MessageSendParam) <-chan types.StreamEvent {
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
	if task == nil {
		task = &types.Task{Id: uuid.New().String()}
	}

	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		return errorStream(err)
	}

	reqCtx := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(task.Id),
		execution.WithContextId(task.ContextId),
		execution.WithTask(task),
	)

	d.execute(ctx, reqCtx, queue)

	resultAggregator := aggregator.NewResultAggregator(taskManager, aggregator.WithBatchSize(10))
	return resultAggregator.ConsumeAndEmit(ctx, queue)
}

// OnCancelTask attempts to cancel the task manged by agentExecutor
func (d *DefaultHandler) OnCancelTask(ctx context.Context, params types.TaskIdParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
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

	reqCtx := execution.NewRequestContext()
	d.cancel(ctx, reqCtx, queue)
	result, err := rg.ConsumeAll(ctx, queue)
	if err != nil {
		return nil, err
	}

	if result.EventType() == "task" {
		return result.(*types.Task), nil
	}
	return nil, errs.ErrInValidResponse
}

func (d *DefaultHandler) OnSetTaskPushNotificationConfig(ctx context.Context, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error) {
	if d.pushNotifier == nil {
		return nil, errs.ErrUnsupportedOperation
	}
	params.TaskId = uuid.New().String()

	err := d.pushNotifier.SetInfo(ctx, params.TaskId, params.Config)
	if err != nil {
		return nil, err
	}
	return &params, nil
}

func (d *DefaultHandler) OnGetTaskPushNotificationConfig(ctx context.Context, params types.TaskIdParams) (*types.TaskPushNotificationConfig, error) {
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

	if task.Id == "" {
		return nil, errs.ErrTaskNotFound
	}

	config, err := d.pushNotifier.GetInfo(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &types.TaskPushNotificationConfig{TaskId: params.Id, Config: config}, nil
}

func (d *DefaultHandler) OnResubscribeToTask(ctx context.Context, params types.TaskIdParams) <-chan types.StreamEvent {
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
	if task == nil || task.Id == "" {
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

func (d *DefaultHandler) execute(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) {
	go func() {
		defer queue.Close()
		err := d.executor.Execute(ctx, reqCtx, queue)
		if err != nil {
			queue.EnqueueError(err)
		}
	}()
}

func (d *DefaultHandler) cancel(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) {
	go func() {
		err := d.executor.Cancel(ctx, reqCtx, queue)
		if err != nil {
			queue.EnqueueError(err)
		}
	}()
}

func (d *DefaultHandler) shouldAddPushInfo(params types.MessageSendParam) bool {
	return d.pushNotifier != nil &&
		params.Configuration != nil &&
		params.Configuration.PushNotificationConfig != nil
}

func (d *DefaultHandler) IsTerminalTaskSates(state types.TaskState) bool {
	return state == types.COMPLETED ||
		state == types.CANCELED ||
		state == types.FAILED ||
		state == types.REJECTED
}

type HandlerOption interface {
	Option(d *DefaultHandler)
}

type HandlerOptionFunc func(d *DefaultHandler)

func (fn HandlerOptionFunc) Option(d *DefaultHandler) {
	fn(d)
}

func WithTaskManger(taskManger *manager.TaskManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.manger = taskManger
	})
}

func WithQueueManger(queueManger event.QueueManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.queueManger = queueManger
	})
}

func WithResultAggregator(rg *aggregator.ResultAggregator) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.resultAggregator = rg
	})
}

func WithPushNotifier(pushNotifier tasks.PushNotifier) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.pushNotifier = pushNotifier
	})
}
