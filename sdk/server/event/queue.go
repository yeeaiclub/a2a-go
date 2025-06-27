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
	"sync"
	"sync/atomic"

	"github.com/yumosx/a2a-go/internal/errs"
	"github.com/yumosx/a2a-go/sdk/types"
)

type Queue struct {
	sync.Mutex
	cap      int
	queue    chan types.Event
	closed   atomic.Bool
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
	select {
	case <-ctx.Done():
		return types.StreamEvent{Err: ctx.Err()}
	case event, ok := <-q.queue:
		if !ok {
			return types.StreamEvent{Closed: true}
		}
		return types.StreamEvent{Event: event}
	default:
		return types.StreamEvent{Err: errs.ErrQueueEmpty}
	}
}

func (q *Queue) Enqueue(event types.Event) bool {
	select {
	case q.queue <- event:
		for _, ch := range q.children {
			ch.Enqueue(event)
		}
		return true
	default:
		return false
	}
}

func (q *Queue) Close() {
	q.closed.Store(true)
	close(q.queue)
	for _, ch := range q.children {
		ch.Close()
	}
}
