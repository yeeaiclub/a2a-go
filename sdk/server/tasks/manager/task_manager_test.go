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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mocktasks "github.com/yeeaiclub/a2a-go/internal/mocks/tasks"
	"github.com/yeeaiclub/a2a-go/sdk/types"
	"go.uber.org/mock/gomock"
)

func TestGetTask(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func(store *mocktasks.MockTaskStore)
		taskId      string
		contextId   string
		initMessage types.Message
		want        *types.Task
		expectErr   bool
	}{
		{
			name:        "get task success",
			taskId:      "1",
			contextId:   "2",
			initMessage: types.Message{Role: "user", TaskID: "1"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(&types.Task{Id: "1"}, nil)
			},
			want: &types.Task{Id: "1"},
		},
		{
			name:      "get task not found",
			taskId:    "not-exist",
			contextId: "2",
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "not-exist").Return(nil, nil)
			},
			want: nil,
		},
		{
			name:      "get task error",
			taskId:    "err-id",
			contextId: "2",
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "err-id").Return(nil, errors.New("db error"))
			},
			want:      nil,
			expectErr: true,
		},
		{
			name:      "empty taskId should error",
			taskId:    "",
			contextId: "2",
			setup:     func(store *mocktasks.MockTaskStore) {},
			want:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocktasks.NewMockTaskStore(ctrl)
			tc.setup(store)

			manager := NewTaskManager(
				store,
				WithTaskId(tc.taskId),
				WithContextId(tc.contextId),
				WithInitMessage(&tc.initMessage),
			)

			task, err := manager.GetTask(context.Background())
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, task)
			}
		})
	}
}

func TestTaskManager_saveTask(t *testing.T) {
	testCases := []struct {
		name      string
		taskId    string
		contextId string
		task      *types.Task
		setup     func(store *mocktasks.MockTaskStore)
		check     func(manager *TaskManager)
		wantErr   bool
	}{
		{
			name:      "save task success",
			taskId:    "1",
			contextId: "2",
			task:      &types.Task{Id: "1", ContextId: "2"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"}).Return(nil)
			},
			check: func(manager *TaskManager) {
				assert.Equal(t, "1", manager.taskId)
				assert.Equal(t, "2", manager.contextId)
			},
		},
		{
			name:   "save task error",
			taskId: "1",
			task:   &types.Task{Id: "1"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1"}).Return(errors.New("save error"))
			},
			check:   func(manager *TaskManager) {},
			wantErr: true,
		},
		{
			name: "create manager without id and save",
			task: &types.Task{Id: "1", ContextId: "2"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"}).Return(nil)
			},
			check: func(manager *TaskManager) {
				assert.Equal(t, "1", manager.taskId)
				assert.Equal(t, "2", manager.contextId)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocktasks.NewMockTaskStore(ctrl)
			tc.setup(store)
			manager := NewTaskManager(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))
			err := manager.saveTask(context.Background(), tc.task)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				if tc.check != nil {
					tc.check(manager)
				}
			}
		})
	}
}

func TestTaskManager_EnsureTask(t *testing.T) {
	testCases := []struct {
		name        string
		taskId      string
		contextId   string
		initMessage types.Message
		event       types.Event
		want        *types.Task
		setup       func(store *mocktasks.MockTaskStore)
		wantErr     bool
	}{
		{
			name:      "get task from store",
			taskId:    "1",
			contextId: "2",
			event:     &types.TaskStatusUpdateEvent{TaskId: "1"},
			want:      &types.Task{Id: "1"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(&types.Task{Id: "1"}, nil)
			},
		},
		{
			name:        "init new task and save",
			taskId:      "1",
			contextId:   "1",
			initMessage: types.Message{Role: types.User, TaskID: "1", ContextID: "1", Kind: "init", MessageID: "msg-1"},
			event:       &types.TaskStatusUpdateEvent{TaskId: "1"},
			want: &types.Task{
				Id:   "1",
				Kind: types.EventTypeTask,
				Status: types.TaskStatus{
					State: types.SUBMITTED,
				},
				History: []*types.Message{
					{Role: types.User, TaskID: "1", ContextID: "1", Kind: "init", MessageID: "msg-1"},
				},
			},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(nil, nil)
				store.EXPECT().Save(gomock.Any(), gomock.Any()).DoAndReturn(
					func(_ context.Context, task *types.Task) error {
						require.Equal(t, "1", task.Id)
						require.Equal(t, types.SUBMITTED, task.Status.State)
						return nil
					},
				)
			},
		},
		{
			name:      "get task error",
			taskId:    "1",
			contextId: "2",
			event:     &types.TaskStatusUpdateEvent{TaskId: "1"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(nil, errors.New("db error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "save new task error",
			taskId:    "1",
			contextId: "2",
			event:     &types.TaskStatusUpdateEvent{TaskId: "1"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Get(gomock.Any(), "1").Return(nil, nil)
				store.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("save error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocktasks.NewMockTaskStore(ctrl)
			if tc.setup != nil {
				tc.setup(store)
			}
			manager := NewTaskManager(store, WithTaskId(tc.taskId), WithContextId(tc.contextId), WithInitMessage(&tc.initMessage))
			task, err := manager.EnsureTask(context.Background(), tc.event)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, task)
			}
		})
	}
}

func TestTaskManager_UpdateWithMessage(t *testing.T) {
	testCases := []struct {
		name      string
		taskId    string
		contextId string
		message   types.Message
		task      types.Task
	}{
		{
			name:      "update with message and status message",
			taskId:    "1",
			contextId: "1",
			message:   types.Message{Role: types.User},
			task:      types.Task{Id: "1", Status: types.TaskStatus{Message: &types.Message{Role: types.Agent}}},
		},
		{
			name:      "update with message only",
			taskId:    "2",
			contextId: "2",
			message:   types.Message{Role: types.User},
			task:      types.Task{Id: "2"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocktasks.NewMockTaskStore(ctrl)
			manager := NewTaskManager(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))

			task := tc.task
			if tc.task.Status.Message != nil {
				result := manager.PushMessageToHistory(&tc.message, &task)
				if assert.Len(t, result.History, 2) {
					assert.Equal(t, tc.task.Status.Message.Role, result.History[0].Role)
					assert.Equal(t, tc.message.Role, result.History[1].Role)
				}
			} else {
				result := manager.PushMessageToHistory(&tc.message, &task)
				if assert.Len(t, result.History, 1) {
					assert.Equal(t, tc.message.Role, result.History[0].Role)
				}
			}
		})
	}
}

func TestSaveTaskEvent(t *testing.T) {
	testCases := []struct {
		name      string
		taskId    string
		contextId string
		event     types.Event
		setup     func(store *mocktasks.MockTaskStore)
		want      *types.Task
		wantErr   bool
	}{
		{
			name:      "save task event with task",
			taskId:    "1",
			contextId: "2",
			event:     &types.Task{Id: "1", ContextId: "2"},
			setup: func(store *mocktasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"}).Return(nil)
			},
			want: &types.Task{Id: "1", ContextId: "2"},
		},
		{
			name:      "save task event with mismatched id",
			taskId:    "1",
			contextId: "2",
			event:     &types.Task{Id: "2"},
			setup:     func(store *mocktasks.MockTaskStore) {},
			want:      nil,
			wantErr:   true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mocktasks.NewMockTaskStore(ctrl)
			if tc.setup != nil {
				tc.setup(store)
			}
			manager := NewTaskManager(store, WithTaskId(tc.taskId), WithContextId(tc.contextId))
			result, err := manager.SaveTaskEvent(context.Background(), tc.event)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, result)
			}
		})
	}
}
