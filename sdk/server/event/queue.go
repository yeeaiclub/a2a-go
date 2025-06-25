// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package event

import (
	"context"
	"sync"

	"github.com/yumosx/a2a-go/internal/errs"
	"github.com/yumosx/a2a-go/sdk/types"
)

type Queue struct {
	sync.Mutex
	cap      int
	queue    chan types.Event
	closed   bool
	children []*Queue
}

const DefaultMaxQueueSize = 1024

func NewQueue(size int) *Queue {
	if size < 0 {
		return nil
	}
	if size == 0 {
		size = DefaultMaxQueueSize
	}
	return &Queue{cap: size, queue: make(chan types.Event, size)}
}

func (q *Queue) DequeueWait(ctx context.Context) types.StreamEvent {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return types.StreamEvent{Err: errs.QueueClosed}
	}

	select {
	case <-ctx.Done():
		return types.StreamEvent{Err: ctx.Err()}
	case event, ok := <-q.queue:
		if !ok {
			return types.StreamEvent{Closed: true}
		}
		return types.StreamEvent{Event: event}
	}
}

func (q *Queue) DequeueNoWait(ctx context.Context) types.StreamEvent {
	q.Lock()
	defer q.Unlock()

	if q.closed || len(q.queue) == 0 {
		return types.StreamEvent{Err: errs.QueueClosed}
	}

	select {
	case <-ctx.Done():
		return types.StreamEvent{Err: ctx.Err()}
	case event, ok := <-q.queue:
		if !ok {
			return types.StreamEvent{Closed: true}
		}
		return types.StreamEvent{Event: event}
	default:
		return types.StreamEvent{Err: errs.QueueEmpty}
	}
}

func (q *Queue) Enqueue(event types.Event) {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return
	}
	select {
	case q.queue <- event:
		for _, ch := range q.children {
			ch.Enqueue(event)
		}
	default:
		return
	}
}

func (q *Queue) Close() {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	close(q.queue)
	for _, ch := range q.children {
		ch.Close()
	}
}
