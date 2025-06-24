package manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yumosx/a2a-go/sdk/types"
	"github.com/yumosx/a2a-go/test/mocks/tasks"
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
			assert.Equal(t, task, tc.want)
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
				assert.Equal(t, manager.taskId, "1")
				assert.Equal(t, manager.contextId, "2")
			},
		},
		{
			name: "create manger without id and save",
			task: &types.Task{Id: "1", ContextId: "2"},
			before: func(store *tasks.MockTaskStore) {
				store.EXPECT().Save(gomock.Any(), &types.Task{Id: "1", ContextId: "2"})
			},
			after: func(manager *TaskManager) {
				assert.Equal(t, manager.taskId, "1")
				assert.Equal(t, manager.contextId, "2")
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
			assert.Equal(t, event, tc.want)
		})
	}
}
