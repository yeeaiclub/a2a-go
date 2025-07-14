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
	"sync/atomic"

	"github.com/yeeaiclub/a2a-go/sdk/types"
)

// Queue is a thread-safe event queue for streaming task events.
// It supports enqueueing data, done, and error events, and allows
// multiple consumers to subscribe to the event stream.
type Queue struct {
	closed atomic.Bool            // Indicates if the queue is closed
	ch     chan types.StreamEvent // Underlying channel for events
	size   uint                   // Buffer size
}

// NewQueue creates a new Queue with the given buffer size.
func NewQueue(size uint) *Queue {
	return &Queue{ch: make(chan types.StreamEvent, size), size: size}
}

// Enqueue adds a data event to the queue.
// Returns false if the queue is closed or full.
func (q *Queue) Enqueue(data types.Event) bool {
	if q.closed.Load() {
		return false
	}
	if data == nil {
		return false
	}

	if data.Done() {
		return q.EnqueueDone(data)
	}
	return q.EnqueueEvent(data)
}

// EnqueueDone adds a done event to the queue.
// Returns false if the queue is closed or full.
func (q *Queue) EnqueueDone(data types.Event) bool {
	if q.closed.Load() {
		return false
	}
	select {
	case q.ch <- types.StreamEvent{Type: types.EventDone, Event: data}:
		return true
	default:
		return false
	}
}

// EnqueueEvent adds an event to the queue
func (q *Queue) EnqueueEvent(data types.Event) bool {
	if q.closed.Load() {
		return false
	}
	select {
	case q.ch <- types.StreamEvent{Type: types.EventData, Event: data}:
		return true
	default:
		return false
	}
}

// EnqueueError adds an error event to the queue.
// Returns false if the queue is closed or full.
func (q *Queue) EnqueueError(err error) bool {
	if q.closed.Load() {
		return false
	}
	select {
	case q.ch <- types.StreamEvent{Type: types.EventError, Err: err}:
		return true
	default:
		return false
	}
}

// Subscribe returns a channel to receive events from the queue.
// The returned channel will be closed when the queue is closed,
// or when a done/error event is received, or when the context is canceled.
func (q *Queue) Subscribe(ctx context.Context) <-chan types.StreamEvent {
	out := make(chan types.StreamEvent, q.size)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				// Send a canceled event and exit
				out <- types.StreamEvent{Type: types.EventCanceled, Err: ctx.Err()}
				return
			case e, ok := <-q.ch:
				if !ok {
					// Queue closed, send closed event and exit
					out <- types.StreamEvent{Type: types.EventClosed}
					return
				}
				out <- e
				// Stop streaming on done or error event
				if e.Type == types.EventDone || e.Type == types.EventError {
					return
				}
			}
		}
	}()
	return out
}

// Close closes the queue and releases all resources.
// Further enqueue operations will fail.
func (q *Queue) Close() {
	if q.closed.CompareAndSwap(false, true) {
		close(q.ch)
	}
}
