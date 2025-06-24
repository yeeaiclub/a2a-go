package event

import (
	"context"
	"time"

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
func (c *Consumer) ConsumeAll(ctx context.Context) <-chan types.StreamEvent {
	eventCh := make(chan types.StreamEvent, 10)
	go func() {
		defer close(eventCh)
		newCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
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
			event := c.queue.DequeueWait(newCtx)
			if event.Err != nil {
				eventCh <- event
				return
			}
			if event.Done() {
				eventCh <- event
				return
			}
			eventCh <- event
		}
	}()
	return eventCh
}
