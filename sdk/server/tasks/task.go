package tasks

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

type TaskStore interface {
	Save(ctx context.Context, task types.Task) error
	Get(ctx context.Context, id string) (types.Task, error)
	Delete(ctx context.Context, id string) error
}
