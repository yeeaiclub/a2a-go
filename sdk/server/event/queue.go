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

type Queue struct {
	closed atomic.Bool
	ch     chan types.StreamEvent
	size   uint
}

func NewQueue(size uint) *Queue {
	return &Queue{ch: make(chan types.StreamEvent, size), size: size}
}

func (q *Queue) Enqueue(data types.Event) bool {
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

func (q *Queue) Subscribe(ctx context.Context) <-chan types.StreamEvent {
	out := make(chan types.StreamEvent, q.size)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				out <- types.StreamEvent{Type: types.EventCanceled, Err: ctx.Err()}
				return
			case e, ok := <-q.ch:
				if !ok {
					out <- types.StreamEvent{Type: types.EventClosed}
					return
				}
				out <- e
				if e.Type == types.EventDone {
					return
				}
			}
		}
	}()
	return out
}

func (q *Queue) Close() {
	if q.closed.CompareAndSwap(false, true) {
		close(q.ch)
	}
}
