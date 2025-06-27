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

package event

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

// Consumer to read events form the agent event queue
type Consumer struct {
	queue *Queue
	err   <-chan error
}

func NewConsumer(queue *Queue, errCh <-chan error) *Consumer {
	return &Consumer{queue: queue, err: errCh}
}

// ConsumeOne consume one event from the agent queue non-blocking
func (c *Consumer) ConsumeOne(ctx context.Context) types.StreamEvent {
	return c.queue.DequeueNoWait(ctx)
}

// ConsumeAll consume all the agents streaming events form agent
func (c *Consumer) ConsumeAll(ctx context.Context, size int) <-chan types.StreamEvent {
	eventCh := make(chan types.StreamEvent, size)
	go func() {
		defer close(eventCh)
		for {
			select {
			case <-ctx.Done():
				eventCh <- types.StreamEvent{Err: ctx.Err()}
				return
			case err, ok := <-c.err:
				if err != nil && !ok {
					eventCh <- types.StreamEvent{Err: err}
					return
				}
			default:
			}
			event := c.queue.DequeueWait(ctx)
			if event.Err != nil {
				eventCh <- event
				return
			}
			if event.Event != nil && event.Done() {
				eventCh <- event
				return
			}
			if event.Event != nil {
				eventCh <- event
			}
		}
	}()
	return eventCh
}
