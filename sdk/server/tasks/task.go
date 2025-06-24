package tasks

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

// TaskStore Agent task store interface
type TaskStore interface {
	// Save or updates a task in the store.
	Save(ctx context.Context, task *types.Task) error
	// Get retrieves a task from the store by ID.
	Get(ctx context.Context, taskId string) (*types.Task, error)
	// Delete a tasks from the store by id
	Delete(ctx context.Context, id string) error
}
