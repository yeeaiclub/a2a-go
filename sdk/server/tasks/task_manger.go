package tasks

import (
	"context"

	"github.com/yumosx/a2a-go/sdk/types"
)

type TaskManger struct {
	taskId      string
	contextId   string
	store       TaskStore
	initMessage types.Message
	currentTask *types.Task
}

func NewTaskManger(taskId string, contextId string, store TaskStore, initMessage types.Message) *TaskManger {
	return &TaskManger{taskId: taskId, contextId: contextId, store: store, initMessage: initMessage}
}

func (t *TaskManger) GetTask(ctx context.Context) (*types.Task, error) {
	if t.currentTask != nil {
		return t.currentTask, nil
	}
	task, err := t.store.Get(ctx, t.taskId)
	if err != nil {
		return &task, err
	}
	t.currentTask = &task
	return &task, nil
}

func (t *TaskManger) UpdateMessage(ctx context.Context, message types.Message, task *types.Task) *types.Task {
	if len(task.History) != 0 {
		task.History = append(task.History, task.Status.Message)
		task.History = append(task.History, message)
	} else {
		task.History = []types.Message{task.Status.Message, message}
	}
	task.Status.Message = types.Message{}
	t.currentTask = task
	return task
}

func (t *TaskManger) saveTask(ctx context.Context, task types.Task) error {
	return t.store.Save(ctx, task)
}

func (t *TaskManger) updateWithMessage(ctx context.Context, message types.Message, task types.Task) types.Task {
	return task
}
