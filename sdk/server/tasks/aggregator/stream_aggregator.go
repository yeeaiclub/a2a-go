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

	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type StreamingConsumer struct {
	manager   *manager.TaskManager
	batchSize int
}

func NewStreamingAggregator(taskManager *manager.TaskManager, batchSize int) *StreamingConsumer {
	return &StreamingConsumer{
		manager:   taskManager,
		batchSize: batchSize,
	}
}

func (s *StreamingConsumer) Consume(ctx context.Context, queue *event.Queue) <-chan types.StreamEvent {
	out := make(chan types.StreamEvent, s.batchSize)

	go func() {
		defer close(out)
		s.processEvents(ctx, queue.Subscribe(ctx), out)
	}()

	return out
}

func (s *StreamingConsumer) processEvents(ctx context.Context, events <-chan types.StreamEvent, out chan<- types.StreamEvent) {
	for e := range events {
		result, shouldReturn := s.handleEvent(ctx, e)
		if out != nil && result != nil {
			select {
			case out <- *result:
			case <-ctx.Done():
				return
			}
		}
		if shouldReturn {
			return
		}
	}
}

func (s *StreamingConsumer) handleEvent(ctx context.Context, e types.StreamEvent) (*types.StreamEvent, bool) {
	switch e.Type {
	case types.EventClosed:
		return nil, true

	case types.EventData, types.EventDone:
		_, err := s.manager.Process(ctx, e.Event)
		if err != nil {
			errorEvent := types.StreamEvent{Type: types.EventError, Err: err}
			return &errorEvent, true
		}
		return &e, e.Type == types.EventDone

	case types.EventCanceled:
		cancelEvent := types.StreamEvent{Type: types.EventCanceled, Err: ctx.Err()}
		return &cancelEvent, true

	case types.EventError:
		return &e, true
	}

	return nil, false
}
