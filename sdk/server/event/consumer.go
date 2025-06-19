package event

import (
	"context"
	"time"

	"github.com/yumosx/a2a-go/sdk/types"
)

type Consumer struct {
	queue *Queue
}

func NewConsumer(queue *Queue) *Consumer {
	return &Consumer{queue: queue}
}

func (c *Consumer) ConsumeOne(ctx context.Context) types.StreamEvent {
	return c.queue.DequeueNoWait(ctx)
}

func (c *Consumer) ConsumeAll(ctx context.Context) <-chan types.StreamEvent {
	eventCh := make(chan types.StreamEvent, 1024)
	go func() {
		newCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				eventCh <- types.StreamEvent{Err: ctx.Err()}
				return
			default:
			}
			event := c.queue.DequeueWait(newCtx)
			if event.Err != nil {
				eventCh <- event
				return
			}
			if event.Done() {
				return
			}
			eventCh <- event
		}
	}()
	return eventCh
}
