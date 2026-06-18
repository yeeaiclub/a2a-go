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

import "github.com/yeeaiclub/a2a-go/sdk/types"

// TaskManagerOption is an option for configuring TaskManager.
type TaskManagerOption interface {
	Option(manager *TaskManager)
}

type TaskManagerOptionFunc func(manager *TaskManager)

func (fn TaskManagerOptionFunc) Option(manger *TaskManager) {
	fn(manger)
}

// WithTaskId sets the task ID for the TaskManager.
func WithTaskId(taskId string) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.taskId = taskId
	})
}

// WithContextId sets the context ID for the TaskManager.
func WithContextId(contextId string) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.contextId = contextId
	})
}

// WithInitMessage sets the initial message for the TaskManager.
func WithInitMessage(message *types.Message) TaskManagerOption {
	return TaskManagerOptionFunc(func(manger *TaskManager) {
		manger.initMessage = message
	})
}
