// Copyright 2025 yumosx
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
	"fmt"

	"github.com/yumosx/a2a-go/sdk/server/tasks"
	"github.com/yumosx/a2a-go/sdk/types"
)

// TaskManager helps mange a task's lifecycle during execution of a request
type TaskManager struct {
	taskId      string
	contextId   string
	store       tasks.TaskStore
	initMessage *types.Message
	currentTask *types.Task
}

type TaskManagerOption interface {
	Option(manager *TaskManager)
}

type TaskManagerOptionFunc func(manger *TaskManager)

func (fn TaskManagerOptionFunc) Option(manger *TaskManager) {
	fn(manger)
}

func WithTaskId(taskId string) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.taskId = taskId
	})
}

func WithContextId(contextId string) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.contextId = contextId
	})
}

func WithInitMessage(message *types.Message) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.initMessage = message
	})
}

func NewTaskManger(store tasks.TaskStore, opts ...TaskManagerOption) *TaskManager {
	manger := &TaskManager{store: store}

	for _, opt := range opts {
		opt.Option(manger)
	}
	return manger
}

// GetTask retrieves the current task object, either from memory or the store
func (t *TaskManager) GetTask(ctx context.Context) (*types.Task, error) {
	if t.taskId == "" {
		return nil, errors.New("task_id is not set, cannot get task")
	}

	if t.currentTask != nil {
		return t.currentTask, nil
	}

	task, err := t.store.Get(ctx, t.taskId)
	if err != nil {
		return task, err
	}
	t.currentTask = task
	return task, nil
}

// SaveTaskEvent Process a tasked-related event
func (t *TaskManager) SaveTaskEvent(ctx context.Context, event types.Event) (*types.Task, error) {
	taskId := event.GetTaskId()

	if t.taskId != taskId {
		return nil, fmt.Errorf("task in event doesn't match TaskManager %s %s", t.taskId, taskId)
	}

	if t.taskId == "" {
		t.taskId = taskId
	}

	if t.contextId != "" && t.contextId != event.GetContextId() {
		t.contextId = event.GetContextId()
	}

	if event.EventType() == "task" {
		return t.handleTaskEvent(ctx, event)
	}
	return t.handleEvent(ctx, event)
}

func (t *TaskManager) handleTaskEvent(ctx context.Context, event types.Event) (*types.Task, error) {
	task, ok := event.(*types.Task)
	if !ok {
		return nil, errors.New("invalid event type for task event")
	}
	if err := t.saveTask(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TaskManager) handleEvent(ctx context.Context, event types.Event) (*types.Task, error) {
	task, err := t.EnsureTask(ctx, event)
	if err != nil {
		return nil, err
	}

	if event.EventType() == "status_update" {
		if task.Status.Message != nil {
			task.History = append(task.History, task.Status.Message)
		}
	}
	err = t.saveTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateWithMessage updates a task object in memory by adding a new initMessage
func (t *TaskManager) UpdateWithMessage(message *types.Message, task *types.Task) *types.Task {
	if task.Status.Message != nil {
		task.History = append(task.History, task.Status.Message)
	}

	task.History = append(task.History, message)
	t.currentTask = task
	return task
}

func (t *TaskManager) Process(ctx context.Context, event types.Event) (types.Event, error) {
	_, err := t.SaveTaskEvent(ctx, event)
	if err != nil {
		return event, err
	}
	return event, nil
}

func (t *TaskManager) saveTask(ctx context.Context, task *types.Task) error {
	err := t.store.Save(ctx, task)
	if err != nil {
		return err
	}
	t.currentTask = task

	if t.taskId == "" {
		t.taskId = task.Id
		t.contextId = task.ContextId
	}
	return nil
}

func (t *TaskManager) EnsureTask(ctx context.Context, event types.Event) (*types.Task, error) {
	task := t.currentTask
	if task == nil && t.taskId != "" {
		newTask, err := t.store.Get(ctx, t.taskId)
		if err != nil {
			return nil, err
		}
		task = newTask
	}
	if task == nil {
		newTask := t.initTask(event.GetTaskId(), event.GetContextId())
		err := t.saveTask(ctx, newTask)
		if err != nil {
			return nil, err
		}
		return newTask, nil
	}
	return task, nil
}

func (t *TaskManager) initTask(taskId string, contextId string) *types.Task {
	task := &types.Task{
		Id:        taskId,
		ContextId: contextId,
		Status:    types.TaskStatus{State: types.SUBMITTED},
	}

	if t.initMessage != nil {
		task.History = append(task.History, t.initMessage)
	}
	return task
}
