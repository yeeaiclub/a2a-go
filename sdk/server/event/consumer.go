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

package event

import (
	"context"

	"github.com/yeeaiclub/a2a-go/sdk/types"
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
		errChanClosed := false
		for {
			select {
			case <-ctx.Done():
				eventCh <- types.StreamEvent{Err: ctx.Err()}
				return
			default:
				if !errChanClosed {
					select {
					case err, ok := <-c.err:
						if !ok {
							errChanClosed = true
							continue
						}
						if err != nil {
							eventCh <- types.StreamEvent{Err: err}
							return
						}
					default:
					}
				}
				event := c.queue.DequeueWait(ctx)
				if c.sendAndCheckDone(eventCh, event) {
					return
				}
			}
		}
	}()
	return eventCh
}

func (c *Consumer) sendAndCheckDone(ch chan types.StreamEvent, s types.StreamEvent) bool {
	if s.Err != nil {
		ch <- types.StreamEvent{Err: s.Err}
		return true
	}

	if s.Event != nil && s.Done() {
		ch <- types.StreamEvent{Event: s.Event}
		return true
	}

	if s.Event != nil {
		ch <- types.StreamEvent{Event: s.Event}
		return false
	}
	return false
}
