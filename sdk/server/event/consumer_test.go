package event

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/types"
)

func TestConsumeAll(t *testing.T) {
	testcases := []struct {
		name   string
		before func(q *Queue)
		want   []types.StreamEvent
	}{
		{
			name: "consume single task",
			before: func(q *Queue) {
				q.Enqueue(&types.Task{Id: "1", Status: types.TaskStatus{State: types.FAILED}})
			},
			want: []types.StreamEvent{
				{Event: &types.Task{Id: "1", Status: types.TaskStatus{State: types.FAILED}}},
			},
		},
		{
			name: "consume multiple tasks",
			before: func(q *Queue) {
				q.Enqueue(&types.Task{Id: "1"})
				q.Enqueue(&types.Task{Id: "2", Status: types.TaskStatus{State: types.FAILED}})
			},
			want: []types.StreamEvent{
				{Event: &types.Task{Id: "1"}},
				{Event: &types.Task{Id: "2", Status: types.TaskStatus{State: types.FAILED}}},
			},
		},
		{
			name: "consume multiple task update event",
			before: func(q *Queue) {
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1"})
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", Final: true})
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1"}},
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", Final: true}},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := NewQueue(2)
			defer queue.Close()
			tc.before(queue)

			consumer := NewConsumer(queue, nil)
			events := consumer.ConsumeAll(context.Background())

			var received []types.StreamEvent
			for event := range events {
				require.NoError(t, event.Err)
				received = append(received, event)
			}
			require.ElementsMatch(t, tc.want, received)
		})
	}
}
