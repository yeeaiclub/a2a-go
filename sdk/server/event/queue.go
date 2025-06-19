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

func NewQueue(size int) *Queue {
	if size < 0 {
		return nil
	}
	if size == 0 {
		size = 1024
	}
	return &Queue{cap: size, queue: make(chan types.Event, size)}
}

func (q *Queue) DequeueWait(ctx context.Context) types.StreamEvent {
	q.Lock()
	defer q.Unlock()
	if q.closed || len(q.queue) == 0 {
		return types.StreamEvent{Err: errs.QueueClosed}
	}

	select {
	case <-ctx.Done():
		return types.StreamEvent{Err: ctx.Err()}
	case event := <-q.queue:
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
	case event := <-q.queue:
		return types.StreamEvent{Event: event}
	default:
		return types.StreamEvent{Err: errs.QueueEmpty}
	}
}

func (q *Queue) Enqueue(event types.Event) {
	q.queue <- event
	for _, ch := range q.children {
		ch.Enqueue(event)
	}
}

func (q *Queue) Close() {
	q.Lock()
	defer q.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	for _, ch := range q.children {
		ch.Close()
	}
}
