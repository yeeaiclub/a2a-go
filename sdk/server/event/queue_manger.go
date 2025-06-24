package event

import "context"

// QueueManager manger for manging the event queue lifecycles per task
type QueueManager interface {
	// Add a new event queue associated with a task ID.
	Add(ctx context.Context, taskId string, queue *Queue) error
	// Get Retrieves the event queue for a task ID.
	Get(ctx context.Context, taskId string) (*Queue, error)
	// Tap Creates a child event queue (tap) for an existing task ID.
	Tap(ctx context.Context, taskId string) (*Queue, error)
	// Close and remove the event queue for a task ID.
	Close(ctx context.Context, taskId string) error
	// CreateOrTap Creates a queue if one doesn't exist, otherwise taps the existing one.
	CreateOrTap(ctx context.Context, taskId string) (*Queue, error)
}
