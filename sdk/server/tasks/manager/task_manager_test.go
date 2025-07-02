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

package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeeaiclub/a2a-go/sdk/types"
	"github.com/yeeaiclub/a2a-go/test/mocks/tasks"
	"go.uber.org/mock/gomock"
)

func TestGetTask(t *testing.T) {
	testcases := []struct {
		name        string
		before      func(store *tasks.MockTaskStore)
		taskId      string
		contextId   string
		initMessage types.Message
		want        *types.Task
	}{
		{
			name:        "get task",
			taskId:      "1",
			contextId:   "2",
			initMessage: types.Message{Role: "user", TaskID: "1"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(&types.Task{Id: "1"}, nil)
			},
			want: &types.Task{Id: "1"},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := tasks.NewMockTaskStore(ctrl)
			tc.before(store)

			manger := NewTaskManger(
				store,
				WithTaskId(tc.taskId),
				WithContextId(tc.contextId),
				WithInitMessage(&tc.initMessage))

			task, err := manger.GetTask(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.want, task)
		})
	}
}

func TestSaveTask(t *testing.T) {
	testcases := []struct {
		name      string
		taskId    string
		contextId string
		task      *types.Task
		before    func(store *tasks.MockTaskStore)
		after     func(manager *TaskManager)
	}{
		{
			name:      "save task",
			taskId:    "1",
			contextId: "2",
			task:      &types.Task{Id: "1", ContextId: "2"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"})
			},
			after: func(manager *TaskManager) {
				assert.Equal(t, "1", manager.taskId)
				assert.Equal(t, "2", manager.contextId)
			},
		},
		{
			name: "create manger without id and save",
			task: &types.Task{Id: "1", ContextId: "2"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"})
			},
			after: func(manager *TaskManager) {
				assert.Equal(t, "1", manager.taskId)
				assert.Equal(t, "2", manager.contextId)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := tasks.NewMockTaskStore(ctrl)
			tc.before(store)
			manager := NewTaskManger(store,
				WithTaskId(tc.taskId),
				WithContextId(tc.contextId))
			err := manager.saveTask(context.Background(), tc.task)
			require.NoError(t, err)
			tc.after(manager)
		})
	}
}

func TestEnsureTask(t *testing.T) {
	testcases := []struct {
		name        string
		taskId      string
		contextId   string
		initMessage types.Message
		event       types.Event
		want        *types.Task
		before      func(store *tasks.MockTaskStore)
	}{
		{
			name:      "get task",
			taskId:    "1",
			contextId: "2",
			event:     &types.TaskStatusUpdateEvent{TaskId: "1"},
			want:      &types.Task{Id: "1"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(&types.Task{Id: "1"}, nil)
			},
		},
		{
			name:      "init task and save",
			taskId:    "1",
			contextId: "1",
			event:     &types.TaskStatusUpdateEvent{TaskId: "1"},
			want: &types.Task{Id: "1", Status: types.TaskStatus{
				State: types.SUBMITTED,
			}},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(nil, nil)
				store.EXPECT().Save(gomock.Any(), &types.Task{
					Id:     "1",
					Status: types.TaskStatus{State: types.SUBMITTED},
				}).Return(nil)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := tasks.NewMockTaskStore(ctrl)
			tc.before(store)

			manager := NewTaskManger(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))
			task, err := manager.EnsureTask(context.Background(), tc.event)
			require.NoError(t, err)
			assert.Equal(t, tc.want, task)
		})
	}
}

func TestUpdateWithMessage(t *testing.T) {
	testcase := []struct {
		name      string
		taskId    string
		contextId string
		message   types.Message
		task      types.Task
	}{
		{
			name:      "update with message",
			taskId:    "1",
			contextId: "1",
			message:   types.Message{Role: types.User},
			task:      types.Task{Id: "1", Status: types.TaskStatus{Message: &types.Message{Role: types.Agent}}},
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := tasks.NewMockTaskStore(ctrl)
			manager := NewTaskManger(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))
			manager.UpdateWithMessage(&tc.message, &tc.task)
			assert.ElementsMatch(t, tc.task.History, []*types.Message{{Role: types.User}, {Role: types.Agent}})
		})
	}
}

func TestSaveStream(t *testing.T) {
	testcases := []struct {
		name      string
		taskId    string
		contextId string
		event     types.Event
		before    func(store *tasks.MockTaskStore)
		want      *types.Task
	}{
		{
			name:      "save task",
			taskId:    "1",
			contextId: "2",
			event:     &types.Task{Id: "1"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1"}).Return(nil)
			},
			want: &types.Task{Id: "1"},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := tasks.NewMockTaskStore(ctrl)
			tc.before(store)
			manger := NewTaskManger(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))
			event, err := manger.SaveTaskEvent(context.Background(), tc.event)
			require.NoError(t, err)
			assert.Equal(t, tc.want, event)
		})
	}
}
