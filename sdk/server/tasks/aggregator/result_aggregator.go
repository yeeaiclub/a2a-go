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
	rg := &ResultAggregator{manager: taskManger}
	for _, opt := range options {
		opt.Option(rg)
	}
	return rg
}

// ConsumeAndEmit process the event stream from the consumer, updates the task store
func (r *ResultAggregator) ConsumeAndEmit(ctx context.Context, consumer *event.Consumer) <-chan types.StreamEvent {
	events := make(chan types.StreamEvent, r.batchSize)
	go func() {
		allEvents := consumer.ConsumeAll(ctx, r.batchSize)
		defer close(events)
		for {
			select {
			case <-ctx.Done():
				events <- types.StreamEvent{Err: ctx.Err()}
				return
			case e, ok := <-allEvents:
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
				ev, err := r.manager.Process(ctx, e.Event)
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
	events := consumer.ConsumeAll(ctx, r.batchSize)
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
				task, err := r.manager.GetTask(ctx)
				if err != nil {
					return nil, err
				}
				return task, nil
			}
			if ev.EventType() == "message" {
				msg := ev.Event.(*types.Message)
				r.message = msg
				return msg, nil
			}
			_, err := r.manager.Process(ctx, ev.Event)
			if err != nil {
				return nil, err
			}
		}
	}
}

// ConsumeAndBreakOnInterrupt process the event stream until completion or an interruptible state is encountered
func (r *ResultAggregator) ConsumeAndBreakOnInterrupt(ctx context.Context, consumer *event.Consumer) (types.Event, error) {
	events := consumer.ConsumeAll(ctx, r.batchSize)
	var result types.Event

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case e, ok := <-events:
			if !ok {
				task, err := r.manager.GetTask(ctx)
				result = task
				return result, err
			}
			if e.Err != nil {
				return nil, e.Err
			}
			if r.IsAuthRequired(e.Event) {
				r.continueConsume(ctx, events)
				return nil, errs.ErrAuthRequired
			}
			if e.Done() {
				task, err := r.manager.GetTask(ctx)
				result = task
				return result, err
			}
			_, err := r.manager.Process(ctx, e.Event)
			if err != nil {
				return nil, err
			}
		}
	}
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
				_, err := r.manager.Process(ctx, ev.Event)
				if err != nil {
					return
				}
			}
		}
	}()
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
