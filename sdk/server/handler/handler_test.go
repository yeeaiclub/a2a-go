// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/server/execution"
	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/server/tasks/updater"
	"github.com/yumosx/a2a-go/sdk/types"
)

type Executor struct{}

func NewExecutor() *Executor {
	return &Executor{}
}

func (e *Executor) Execute(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	u.Complete(updater.WithFinal(true))
	return nil
}

func (e *Executor) Cancel(ctx context.Context, requestContext *execution.RequestContext, queue *event.Queue) error {
	u := updater.NewTaskUpdater(queue, requestContext.TaskId, requestContext.ContextId)
	u.Complete(updater.WithFinal(true))
	return nil
}

type QueueManger struct{}

func (q QueueManger) Add(ctx context.Context, taskId string, queue *event.Queue) error {
	panic("implement me")
}

func (q QueueManger) Get(ctx context.Context, taskId string) (*event.Queue, error) {
	panic("implement me")
}

func (q QueueManger) Tap(ctx context.Context, taskId string) (*event.Queue, error) {
	panic("implement me")
}

func (q QueueManger) Close(ctx context.Context, taskId string) error {
	panic("implement me")
}

func (q QueueManger) CreateOrTap(ctx context.Context, taskId string) (*event.Queue, error) {
	return event.NewQueue(10), nil
}

func TestGatTask(t *testing.T) {
	testcases := []struct {
		name   string
		input  types.TaskQueryParams
		before func(store tasks.TaskStore)
		want   *types.Task
	}{
		{
			name:  "on get card",
			input: types.TaskQueryParams{Id: "1"},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: &types.Task{Id: "1", ContextId: "2"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := NewExecutor()
			handler := NewDefaultHandler(store, executor)
			task, err := handler.OnGetTask(context.Background(), tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.want, task)
		})
	}
}

func TestOnMessageSend(t *testing.T) {
	testcases := []struct {
		name   string
		input  types.MessageSendParam
		before func(store tasks.TaskStore)
		want   types.Event
	}{
		{
			name:  "on message send",
			input: types.MessageSendParam{Message: &types.Message{TaskID: "1", ContextID: "2"}},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: &types.Task{Id: "1", ContextId: "2", History: []*types.Message{{TaskID: "1", ContextID: "2"}}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := NewExecutor()
			manger := QueueManger{}
			handler := NewDefaultHandler(store, executor, WithQueueManger(manger))
			ev, err := handler.OnMessageSend(context.Background(), tc.input)
			require.NoError(t, err)
			assert.Equal(t, ev, tc.want)
		})
	}
}

func TestOnMessageSendStream(t *testing.T) {
	testcases := []struct {
		name   string
		input  types.MessageSendParam
		before func(store tasks.TaskStore)
		after  func(events []types.Event)
		want   []types.Event
	}{
		{
			name:  "on message stream",
			input: types.MessageSendParam{Message: &types.Message{TaskID: "1", ContextID: "2"}},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: []types.Event{&types.TaskStatusUpdateEvent{TaskId: "1", ContextId: "2", Final: true, Status: types.TaskStatus{State: types.COMPLETED}}},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := NewExecutor()
			manager := QueueManger{}
			handler := NewDefaultHandler(store, executor, WithQueueManger(manager))
			events := handler.OnMessageSendStream(context.Background(), tc.input)
			var received []types.Event
			for ev := range events {
				require.NoError(t, ev.Err)
				received = append(received, ev.Event)
			}

			for i := range received {
				if ev, ok := received[i].(*types.TaskStatusUpdateEvent); ok {
					ev.Status.TimeStamp = ""
				}
			}

			assert.ElementsMatch(t, tc.want, received)
		})
	}
}

func TestOnCancelTask(t *testing.T) {
	testcases := []struct {
		name   string
		input  types.TaskIdParams
		before func(store tasks.TaskStore)
		want   *types.Task
	}{
		{
			name:  "cancel task",
			input: types.TaskIdParams{Id: "1"},
			before: func(store tasks.TaskStore) {
				err := store.Save(context.Background(), &types.Task{Id: "1", ContextId: "2"})
				require.NoError(t, err)
			},
			want: &types.Task{Id: "1", ContextId: "2"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			store := tasks.NewInMemoryTaskStore()
			tc.before(store)
			executor := NewExecutor()
			manger := QueueManger{}
			handler := NewDefaultHandler(store, executor, WithQueueManger(manger))
			task, err := handler.OnCancelTask(context.Background(), tc.input)
			require.NoError(t, err)
			assert.Equal(t, task, tc.want)
		})
	}
}
