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

package updater

import (
	"time"

	"github.com/google/uuid"
	"github.com/yumosx/a2a-go/sdk/server/event"
	"github.com/yumosx/a2a-go/sdk/types"
)

type TaskUpdater struct {
	queue     *event.Queue
	taskId    string
	contextId string
}

func NewTaskUpdater(queue *event.Queue, taskId string, contextId string) *TaskUpdater {
	return &TaskUpdater{queue: queue, taskId: taskId, contextId: contextId}
}

func (t *TaskUpdater) UpdateStatus(state types.TaskState, opts ...TaskUpdaterOption) {
	option := &TaskUpdaterOptions{}
	for _, opt := range opts {
		opt.Option(option)
	}

	if option.timeStamp == "" {
		option.timeStamp = time.Now().Format(time.RFC3339)
	}

	t.queue.Enqueue(&types.TaskStatusUpdateEvent{
		TaskId:    t.taskId,
		ContextId: t.contextId,
		Final:     option.final,
		Status: types.TaskStatus{
			Message:   option.message,
			State:     state,
			TimeStamp: option.timeStamp,
		},
	})
}

func (t *TaskUpdater) AddArtifact(parts []types.Part, opts ...TaskUpdaterOption) {
	option := &TaskUpdaterOptions{}
	for _, opt := range opts {
		opt.Option(option)
	}

	if option.artifactId == "" {
		option.artifactId = uuid.New().String()
	}

	t.queue.Enqueue(&types.TaskArtifactUpdateEvent{
		TaskId:    t.taskId,
		ContextId: t.contextId,
		Artifact: &types.Artifact{
			ArtifactId: option.artifactId,
			Name:       option.name,
			Parts:      parts,
			Metadata:   option.metadata,
		},
	})
}

func (t *TaskUpdater) Complete(opts ...TaskUpdaterOption) {
	t.UpdateStatus(types.COMPLETED, opts...)
}

func (t *TaskUpdater) Failed(opts ...TaskUpdaterOption) {
	t.UpdateStatus(types.FAILED, opts...)
}

func (t *TaskUpdater) Reject(opts ...TaskUpdaterOption) {
	t.UpdateStatus(types.REJECTED, opts...)
}

func (t *TaskUpdater) Submit(opts ...TaskUpdaterOption) {
	t.UpdateStatus(types.SUBMITTED, opts...)
}

func (t *TaskUpdater) StartWork(opts ...TaskUpdaterOption) {
	t.UpdateStatus(types.WORKING, opts...)
}

// NewAgentMessage create a new message object sent by the agent for this task/context
func (t *TaskUpdater) NewAgentMessage(parts []types.Part, opts ...TaskUpdaterOption) types.Message {
	options := &TaskUpdaterOptions{}

	for _, opt := range opts {
		opt.Option(options)
	}

	return types.Message{
		Role:      types.Agent,
		TaskID:    t.taskId,
		ContextID: t.contextId,
		MessageID: uuid.New().String(),
		Parts:     parts,
		Metadata:  options.metadata,
	}
}

type TaskUpdaterOptions struct {
	message    *types.Message
	metadata   map[string]any
	final      bool
	artifactId string
	name       string
	timeStamp  string
}

type TaskUpdaterOption interface {
	Option(t *TaskUpdaterOptions)
}

type TaskUpdaterOptionFunc func(t *TaskUpdaterOptions)

func (fn TaskUpdaterOptionFunc) Option(t *TaskUpdaterOptions) {
	fn(t)
}

func WithMessage(message *types.Message) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.message = message
	})
}

func WithMetadata(metadata map[string]any) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.metadata = metadata
	})
}

func WithFinal(final bool) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.final = final
	})
}

func WithArtifactId(id string) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.artifactId = id
	})
}

func WithName(name string) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.name = name
	})
}
func WithTimestamp(timeStamp string) TaskUpdaterOption {
	return TaskUpdaterOptionFunc(func(t *TaskUpdaterOptions) {
		t.timeStamp = timeStamp
	})
}
