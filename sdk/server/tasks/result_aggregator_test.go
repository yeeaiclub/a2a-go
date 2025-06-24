package tasks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/types"
)

func TestConsumeAll(t *testing.T) {
	testcases := []struct {
		name      string
		before    func(q *event.Queue, store *InMemoryTaskStore)
		contextId string
	}{
		{
			name: "consumer all",
			before: func(q *event.Queue, store *InMemoryTaskStore) {
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true})
				err := store.Save(context.Background(), &types.Task{Id: "1"})
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(10)
			defer queue.Close()
			store := NewInMemoryTaskStore()
			tc.before(queue, store)

			manager := NewTaskManger(store, WithTaskId("1"), WithContextId("2"))
			consumer := event.NewConsumer(queue, nil)
			aggregator := NewResultAggregator(manager, nil)
			all, err := aggregator.ConsumeAll(context.Background(), consumer)
			require.NoError(t, err)
			assert.Equal(t, all.GetTaskId(), "1")
		})
	}
}
