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

type InterruptibleConsumer struct {
	manager *manager.TaskManager
}

func NewInterruptibleConsumer(manager *manager.TaskManager) *InterruptibleConsumer {
	return &InterruptibleConsumer{manager: manager}
}

func (r *InterruptibleConsumer) Consume(ctx context.Context, queue *event.Queue) (types.Event, error) {
	for e := range queue.Subscribe(ctx) {
		switch e.Type {
		case types.EventCanceled:
			return nil, ctx.Err()
		case types.EventError:
			return nil, e.Err
		case types.EventClosed:
			task, err := r.manager.GetTask(ctx)
			if err != nil {
				return nil, err
			}
			return task, nil
		case types.EventDone:
			_, err := r.manager.Process(ctx, e.Event)
			if err != nil {
				return nil, err
			}
			task, err := r.manager.GetTask(ctx)
			if err != nil {
				return nil, err
			}
			return task, nil
		case types.EventData:
			if msg, ok := e.Event.(*types.Message); ok {
				return msg, nil
			}

			// Check for auth_required before processing the event
			if r.IsAuthRequired(e.Event) {
				// Start background processing
				go r.continueConsume(ctx, queue.Subscribe(ctx))
				return nil, errs.ErrAuthRequired
			}

			_, err := r.manager.Process(ctx, e.Event)
			if err != nil {
				return nil, err
			}
		}
	}

	// If we reach here, get the final task
	task, err := r.manager.GetTask(ctx)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *InterruptibleConsumer) IsAuthRequired(event types.Event) bool {
	if event.GetKind() == types.EventTypeStatusUpdate {
		updateEvent := event.(*types.TaskStatusUpdateEvent)
		if updateEvent.Status.State == types.AuthRequired {
			log.Debug("Encountered an auth-required task: breaking synchronous message/send flow.")
			return true
		}
	}

	if event.GetKind() == types.EventTypeTask {
		taskEvent := event.(*types.Task)
		if taskEvent.Status.State == types.AuthRequired {
			log.Debug("Encountered an auth-required task: breaking synchronous message/send flow.")
			return true
		}
	}
	return false
}

// continueConsume processing an event stream in backhand task
func (r *InterruptibleConsumer) continueConsume(ctx context.Context, events <-chan types.StreamEvent) {
	for e := range events {
		switch e.Type {
		case types.EventClosed, types.EventCanceled, types.EventError:
			return
		case types.EventDone, types.EventData:
			if _, err := r.manager.Process(ctx, e.Event); err != nil {
				return
			}
			if e.Type == types.EventDone {
				return
			}
		}
	}
}
