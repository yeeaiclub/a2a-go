package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yumosx/a2a-go/internal/errs"
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/server/execution"
	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/types"
)

// Handler a2a request handler interface
type Handler interface {
	OnGetTask(ctx context.Context, params types.TaskQueryParams) (*types.Task, error)
	OnMessageSend(ctx context.Context, params types.MessageSendParam) (types.Event, error)
	OnMessageSendStream(ctx context.Context, params types.MessageSendParam) chan types.StreamEvent
	OnCancelTask(ctx context.Context, params types.TaskIdParams) (*types.Task, error)
	OnSetTaskPushNotificationConfig(ctx context.Context, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error)
	OnGetTaskPushNotificationConfig(ctx context.Context, params types.TaskIdParams) (*types.TaskPushNotificationConfig, error)
	OnResubscribeToTask(ctx context.Context, params types.TaskIdParams) chan types.StreamEvent
}

type DefaultHandler struct {
	manger           *tasks.TaskManager
	store            tasks.TaskStore
	queueManger      event.QueueManager
	executor         execution.AgentExecutor
	resultAggregator *tasks.ResultAggregator
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
	manger := tasks.NewTaskManger(
		d.store,
		tasks.WithTaskId(params.Message.TaskID),
		tasks.WithContextId(params.Message.ContextID),
		tasks.WithInitMessage(params.Message),
	)
	task, err := manger.GetTask(ctx)
	if err != nil {
		return nil, err
	}
	if task == nil || task.Id == "" {
		return nil, errs.TaskNotFound
	}
	if task.Status.State == types.COMPLETED {
		return nil, fmt.Errorf("task %s is in terminal state: %s", task.Id, task.Status.State)
	}
	task = manger.UpdateWithMessage(params.Message, task)

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

	resultAggregator := tasks.NewResultAggregator(manger, nil)
	ev, _, err := resultAggregator.ConsumeAndBreakOnInterrupt(ctx, consumer)
	if err != nil {
		return nil, err
	}

	if ev.EventType() == "task" && ev.GetTaskId() != task.Id {
		return nil, errors.New("task ID mismatch in agent response")
	}
	return ev, nil
}

func (d *DefaultHandler) OnMessageSendStream(ctx context.Context, params types.MessageSendParam) <-chan types.StreamEvent {
	ch := make(chan types.StreamEvent, 1)

	manger := tasks.NewTaskManger(
		d.store,
		tasks.WithTaskId(params.Message.TaskID),
		tasks.WithContextId(params.Message.ContextID),
		tasks.WithInitMessage(params.Message),
	)

	task, err := manger.GetTask(ctx)
	if err != nil {
		ch <- types.StreamEvent{Err: err}
		return ch
	}
	if task == nil {
		ch <- types.StreamEvent{Err: errs.TaskNotFound}
		return ch
	}
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		ch <- types.StreamEvent{Err: err}
		return ch
	}
	defer queue.Close()

	rg := tasks.NewResultAggregator(manger, nil)
	reqCtx := execution.NewRequestContext(
		execution.WithParams(params),
		execution.WithTaskId(task.Id),
		execution.WithContextId(task.ContextId),
		execution.WithTask(task),
	)

	errCh := d.execute(ctx, reqCtx, queue)
	consumer := event.NewConsumer(queue, errCh)

	return rg.ConsumeAndEmit(ctx, consumer)
}

// CancelTask attempts to cancel the task manged by agentExecutor
func (d *DefaultHandler) CancelTask(ctx context.Context, params types.TaskIdParams) (*types.Task, error) {
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	manger := tasks.NewTaskManger(
		d.store,
		tasks.WithTaskId(task.Id),
		tasks.WithContextId(task.ContextId),
	)

	rg := tasks.NewResultAggregator(manger, nil)
	queue, err := d.queueManger.CreateOrTap(ctx, task.Id)
	if err != nil {
		return nil, err
	}
	defer queue.Close()

	if queue == nil {
		queue = event.NewQueue(0)
	}

	reqCtx := execution.NewRequestContext()
	done := d.execute(ctx, reqCtx, queue)
	consumer := event.NewConsumer(queue, done)
	result, err := rg.ConsumeAll(ctx, consumer)
	if err != nil {
		return nil, err
	}

	if result.EventType() == "task" {
		return result.(*types.Task), nil
	}
	return nil, errors.New("agent did not return valid response for cancel")
}

func (d *DefaultHandler) OnSetTaskPushNotificationConfig(ctx context.Context, params types.TaskPushNotificationConfig) (*types.TaskPushNotificationConfig, error) {
	if d.pushNotifier == nil {
		return nil, errs.UnSupportedOperation
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
		return nil, errs.UnSupportedOperation
	}

	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errs.TaskNotFound
	}

	if task.Id == "" {
		return nil, errs.TaskNotFound
	}

	config, err := d.pushNotifier.GetInfo(ctx, params.Id)
	if err != nil {
		return nil, err
	}
	return &types.TaskPushNotificationConfig{TaskId: params.Id, Config: config}, nil
}

func (d *DefaultHandler) OnReSubscribeToTask(ctx context.Context, params types.TaskIdParams) <-chan types.StreamEvent {
	errCh := make(chan types.StreamEvent, 1)
	task, err := d.store.Get(ctx, params.Id)
	if err != nil {
		errCh <- types.StreamEvent{Err: err}
		return errCh
	}

	if task == nil {
		errCh <- types.StreamEvent{Err: errs.TaskNotFound}
		return errCh
	}

	manger := tasks.NewTaskManger(
		d.store,
		tasks.WithTaskId(task.Id),
		tasks.WithContextId(task.ContextId),
	)
	rg := tasks.NewResultAggregator(manger, nil)
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

type HandlerOption interface {
	Option(d *DefaultHandler)
}

type HandlerOptionFunc func(d *DefaultHandler)

func (fn HandlerOptionFunc) Option(d *DefaultHandler) {
	fn(d)
}

func WithTaskManger(taskManger *tasks.TaskManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.manger = taskManger
	})
}

func WithQueueManger(queueManger event.QueueManager) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.queueManger = queueManger
	})
}

func WithResultAggregator(rg *tasks.ResultAggregator) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.resultAggregator = rg
	})
}

func WithPushNotifier(pushNotifier tasks.PushNotifier) HandlerOption {
	return HandlerOptionFunc(func(d *DefaultHandler) {
		d.pushNotifier = pushNotifier
	})
}
