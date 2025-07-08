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

package aggregator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/internal/errs"
	"github.com/yeeaiclub/a2a-go/sdk/server/event"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/server/tasks/manager"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

func TestConsumeAll(t *testing.T) {
	testcases := []struct {
		name      string
		before    func(q *event.Queue, store *tasks.InMemoryTaskStore)
		contextId string
		want      *types.Task
	}{
		{
			name: "consumer all",
			before: func(q *event.Queue, store *tasks.InMemoryTaskStore) {
				task := &types.Task{Id: "1", ContextId: "2"}
				err := store.Save(context.Background(), task)
				require.NoError(t, err)

				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				q.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true})
			},
			want: &types.Task{Id: "1", ContextId: "2"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(10)
			defer queue.Close()
			store := tasks.NewInMemoryTaskStore()
			tc.before(queue, store)

			taskManger := manager.NewTaskManger(
				store,
				manager.WithTaskId("1"),
				manager.WithContextId("2"))

			aggregator := NewResultAggregator(taskManger)
			all, err := aggregator.ConsumeAll(context.Background(), queue)
			require.NoError(t, err)
			assert.Equal(t, tc.want, all)
		})
	}
}

func TestConsumeAndEmit(t *testing.T) {
	testcases := []struct {
		name   string
		before func(q *event.Queue, store *tasks.InMemoryTaskStore)
		want   []types.StreamEvent
	}{
		{
			name: "consumer and emit",
			before: func(q *event.Queue, store *tasks.InMemoryTaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1"})
				require.NoError(t, err)
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				q.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true})
			},
			want: []types.StreamEvent{
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false}, Type: types.EventData},
				{Event: &types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true}, Type: types.EventDone},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(2)
			defer queue.Close()
			store := tasks.NewInMemoryTaskStore()
			tc.before(queue, store)

			taskManger := manager.NewTaskManger(
				store,
				manager.WithTaskId("1"),
				manager.WithContextId("2"))
			aggregator := NewResultAggregator(taskManger)

			events := aggregator.ConsumeAndEmit(context.Background(), queue)

			var received []types.StreamEvent
			for ev := range events {
				require.NoError(t, ev.Err)
				received = append(received, ev)
			}

			require.ElementsMatch(t, tc.want, received)
		})
	}
}

func TestConsumeAndBreakOnInterrupt(t *testing.T) {
	testcases := []struct {
		name        string
		before      func(q *event.Queue, store *tasks.InMemoryTaskStore)
		contextId   string
		want        *types.Task
		expectError error
	}{
		{
			name: "consume",
			before: func(q *event.Queue, store *tasks.InMemoryTaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1"})
				require.NoError(t, err)
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: false})
				q.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "3", Final: true})
			},
			want:        &types.Task{Id: "1"},
			expectError: nil,
		},
		{
			name: "interrupt",
			before: func(q *event.Queue, store *tasks.InMemoryTaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1"})
				require.NoError(t, err)
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "1", Final: false})
				q.Enqueue(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: false, Status: types.TaskStatus{State: types.AUTH_REQUIRED}})
				q.EnqueueDone(&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "3", Final: true})
			},
			want:        &types.Task{Id: "1"},
			expectError: errs.ErrAuthRequired,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			queue := event.NewQueue(10)
			defer queue.Close()
			store := tasks.NewInMemoryTaskStore()
			tc.before(queue, store)

			taskManger := manager.NewTaskManger(
				store,
				manager.WithTaskId("1"),
				manager.WithContextId("2"))
			aggregator := NewResultAggregator(taskManger)
			events, err := aggregator.ConsumeAndBreakOnInterrupt(context.Background(), queue)
			if tc.expectError != nil {
				require.ErrorIs(t, err, tc.expectError)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, events)
		})
	}
}
