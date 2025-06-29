// Copyright 2025 yumosx
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
	"github.com/yumosx/a2a-go/internal/errs"
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/server/execution"
	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/server/tasks/aggregator"
	"github.com/yumosx/a2a-go/sdk/server/tasks/manager"
	"github.com/yumosx/a2a-go/sdk/types"
)

// Handler a2a request handler interface
type Handler interface {
	OnGetTask(ctx context.Context, params types.TaskQueryParams) (*types.Task, error)
	OnMessageSend(ctx context.Context, params types.MessageSendParam) (types.Event, error)
	OnMessageSendStream(ctx context.Context, params types.MessageSendParam) <-chan types.StreamEvent
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
	taskManger := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(params.Message.TaskID),
		manager.WithContextId(params.Message.ContextID),
		manager.WithInitMessage(params.Message),
	)

	task, err := d.getTask(ctx, params.Message.TaskID)
	if err != nil {
		return nil, err
	}

	if d.IsTerminalTaskSates(task.Status.State) {
		return nil, fmt.Errorf("task %s is in terminal state: %s", task.Id, task.Status.State)
	}

	if d.shouldAddPushInfo(params) {
		err = d.pushNotifier.SetInfo(ctx, task.Id, params.Configuration.PushNotificationConfig)
		if err != nil {
			return nil, err
		}
	}

	task = taskManger.UpdateWithMessage(params.Message, task)

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
	defer queue.Close()

	consumer := event.NewConsumer(queue, nil)
	err = d.executor.Execute(ctx, reqContext, queue)
	if err != nil {
		return nil, err
	}

	resultAggregator := aggregator.NewResultAggregator(taskManger)
	ev, err := resultAggregator.ConsumeAndBreakOnInterrupt(ctx, consumer)
	if err != nil {
		return nil, err
	}

	if ev.EventType() == "task" && ev.GetTaskId() != task.Id {
		return nil, errs.ErrTaskIdMissingMatch
	}
	return ev, nil
}

func (d *DefaultHandler) OnMessageSendStream(ctx context.Context, params types.MessageSendParam) <-chan types.StreamEvent {
	ch := make(chan types.StreamEvent, 1)

	taskManger := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(params.Message.TaskID),
		manager.WithContextId(params.Message.ContextID),
		manager.WithInitMessage(params.Message),
	)

	task, err := taskManger.GetTask(ctx)
	if err != nil {
		ch <- types.StreamEvent{Err: err}
		return ch
	}
	if task == nil {
		ch <- types.StreamEvent{Err: errs.ErrTaskNotFound}
		return ch
	}
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		ch <- types.StreamEvent{Err: err}
		return ch
	}

	reqCtx := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(task.Id),
		execution.WithContextId(task.ContextId),
		execution.WithTask(task),
	)

	errCh := d.execute(ctx, reqCtx, queue)
	consumer := event.NewConsumer(queue, errCh)

	rg := aggregator.NewResultAggregator(taskManger)
	events := rg.ConsumeAndEmit(ctx, consumer)
	return events
}

// OnCancelTask attempts to cancel the task manged by agentExecutor
func (d *DefaultHandler) OnCancelTask(ctx context.Context, params types.TaskIdParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	taskManger := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(task.Id),
		manager.WithContextId(task.ContextId),
	)

	rg := aggregator.NewResultAggregator(taskManger)
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		return nil, err
	}
	if queue == nil {
		queue = event.NewQueue(0)
		defer queue.Close()
	}
	defer queue.Close()

	reqCtx := execution.NewRequestContext()
	done := d.cancel(ctx, reqCtx, queue)
	consumer := event.NewConsumer(queue, done)
	result, err := rg.ConsumeAll(ctx, consumer)
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
		return nil, errs.ErrUnSupportedOperation
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
		return nil, errs.ErrUnSupportedOperation
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
	errCh := make(chan types.StreamEvent, 1)
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		errCh <- types.StreamEvent{Err: err}
		return errCh
	}

	if task == nil {
		errCh <- types.StreamEvent{Err: errs.ErrTaskNotFound}
		return errCh
	}

	manger := manager.NewTaskManger(
		d.store,
		manager.WithTaskId(task.Id),
		manager.WithContextId(task.ContextId),
	)
	rg := aggregator.NewResultAggregator(manger)
	queue, err := d.queueManger.Tap(ctx, task.Id)
	if err != nil {
		errCh <- types.StreamEvent{Err: err}
		return errCh
	}
	consumer := event.NewConsumer(queue, nil)
	return rg.ConsumeAndEmit(ctx, consumer)
}

func (d *DefaultHandler) execute(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- d.executor.Execute(ctx, reqCtx, queue)
	}()
	return ch
}

func (d *DefaultHandler) cancel(ctx context.Context, reqCtx *execution.RequestContext, queue *event.Queue) chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		ch <- d.executor.Cancel(ctx, reqCtx, queue)
	}()
	return ch
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

func (d *DefaultHandler) getTask(ctx context.Context, taskId string) (*types.Task, error) {
	task, err := d.store.Get(ctx, taskId)
	if err != nil {
		return nil, err
	}
	if task == nil || task.Id == "" {
		return nil, errs.ErrTaskNotFound
	}
	return task, nil
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
