package tasks

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/types"
)

// ResultAggregator is used to process the event streams from an AgentExecutor
type ResultAggregator struct {
	Manager *TaskManager
	Message *types.Message
}

func NewResultAggregator(taskManger *TaskManager, initMessage *types.Message) *ResultAggregator {
	return &ResultAggregator{Manager: taskManger, Message: initMessage}
}

// ConsumeAndEmit process the event stream from the consumer, updates the task store
func (r *ResultAggregator) ConsumeAndEmit(ctx context.Context, consumer *event.Consumer) <-chan types.StreamEvent {
	events := make(chan types.StreamEvent, 1024)
	go func() {
		defer close(events)
		for {
			select {
			case <-ctx.Done():
				events <- types.StreamEvent{Err: ctx.Err()}
				return
			case e, ok := <-consumer.ConsumeAll(ctx):
				if !ok {
					return
				}
				if e.Err != nil {
					events <- types.StreamEvent{Err: e.Err}
					return
				}
				if e.Done() {
					events <- e
					return
				}
				ev, err := r.Manager.Process(ctx, e.Event)
				if err != nil {
					events <- types.StreamEvent{Err: err}
					return
				}
				events <- types.StreamEvent{Event: ev}
			}
		}
	}()
	return events
}

// ConsumeAll process the entire event stream from consumer and return the final result.
func (r *ResultAggregator) ConsumeAll(ctx context.Context, consumer *event.Consumer) (types.Event, error) {
	events := consumer.ConsumeAll(ctx)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case ev, ok := <-events:
			if !ok {
				return nil, nil
			}
			if ev.Err != nil {
				return nil, ev.Err
			}
			if ev.Done() {
				task, err := r.Manager.GetTask(ctx)
				if err != nil {
					return nil, err
				}
				return task, nil
			}
			if ev.EventType() == "message" {
				msg := ev.Event.(*types.Message)
				r.Message = msg
				return msg, nil
			}
			_, err := r.Manager.Process(ctx, ev.Event)
			if err != nil {
				return nil, err
			}
		}
	}
}

// ConsumeAndBreakOnInterrupt process the event stream until completion or an interruptible state is encountered
func (r *ResultAggregator) ConsumeAndBreakOnInterrupt(ctx context.Context, consumer *event.Consumer) (types.Event, bool, error) {
	events := consumer.ConsumeAll(ctx)
	var result types.Event

	for {
		select {
		case <-ctx.Done():
			return nil, false, ctx.Err()
		case e, ok := <-events:
			if !ok {
				task, err := r.Manager.GetTask(ctx)
				result = task
				return result, false, err
			}
			if e.Err != nil {
				return nil, false, e.Err
			}
			if e.Done() {
				task, err := r.Manager.GetTask(ctx)
				result = task
				return result, false, err
			}
			_, err := r.Manager.Process(ctx, e.Event)
			if err != nil {
				return nil, false, err
			}

			if e.EventType() == "update_event" {
				updateEvent := e.Event.(*types.TaskStatusUpdateEvent)
				if updateEvent.Status.State == types.AUTH_REQUIRED {
					r.continueConsume(ctx, events)
					return nil, true, err
				}
			}

			if e.EventType() == "task" {
				taskEvent := e.Event.(*types.Task)
				if taskEvent.Status.State == types.AUTH_REQUIRED {
					r.continueConsume(ctx, events)
					return nil, true, err
				}
			}
		}
	}
}

// continueConsume processing an event stream in backhand task
func (r *ResultAggregator) continueConsume(ctx context.Context, events <-chan types.StreamEvent) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev, ok := <-events:
				if !ok {
					return
				}
				if ev.Done() {
					return
				}
				_, err := r.Manager.Process(ctx, ev.Event)
				if err != nil {
					return
				}
			}
		}
	}()
}
