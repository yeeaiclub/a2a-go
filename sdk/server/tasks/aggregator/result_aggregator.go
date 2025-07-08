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

package aggregator

import (
	"context"

	"github.com/yeeaiclub/a2a-go/internal/errs"
	log "github.com/yeeaiclub/a2a-go/internal/logger"
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// ResultAggregator is used to process the event streams from an AgentExecutor
type ResultAggregator struct {
	manager   *manager.TaskManager
	message   *types.Message
	batchSize int
}

func NewResultAggregator(taskManger *manager.TaskManager, options ...ResultAggregatorOption) *ResultAggregator {
	rg := &ResultAggregator{manager: taskManger, batchSize: 10}
	for _, opt := range options {
		opt.Option(rg)
	}
	return rg
}

// ConsumeAndEmit process the event stream from the queue, updates the task store
func (r *ResultAggregator) ConsumeAndEmit(ctx context.Context, queue *event.Queue) <-chan types.StreamEvent {
	out := make(chan types.StreamEvent, r.batchSize)
	go func() {
		defer close(out)
		for e := range queue.Subscribe(ctx) {
			switch e.Type {
			case types.EventClosed:
				return
			case types.EventData:
				_, err := r.manager.Process(ctx, e.Event)
				if err != nil {
					out <- types.StreamEvent{Type: types.EventError, Err: err}
					return
				}
				out <- e
			case types.EventCanceled, types.EventError, types.EventDone:
				out <- e
				return
			}
		}
	}()
	return out
}

// ConsumeAll process the entire event stream from queue and return the final result.
func (r *ResultAggregator) ConsumeAll(ctx context.Context, queue *event.Queue) (types.Event, error) {
	for e := range queue.Subscribe(ctx) {
		switch e.Type {
		case types.EventCanceled:
			return nil, ctx.Err()
		case types.EventError:
			return nil, e.Err
		case types.EventDone:
			task, err := r.manager.GetTask(ctx)
			if err != nil {
				return nil, err
			}
			return task, nil
		case types.EventData:
			if msg, ok := e.Event.(*types.Message); ok {
				r.message = msg
				return msg, nil
			}
			_, err := r.manager.Process(ctx, e.Event)
			if err != nil {
				return nil, err
			}
		case types.EventClosed:
			return nil, nil
		}
	}
	return nil, nil
}

// ConsumeAndBreakOnInterrupt process the event stream until completion or an interruptible state is encountered
func (r *ResultAggregator) ConsumeAndBreakOnInterrupt(ctx context.Context, queue *event.Queue) (types.Event, error) {
	var result types.Event
	for e := range queue.Subscribe(ctx) {
		switch e.Type {
		case types.EventCanceled:
			return nil, ctx.Err()
		case types.EventError:
			return nil, e.Err
		case types.EventClosed:
			task, err := r.manager.GetTask(ctx)
			result = task
			return result, err
		case types.EventDone:
			task, err := r.manager.GetTask(ctx)
			result = task
			return result, err
		case types.EventData:
			if r.IsAuthRequired(e.Event) {
				go r.continueConsume(ctx, queue.Subscribe(ctx))
				return nil, errs.ErrAuthRequired
			}
			_, err := r.manager.Process(ctx, e.Event)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func (r *ResultAggregator) IsAuthRequired(event types.Event) bool {
	if event.EventType() == "status_update" {
		updateEvent := event.(*types.TaskStatusUpdateEvent)
		if updateEvent.Status.State == types.AUTH_REQUIRED {
			log.Debug("Encountered an auth-required task: breaking synchronous message/send flow.")
			return true
		}
	}

	if event.EventType() == "task" {
		taskEvent := event.(*types.Task)
		if taskEvent.Status.State == types.AUTH_REQUIRED {
			log.Debug("Encountered an auth-required task: breaking synchronous message/send flow.")
			return true
		}
	}
	return false
}

// continueConsume processing an event stream in backhand task
func (r *ResultAggregator) continueConsume(ctx context.Context, events <-chan types.StreamEvent) {
	for e := range events {
		if e.Type == types.EventDone || e.Type == types.EventClosed || e.Type == types.EventCanceled {
			return
		}
		_, err := r.manager.Process(ctx, e.Event)
		if err != nil {
			return
		}
	}
}

type ResultAggregatorOption interface {
	Option(rg *ResultAggregator)
}

type ResultAggregatorOptionFunc func(rg *ResultAggregator)

func (fn ResultAggregatorOptionFunc) Option(rg *ResultAggregator) {
	fn(rg)
}

func WithMessage(message *types.Message) ResultAggregatorOption {
	return ResultAggregatorOptionFunc(func(rg *ResultAggregator) {
		rg.message = message
	})
}

func WithBatchSize(batch int) ResultAggregatorOptionFunc {
	return ResultAggregatorOptionFunc(func(rg *ResultAggregator) {
		rg.batchSize = batch
	})
}
